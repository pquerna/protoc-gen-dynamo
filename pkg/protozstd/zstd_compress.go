package protozstd

import (
	"github.com/klauspost/compress/zstd"
)

// getEncoder gets a zstd encoder from the pool
func (o *MarshalOptions) getEncoder() *zstd.Encoder {
	return o.EncoderPool.Get().(*zstd.Encoder)
}

// putEncoder returns a zstd encoder to the pool
func (o *MarshalOptions) putEncoder(enc *zstd.Encoder) {
	enc.Reset(nil)
	o.EncoderPool.Put(enc)
}

// getDecoder gets a zstd decoder from the pool
func (o *UnmarshalOptions) getDecoder() *zstd.Decoder {
	return o.DecoderPool.Get().(*zstd.Decoder)
}

// putDecoder returns a zstd decoder to the pool
func (o *UnmarshalOptions) putDecoder(dec *zstd.Decoder) {
	dec.Reset(nil)
	o.DecoderPool.Put(dec)
}

// IsCompressed checks if the data is zstd compressed
func (o *UnmarshalOptions) IsCompressed(data []byte) bool {
	return o.isCompressed(data)
}

// isCompressed is an internal helper function to check if data is compressed
func (o *UnmarshalOptions) isCompressed(data []byte) bool {
	// Check for zstd magic bytes (0x28 0xB5 0x2F 0xFD) in little-endian format
	return len(data) >= 4 &&
		data[0] == 0x28 &&
		data[1] == 0xB5 &&
		data[2] == 0x2F &&
		data[3] == 0xFD
}

// CompressValue compresses data using zstd if it exceeds MinCompressSize
func (o *MarshalOptions) compressValue(data []byte) ([]byte, error) {
	if o.DisableCompression || len(data) <= o.MinimumSizeToCompress {
		return data, nil
	}

	enc := o.getEncoder()
	defer o.putEncoder(enc)

	return enc.EncodeAll(data, nil), nil
}

// DecompressValue decompresses zstd compressed data
func (o *UnmarshalOptions) decompressValue(data []byte) ([]byte, error) {
	if !o.isCompressed(data) {
		return data, nil
	}

	dec := o.getDecoder()
	defer o.putDecoder(dec)

	return dec.DecodeAll(data, nil)
}
