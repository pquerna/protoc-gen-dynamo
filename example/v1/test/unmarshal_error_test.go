package v1_test

import (
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/klauspost/compress/zstd"

	v1 "github.com/pquerna/protoc-gen-dynamo/example/v1"
)

// oversizedZstdFrame builds a zstd frame header declaring fcs bytes of content and
// carrying no blocks. Single_Segment is left clear so the window descriptor supplies
// a legal 1MiB window; with it set, WindowSize takes the declared content size and
// the window check rejects the frame before the decoded-size check is reached.
//
// Declaring more than the decoder's own maximum means the frame is rejected before
// an output buffer is sized, so this costs no allocation.
func oversizedZstdFrame(fcs uint64) []byte {
	const (
		fcs8Byte        = 3 << 6
		window1MiBExp10 = 10 << 3
	)
	frame := []byte{0x28, 0xB5, 0x2F, 0xFD, fcs8Byte, window1MiBExp10}
	return binary.LittleEndian.AppendUint64(frame, fcs)
}

func itemFor(typeName string, value []byte) types.AttributeValue {
	return &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"typ":   &types.AttributeValueMemberS{Value: typeName},
		"value": &types.AttributeValueMemberB{Value: value},
	}}
}

// A decode failure used to return protozstd.Unmarshal's error unwrapped, which
// names neither the message nor the input. A caller saw only "decompressed size
// exceeds configured limit" with nothing to identify what had failed.
func TestUnmarshalDecodeErrorNamesMessageAndSize(t *testing.T) {
	frame := oversizedZstdFrame(1 << 40)

	p := &v1.Store{}
	err := p.UnmarshalDynamoDBAttributeValue(itemFor("examplepb.v1.Store", frame))
	if err == nil {
		t.Fatal("expected a decode error")
	}

	for _, want := range []string{"examplepb.v1.Store", "14 compressed bytes"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
}

// Wrapping must not cost callers the ability to classify the failure: a
// decoded-size rejection has to stay distinguishable from a corrupt frame.
func TestUnmarshalDecodeErrorStaysUnwrappable(t *testing.T) {
	frame := oversizedZstdFrame(1 << 40)

	p := &v1.Store{}
	err := p.UnmarshalDynamoDBAttributeValue(itemFor("examplepb.v1.Store", frame))
	if !errors.Is(err, zstd.ErrDecoderSizeExceeded) {
		t.Fatalf("errors.Is lost zstd.ErrDecoderSizeExceeded through the wrap: %v", err)
	}
}

// UnmarshalDynamoItem and UnmarshalDynamo both delegate to
// UnmarshalDynamoDBAttributeValue, so they inherit the context.
func TestUnmarshalDynamoItemInheritsTheContext(t *testing.T) {
	frame := oversizedZstdFrame(1 << 40)

	p := &v1.Store{}
	err := p.UnmarshalDynamoItem(map[string]types.AttributeValue{
		"typ":   &types.AttributeValueMemberS{Value: "examplepb.v1.Store"},
		"value": &types.AttributeValueMemberB{Value: frame},
	})
	if err == nil {
		t.Fatal("expected a decode error")
	}
	if !strings.Contains(err.Error(), "examplepb.v1.Store") {
		t.Errorf("error %q does not name the message", err)
	}
}
