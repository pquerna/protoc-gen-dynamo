package protozstd

import (
	"os"
	"strconv"
	"sync"

	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

type MarshalOptions struct {
	proto.MarshalOptions
	DisableCompression    bool
	MinimumSizeToCompress int

	EncoderPool    *sync.Pool
	EncoderOptions []zstd.EOption
}

type UnmarshalOptions struct {
	proto.UnmarshalOptions
	DecoderPool    *sync.Pool
	DecoderOptions []zstd.DOption
}

var DefaultMarshalOptions = NewMarshalOptions()
var DefaultUnmarshalOptions = NewUnmarshalOptions()

func NewMarshalOptions() *MarshalOptions {
	mo := &MarshalOptions{
		MarshalOptions:     proto.MarshalOptions{},
		DisableCompression: false,
		EncoderOptions: []zstd.EOption{
			zstd.WithEncoderLevel(zstd.SpeedFastest),
		},
	}

	mo.EncoderPool = &sync.Pool{
		New: func() any {
			return mo.encoderConstruct()
		},
	}
	return mo
}

func NewUnmarshalOptions() *UnmarshalOptions {
	uo := &UnmarshalOptions{
		UnmarshalOptions: proto.UnmarshalOptions{},
		DecoderOptions: []zstd.DOption{
			zstd.WithDecoderConcurrency(0),
		},
	}
	uo.DecoderPool = &sync.Pool{
		New: func() any {
			return uo.decoderConstruct()
		},
	}
	return uo
}

func (o *MarshalOptions) encoderConstruct() *zstd.Encoder {
	encoder, err := zstd.NewWriter(nil,
		o.EncoderOptions...,
	)
	if err != nil {
		panic(err)
	}
	return encoder
}

func (o *UnmarshalOptions) decoderConstruct() *zstd.Decoder {
	decoder, err := zstd.NewReader(nil,
		o.DecoderOptions...,
	)
	if err != nil {
		panic(err)
	}
	return decoder
}

func (o *UnmarshalOptions) Unmarshal(data []byte, m proto.Message) error {
	if o.isCompressed(data) {
		var err error
		data, err = o.decompressValue(data)
		if err != nil {
			return err
		}
	}
	return o.UnmarshalOptions.Unmarshal(data, m)
}

func (o *MarshalOptions) Marshal(m proto.Message) ([]byte, error) {
	data, err := o.MarshalOptions.Marshal(m)
	if err != nil {
		return nil, err
	}
	return o.compressValue(data)
}

const (
	envDisableCompression = "PROTOZSTD_DISABLE_COMPRESSION"
)

func init() {
	// disableCompression is a flag to disable compression
	//
	// NOTE: we use an environment variable here because
	// there is not a context or other way to propagate what we
	// want down to the Marshaler interfaces. Not ideal!
	//
	// NOTE: only disables compression of new data --
	// any previously compressed data will still be decompressed
	// on read.
	var disableCompression, _ = strconv.ParseBool(os.Getenv(envDisableCompression))
	DefaultMarshalOptions.DisableCompression = disableCompression
}

func Marshal(m proto.Message) ([]byte, error) {
	return DefaultMarshalOptions.Marshal(m)
}

func Unmarshal(data []byte, m proto.Message) error {
	return DefaultUnmarshalOptions.Unmarshal(data, m)
}
