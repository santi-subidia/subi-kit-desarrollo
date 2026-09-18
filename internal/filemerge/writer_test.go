package filemerge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomic_NewFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-atomic-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	data := []byte("hello world")

	res, err := WriteFileAtomic(filePath, data, 0644)
	if err != nil {
		t.Fatalf("WriteFileAtomic failed: %v", err)
	}

	if !res.Created || !res.Changed {
		t.Errorf("expected Created=true, Changed=true, got %+v", res)
	}

	readBack, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if !bytes.Equal(readBack, data) {
		t.Errorf("expected %q, got %q", string(data), string(readBack))
	}
}

func TestWriteFileAtomic_Idempotent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-atomic-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	data := []byte("hello world")

	// First write
	_, err = WriteFileAtomic(filePath, data, 0644)
	if err != nil {
		t.Fatalf("first WriteFileAtomic failed: %v", err)
	}

	// Second write with identical content
	res, err := WriteFileAtomic(filePath, data, 0644)
	if err != nil {
		t.Fatalf("second WriteFileAtomic failed: %v", err)
	}

	if res.Changed || res.Created {
		t.Errorf("expected idempotent no-op (Changed=false, Created=false), got %+v", res)
	}
}

func TestWriteFileAtomic_Update(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-atomic-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	data1 := []byte("hello world")
	data2 := []byte("updated content")

	_, _ = WriteFileAtomic(filePath, data1, 0644)

	res, err := WriteFileAtomic(filePath, data2, 0644)
	if err != nil {
		t.Fatalf("update WriteFileAtomic failed: %v", err)
	}

	if !res.Changed || res.Created {
		t.Errorf("expected Changed=true, Created=false, got %+v", res)
	}

	readBack, _ := os.ReadFile(filePath)
	if !bytes.Equal(readBack, data2) {
		t.Errorf("expected %q, got %q", string(data2), string(readBack))
	}
}

func TestWriteFileAtomic_NestedDirCreation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-atomic-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "nested", "deep", "dir", "test.txt")
	data := []byte("deep content")

	res, err := WriteFileAtomic(filePath, data, 0644)
	if err != nil {
		t.Fatalf("WriteFileAtomic failed: %v", err)
	}

	if !res.Created || !res.Changed {
		t.Errorf("expected Created=true, Changed=true, got %+v", res)
	}

	readBack, _ := os.ReadFile(filePath)
	if !bytes.Equal(readBack, data) {
		t.Errorf("expected %q, got %q", string(data), string(readBack))
	}
}

func TestWriteStreamAtomic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-atomic-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "stream.txt")
	data := []byte("streamed content to verify SHA256 digest")
	expectedDigest := sha256.Sum256(data)
	expectedDigestHex := hex.EncodeToString(expectedDigest[:])

	res, err := WriteStreamAtomic(filePath, bytes.NewReader(data), 0644)
	if err != nil {
		t.Fatalf("WriteStreamAtomic failed: %v", err)
	}

	if res.Bytes != int64(len(data)) {
		t.Errorf("expected %d bytes, got %d", len(data), res.Bytes)
	}
	if res.Digest != expectedDigestHex {
		t.Errorf("expected digest %s, got %s", expectedDigestHex, res.Digest)
	}
}

func TestFileDigest(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-atomic-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "digest.txt")
	data := []byte("compute my digest accurately")
	expectedDigest := sha256.Sum256(data)
	expectedDigestHex := hex.EncodeToString(expectedDigest[:])

	_, err = WriteFileAtomic(filePath, data, 0644)
	if err != nil {
		t.Fatalf("WriteFileAtomic failed: %v", err)
	}

	digest, size, err := FileDigest(filePath)
	if err != nil {
		t.Fatalf("FileDigest failed: %v", err)
	}

	if size != int64(len(data)) {
		t.Errorf("expected size %d, got %d", len(data), size)
	}
	if digest != expectedDigestHex {
		t.Errorf("expected digest %s, got %s", expectedDigestHex, digest)
	}
}
