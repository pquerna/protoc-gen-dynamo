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
		New: func() interface{} {
			encoder, err := zstd.NewWriter(nil,
				zstd.WithEncoderLevel(zstd.SpeedFastest),
				zstd.WithEncoderConcurrency(1),
			)
			if err != nil {
				panic(err)
			}
			return encoder
		},
	}

	// decoderPool is a pool of zstd decoders
	decoderPool = sync.Pool{
		New: func() interface{} {
			decoder, err := zstd.NewReader(nil,
				zstd.WithDecoderConcurrency(1),
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

// ZstdCompress compresses data using zstd if it exceeds MinCompressSize
func ZstdCompress(data []byte) ([]byte, error) {
	if len(data) <= minCompressSize {
		return data, nil
	}

	enc := getEncoder()
	defer putEncoder(enc)

	return enc.EncodeAll(data, nil), nil
}

// ZstdDecompress decompresses zstd compressed data
func ZstdDecompress(data []byte) ([]byte, error) {
	if !zstdIsCompressed(data) {
		return data, nil
	}

	dec := getDecoder()
	defer putDecoder(dec)

	return dec.DecodeAll(data, nil)
}
