package protozstd

import (
	"bytes"
	"sync"
	"testing"

	"github.com/klauspost/compress/zstd"
)

func TestMarshalOptions_SmallData(t *testing.T) {
	// Create options with minimum size set to ensure small data won't be compressed
	opts := NewMarshalOptions()
	opts.MinimumSizeToCompress = 200

	// Data smaller than MinimumSizeToCompress should not be compressed
	smallData := []byte("This is a small test string that is less than 200 bytes")

	compressed, err := opts.compressValue(smallData)
	if err != nil {
		t.Fatalf("compressValue failed: %v", err)
	}

	// The returned data should be identical to the input
	if !bytes.Equal(compressed, smallData) {
		t.Errorf("Small data was compressed when it should not have been")
	}

	// Create unmarshal options
	uOpts := NewUnmarshalOptions()

	// Verify it's not detected as compressed
	if uOpts.isCompressed(compressed) {
		t.Errorf("Small data incorrectly identified as compressed")
	}
}

func TestMarshalOptions_LargeData(t *testing.T) {
	// Create options with minimum size set to ensure large data will be compressed
	opts := NewMarshalOptions()
	opts.MinimumSizeToCompress = 200

	// Create data larger than MinimumSizeToCompress
	largeData := make([]byte, 500)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	compressed, err := opts.compressValue(largeData)
	if err != nil {
		t.Fatalf("compressValue failed: %v", err)
	}

	// The compressed data should be different from the input
	if bytes.Equal(compressed, largeData) {
		t.Errorf("Large data was not compressed")
	}

	// Create unmarshal options
	uOpts := NewUnmarshalOptions()

	// Verify it's detected as compressed
	if !uOpts.isCompressed(compressed) {
		t.Errorf("Compressed data not correctly identified")
	}

	// Verify we can decompress it back
	decompressed, err := uOpts.decompressValue(compressed)
	if err != nil {
		t.Fatalf("decompressValue failed: %v", err)
	}

	// The decompressed data should match the original
	if !bytes.Equal(decompressed, largeData) {
		t.Errorf("Decompressed data doesn't match original")
	}
}

func TestUnmarshalOptions_UncompressedData(t *testing.T) {
	// Create unmarshal options
	opts := NewUnmarshalOptions()

	// Uncompressed data should be returned as-is
	originalData := []byte("This is not compressed data")

	result, err := opts.decompressValue(originalData)
	if err != nil {
		t.Fatalf("decompressValue failed: %v", err)
	}

	if !bytes.Equal(result, originalData) {
		t.Errorf("Uncompressed data was modified by decompressValue")
	}
}

func TestUnmarshalOptions_IsCompressed(t *testing.T) {
	// Create unmarshal options
	opts := NewUnmarshalOptions()

	// Test with invalid/empty data
	if opts.isCompressed(nil) {
		t.Errorf("nil data incorrectly identified as compressed")
	}

	if opts.isCompressed([]byte{}) {
		t.Errorf("Empty data incorrectly identified as compressed")
	}

	if opts.isCompressed([]byte{0x01, 0x02, 0x03}) {
		t.Errorf("Too short data incorrectly identified as compressed")
	}

	// Test with data that has incorrect magic bytes
	wrongMagic := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0x01, 0x02}
	if opts.isCompressed(wrongMagic) {
		t.Errorf("Data with incorrect magic bytes incorrectly identified as compressed")
	}

	// Test with data that has correct magic bytes
	correctMagic := []byte{0x28, 0xB5, 0x2F, 0xFD, 0x01, 0x02}
	if !opts.isCompressed(correctMagic) {
		t.Errorf("Data with correct magic bytes not identified as compressed")
	}
}

// Add tests for DisableCompression flag
func TestMarshalOptions_DisableCompression(t *testing.T) {
	opts := NewMarshalOptions()
	opts.DisableCompression = true
	opts.MinimumSizeToCompress = 10 // Set low to ensure compression would happen otherwise

	// Create data that would normally be compressed
	largeData := make([]byte, 500)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	compressed, err := opts.compressValue(largeData)
	if err != nil {
		t.Fatalf("compressValue failed: %v", err)
	}

	// Data should not be compressed when DisableCompression is true
	if !bytes.Equal(compressed, largeData) {
		t.Errorf("Data was compressed when DisableCompression was set to true")
	}

	// Create unmarshal options
	uOpts := NewUnmarshalOptions()

	// Verify it's not detected as compressed
	if uOpts.isCompressed(compressed) {
		t.Errorf("Data incorrectly identified as compressed when compression was disabled")
	}
}

// Define test proto message type for testing
type testProtoMessage struct {
	data []byte
}

