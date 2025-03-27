package v1

import (
	"sync"

	"github.com/klauspost/compress/zstd"
)

const (
	// minCompressSize is the minimum size of a protobuf message before compression is applied
	minCompressSize = 200
)

var (
	// encoderPool is a pool of zstd encoders
	encoderPool = sync.Pool{
		New: func() any {
			encoder, err := zstd.NewWriter(nil,
				zstd.WithEncoderLevel(zstd.SpeedFastest),
				// NOTE: similair to decoder concurrency
				// we want to use up to all cores, but its behavoir
				// is a little different in its default value vs 0,
				// so we don't set it here
				// zstd.WithEncoderConcurrency(1),
			)
			if err != nil {
				panic(err)
			}
			return encoder
		},
	}

	// decoderPool is a pool of zstd decoders
	decoderPool = sync.Pool{
		New: func() any {
			decoder, err := zstd.NewReader(nil,
				// allow concurrent decoders to use all cores
				zstd.WithDecoderConcurrency(0),
			)
			if err != nil {
				panic(err)
			}
			return decoder
		},
	}
)

// getEncoder gets a zstd encoder from the pool
func getEncoder() *zstd.Encoder {
	return encoderPool.Get().(*zstd.Encoder)
}

// putEncoder returns a zstd encoder to the pool
func putEncoder(enc *zstd.Encoder) {
	enc.Reset(nil)
	encoderPool.Put(enc)
}

// getDecoder gets a zstd decoder from the pool
func getDecoder() *zstd.Decoder {
	return decoderPool.Get().(*zstd.Decoder)
}

// putDecoder returns a zstd decoder to the pool
func putDecoder(dec *zstd.Decoder) {
	dec.Reset(nil)
	decoderPool.Put(dec)
}

// zstdIsCompressed checks if the data is zstd compressed
func zstdIsCompressed(data []byte) bool {
	// Check for zstd magic bytes (0x28 0xB5 0x2F 0xFD) in little-endian format
	return len(data) >= 4 &&
		data[0] == 0x28 &&
		data[1] == 0xB5 &&
		data[2] == 0x2F &&
		data[3] == 0xFD
}

// CompressValue compresses data using zstd if it exceeds MinCompressSize
func CompressValue(data []byte) ([]byte, error) {
	if len(data) <= minCompressSize {
		return data, nil
	}

	enc := getEncoder()
	defer putEncoder(enc)

	return enc.EncodeAll(data, nil), nil
}

// DecompressValue decompresses zstd compressed data
func DecompressValue(data []byte) ([]byte, error) {
	if !zstdIsCompressed(data) {
		return data, nil
	}

	dec := getDecoder()
	defer putDecoder(dec)

	return dec.DecodeAll(data, nil)
}
