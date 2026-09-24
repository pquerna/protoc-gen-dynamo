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

// codecConcurrency pins both pools to single-worker codecs.
//
// This package only ever calls EncodeAll / DecodeAll, which take one codec off the zstd
// object's internal channel and run synchronously in the caller's goroutine, so per-codec
// concurrency buys nothing. Left unset it costs a full set of internal states per pooled
// object -- GOMAXPROCS encoder states (~270 KB each at SpeedFastest) plus GOMAXPROCS block
// decoders -- and sync.Pool discards those on every GC, so a busy process rebuilds them
// continuously. Concurrency is not part of the zstd frame format, so output is
// byte-identical either way.
//
// Spelled out rather than left to zstd, whose encoder defaults to GOMAXPROCS and whose
// decoder defaults to min(4, GOMAXPROCS). Do not reach for 0 here: the decoder reads it as
// GOMAXPROCS rather than "library default", and the encoder rejects it outright, which
// encoderConstruct turns into a panic.
const codecConcurrency = 1

func NewMarshalOptions() *MarshalOptions {
	mo := &MarshalOptions{
		MarshalOptions:     proto.MarshalOptions{},
		DisableCompression: false,
		EncoderOptions: []zstd.EOption{
			zstd.WithEncoderLevel(zstd.SpeedFastest),
			zstd.WithEncoderConcurrency(codecConcurrency),
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
			zstd.WithDecoderConcurrency(codecConcurrency),
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
