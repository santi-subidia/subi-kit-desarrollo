package filemerge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

var (
	runtimeGOOS = func() string { return runtime.GOOS }
	renameFn    = os.Rename
	syncDirFn   = func(dir string) error {
		fd, err := os.Open(dir)
		if err != nil {
			return fmt.Errorf("open parent directory %q: %w", dir, err)
		}
		defer fd.Close()
		return fd.Sync()
	}
)

type stagedFile interface {
	io.Writer
	Name() string
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
}

var createStagedFile = func(dir string) (stagedFile, error) {
	return os.CreateTemp(dir, ".subikit-*.tmp")
}

// WriteResult reports what happened to the destination path.
type WriteResult struct {
	Changed bool
	Created bool
}

// StreamResult describes the bytes that landed at the destination.
type StreamResult struct {
	Bytes  int64
	Digest string
}

// WriteFileAtomic replaces path with content and reports whether the replacement actually occurred.
// It is idempotent: if the file already exists and its contents are identical, no write occurs.
// It writes to a temporary file in the same directory, syncs to disk, atomically renames to path,
// verifies by reading back the SHA-256 digest, and syncs the parent directory.
func WriteFileAtomic(path string, content []byte, perm fs.FileMode) (WriteResult, error) {
	if perm == 0 {
		perm = 0o644
	}

	created := false
	existing, err := os.ReadFile(path)
	if err == nil {
		if bytes.Equal(existing, content) {
			return WriteResult{Changed: false, Created: false}, nil
		}
	} else if !os.IsNotExist(err) {
		return WriteResult{}, fmt.Errorf("read existing file %q: %w", path, err)
	} else {
		created = true
	}

	landed, _, err := replaceDurably(path, bytes.NewReader(content), perm)
	result := WriteResult{Changed: landed, Created: created && landed}
	if err != nil {
		return result, err
	}
	return result, nil
}

// WriteStreamAtomic replaces path with everything readable from src and returns
// the size and SHA-256 of the bytes read back from path afterwards.
func WriteStreamAtomic(path string, src io.Reader, perm fs.FileMode) (StreamResult, error) {
	_, result, err := replaceDurably(path, src, perm)
	return result, err
}

func replaceDurably(path string, src io.Reader, perm fs.FileMode) (landed bool, result StreamResult, err error) {
	if perm == 0 {
		perm = 0o644
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, StreamResult{}, fmt.Errorf("create parent directory %q: %w", dir, err)
	}

	tmp, err := createStagedFile(dir)
	if err != nil {
		return false, StreamResult{}, fmt.Errorf("create temp file for %q: %w", path, err)
	}

	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	staged := sha256.New()
	written, err := io.Copy(io.MultiWriter(tmp, staged), src)
	if err != nil {
		_ = tmp.Close()
		return false, StreamResult{}, fmt.Errorf("write temp file for %q: %w", path, err)
	}
	stagedDigest := hex.EncodeToString(staged.Sum(nil))

	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return false, StreamResult{}, fmt.Errorf("set permissions on temp file for %q: %w", path, err)
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return false, StreamResult{}, fmt.Errorf("sync temp file for %q: %w", path, err)
	}

	if err := tmp.Close(); err != nil {
		return false, StreamResult{}, fmt.Errorf("close temp file for %q: %w", path, err)
	}

	if err := renameFn(tmpPath, path); err != nil {
		return false, StreamResult{}, fmt.Errorf("replace %q atomically: %w", path, err)
	}

	// Verify the destination on disk
	diskDigest, diskBytes, err := digestFileOnDisk(path)
	if err != nil {
		return false, StreamResult{}, fmt.Errorf("read back %q after replacement: %w", path, err)
	}
	if diskBytes != written || diskDigest != stagedDigest {
		return false, StreamResult{}, fmt.Errorf(
			"replace %q atomically: verification failed. Disk holds %d bytes (%s); %d bytes (%s) were written",
			path, diskBytes, diskDigest, written, stagedDigest,
		)
	}

	cleanup = false
	result = StreamResult{Bytes: diskBytes, Digest: diskDigest}

	if err := SyncDir(dir); err != nil {
		return true, result, fmt.Errorf("sync parent directory for %q: %w", path, err)
	}

	return true, result, nil
}

// SyncDir flushes dir entries to stable storage.
// Tolerates ErrPermission on Windows because NTFS refuses directory handle Sync().
func SyncDir(dir string) error {
	err := syncDirFn(dir)
	if err == nil {
		return nil
	}
	if runtimeGOOS() == "windows" && errors.Is(err, os.ErrPermission) {
		return nil
	}
	return err
}

// FileDigest returns the SHA-256 and size of the file at path read from disk.
func FileDigest(path string) (digest string, size int64, err error) {
	return digestFileOnDisk(path)
}

func digestFileOnDisk(path string) (string, int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", 0, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", 0, fmt.Errorf("destination %q is a symlink", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	hasher := sha256.New()
	bytesCopied, err := io.Copy(hasher, file)
	if err != nil {
		return "", 0, fmt.Errorf("read %q for digest: %w", path, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), bytesCopied, nil
}
