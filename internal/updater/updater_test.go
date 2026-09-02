package updater

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"v0.4.0", "v0.4.0", 0},
		{"0.4.0", "0.4.0", 0},
		{"v0.5.0", "v0.4.0", 1},
		{"v0.4.0", "v0.5.0", -1},
		{"v1.0.0", "v0.9.9", 1},
		{"v0.4.1", "v0.4.0", 1},
		{"v0.4.0", "v0.4.1", -1},
		{"v1.2.3-alpha", "v1.2.2", 1},
		{"v1.2.3", "v1.2.3-beta", 0},
	}

	for _, tt := range tests {
		got := CompareSemver(tt.v1, tt.v2)
		if got != tt.expected {
			t.Errorf("CompareSemver(%q, %q) = %d; want %d", tt.v1, tt.v2, got, tt.expected)
		}
	}
}

func TestFindAssetForPlatform(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "subikit-windows-amd64.zip", BrowserDownloadURL: "https://example.com/win64.zip"},
		{Name: "subikit-linux-amd64.tar.gz", BrowserDownloadURL: "https://example.com/linux64.tar.gz"},
		{Name: "subikit-linux-arm64.tar.gz", BrowserDownloadURL: "https://example.com/linuxarm64.tar.gz"},
		{Name: "subikit-darwin-arm64.tar.gz", BrowserDownloadURL: "https://example.com/macarm64.tar.gz"},
	}

	// Test Windows amd64
	asset, err := FindAssetForPlatform(assets, "windows", "amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.Name != "subikit-windows-amd64.zip" {
		t.Errorf("got asset %q, want subikit-windows-amd64.zip", asset.Name)
	}

	// Test Linux arm64
	asset, err = FindAssetForPlatform(assets, "linux", "arm64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.Name != "subikit-linux-arm64.tar.gz" {
		t.Errorf("got asset %q, want subikit-linux-arm64.tar.gz", asset.Name)
	}

	// Test Incompatible platform
	_, err = FindAssetForPlatform(assets, "freebsd", "386")
	if err == nil {
		t.Errorf("expected error for unsupported platform, got nil")
	}
}

func TestExtractBinaryZip(t *testing.T) {
	// Create an in-memory zip containing subikit.exe
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, err := w.Create("subikit.exe")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	expectedContent := "mock-binary-content-123"
	if _, err := f.Write([]byte(expectedContent)); err != nil {
		t.Fatalf("failed to write zip entry: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	extracted, err := extractBinary(buf.Bytes(), "subikit-windows-amd64.zip", "subikit.exe")
	if err != nil {
		t.Fatalf("extractBinary failed: %v", err)
	}

	if string(extracted) != expectedContent {
		t.Errorf("got content %q, want %q", string(extracted), expectedContent)
	}
}