// Implement the proto.Message interface
func (m *testProtoMessage) ProtoMessage()  {}
func (m *testProtoMessage) Reset()         { m.data = nil }
func (m *testProtoMessage) String() string { return string(m.data) }

// Test Marshal/Unmarshal integration with binary data
func TestMarshalUnmarshalIntegration(t *testing.T) {
	// Rather than trying to mock proto.Message which requires reflection support,
	// we'll test the compression/decompression functionality directly

	// Create Marshal/Unmarshal options
	mOpts := NewMarshalOptions()
	mOpts.MinimumSizeToCompress = 100 // Set small enough to ensure compression

	uOpts := NewUnmarshalOptions()

	// Create large test data
	testData := []byte("Test data for compression")
	for i := 0; i < 5; i++ {
		testData = append(testData, testData...) // Make it large
	}

	// Test compression
	compressed, err := mOpts.compressValue(testData)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}

	// Verify it was actually compressed
	if !uOpts.isCompressed(compressed) {
		t.Errorf("Data not correctly compressed")
	}

	// Decompress
	decompressed, err := uOpts.decompressValue(compressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}

	// Verify data integrity
	if !bytes.Equal(testData, decompressed) {
		t.Errorf("Data does not match after compression/decompression cycle")
	}

	// Test end-to-end simulating Marshal and Unmarshal without using actual proto messages
	// This effectively tests what would happen with integration with proto.Marshal/Unmarshal

	t.Run("End-to-end simulation", func(t *testing.T) {
		// 1. Simulate proto marshaling
		marshaledData := testData

		// 2. Apply our compression layer
		compressedData, err := mOpts.compressValue(marshaledData)
		if err != nil {
			t.Fatalf("compressValue failed: %v", err)
		}

		// 3. Simulate proto unmarshaling
		// First decompress
		decompressedData, err := uOpts.decompressValue(compressedData)
		if err != nil {
			t.Fatalf("decompressValue failed: %v", err)
		}

		// 4. Verify the unmarshaled result matches the original
		if !bytes.Equal(marshaledData, decompressedData) {
			t.Errorf("End-to-end cycle failed - data doesn't match")
		}
	})
}

// Add benchmark for compression/decompression
func BenchmarkCompressDecompress(b *testing.B) {
	mOpts := NewMarshalOptions()
	uOpts := NewUnmarshalOptions()

	// Create test data (1KB)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compressed, err := mOpts.compressValue(data)
		if err != nil {
			b.Fatalf("Compression failed: %v", err)
		}

		_, err = uOpts.decompressValue(compressed)
		if err != nil {
			b.Fatalf("Decompression failed: %v", err)
		}
	}
}

// Test custom encoder/decoder options
func TestCustomEncoderDecoderOptions(t *testing.T) {
	// Custom marshal options with different compression level
	mOpts := &MarshalOptions{
		MinimumSizeToCompress: 100,
		EncoderOptions: []zstd.EOption{
			zstd.WithEncoderLevel(zstd.SpeedDefault), // Different from default
		},
	}

	// Initialize the pool
	mOpts.EncoderPool = &sync.Pool{
		New: func() any {
			return mOpts.encoderConstruct()
		},
	}

	// Test compression with custom options
	data := make([]byte, 200)

	compressed, err := mOpts.compressValue(data)
	if err != nil {
		t.Fatalf("Compression with custom options failed: %v", err)
	}

	// Should be compressed
	uOpts := NewUnmarshalOptions()
	if !uOpts.isCompressed(compressed) {
		t.Errorf("Data compressed with custom options not detected as compressed")
	}
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	// Test corrupt data handling
	t.Run("Corrupt data handling", func(t *testing.T) {
		// Create corrupt compressed data - valid magic bytes but corrupt data
		corruptData := []byte{0x28, 0xB5, 0x2F, 0xFD, 0xFF, 0xFF, 0xFF, 0xFF}

		uOpts := NewUnmarshalOptions()

		// This should result in an error
		_, err := uOpts.decompressValue(corruptData)
		if err == nil {
			t.Errorf("Expected error when decompressing corrupt data but got nil")
		}
	})

	// Test malformed zstd data
	t.Run("Malformed zstd data", func(t *testing.T) {
		// Create data with invalid zstd format but valid magic
		invalidData := []byte{0x28, 0xB5, 0x2F, 0xFD, 0x00, 0x00}

		uOpts := NewUnmarshalOptions()

		// Decompression should fail
		_, err := uOpts.decompressValue(invalidData)
		if err == nil {
			t.Errorf("Expected error with malformed zstd data but got nil")
		}
	})
}
