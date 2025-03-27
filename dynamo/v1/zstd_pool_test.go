package v1

import (
	"bytes"
	"testing"
)

func TestCompressValue_SmallData(t *testing.T) {
	// Data smaller than minCompressSize should not be compressed
	smallData := []byte("This is a small test string that is less than 200 bytes")

	compressed, err := CompressValue(smallData)
	if err != nil {
		t.Fatalf("CompressValue failed: %v", err)
	}

	// The returned data should be identical to the input
	if !bytes.Equal(compressed, smallData) {
		t.Errorf("Small data was compressed when it should not have been")
	}

	// Verify it's not detected as compressed
	if zstdIsCompressed(compressed) {
		t.Errorf("Small data incorrectly identified as compressed")
	}
}

func TestCompressValue_LargeData(t *testing.T) {
	// Create data larger than minCompressSize
	largeData := make([]byte, 500)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	compressed, err := CompressValue(largeData)
	if err != nil {
		t.Fatalf("CompressValue failed: %v", err)
	}

	// The compressed data should be different from the input
	if bytes.Equal(compressed, largeData) {
		t.Errorf("Large data was not compressed")
	}

	// Verify it's detected as compressed
	if !zstdIsCompressed(compressed) {
		t.Errorf("Compressed data not correctly identified")
	}

	// Verify we can decompress it back
	decompressed, err := DecompressValue(compressed)
	if err != nil {
		t.Fatalf("DecompressValue failed: %v", err)
	}

	// The decompressed data should match the original
	if !bytes.Equal(decompressed, largeData) {
		t.Errorf("Decompressed data doesn't match original")
	}
}

func TestDecompressValue_UncompressedData(t *testing.T) {
	// Uncompressed data should be returned as-is
	originalData := []byte("This is not compressed data")

	result, err := DecompressValue(originalData)
	if err != nil {
		t.Fatalf("DecompressValue failed: %v", err)
	}

	if !bytes.Equal(result, originalData) {
		t.Errorf("Uncompressed data was modified by DecompressValue")
	}
}

func TestZstdIsCompressed(t *testing.T) {
	// Test with invalid/empty data
	if zstdIsCompressed(nil) {
		t.Errorf("nil data incorrectly identified as compressed")
	}

	if zstdIsCompressed([]byte{}) {
		t.Errorf("Empty data incorrectly identified as compressed")
	}

	if zstdIsCompressed([]byte{0x01, 0x02, 0x03}) {
		t.Errorf("Too short data incorrectly identified as compressed")
	}

	// Test with data that has incorrect magic bytes
	wrongMagic := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0x01, 0x02}
	if zstdIsCompressed(wrongMagic) {
		t.Errorf("Data with wrong magic bytes incorrectly identified as compressed")
	}

	// Test with data that has correct magic bytes
	correctMagic := []byte{0x28, 0xB5, 0x2F, 0xFD, 0x01, 0x02}
	if !zstdIsCompressed(correctMagic) {
		t.Errorf("Data with correct magic bytes not identified as compressed")
	}
}
