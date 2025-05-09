package v1_test

import (
	"strconv"
	"testing"
	"time"

	v1 "github.com/pquerna/protoc-gen-dynamo/example/v1"
	"github.com/pquerna/protoc-gen-dynamo/pkg/protozstd"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Test with a real protobuf message from example.pb.go
func TestRealProtoRoundTrip(t *testing.T) {
	// Create a Store message from the example proto
	now := time.Now()
	store := v1.Store_builder{
		Id:      proto.String("store123"),
		Country: proto.String("USA"),
		Region:  proto.String("West"),
		State:   proto.String("California"),
		City:    proto.String("San Francisco"),
		Closed:  proto.Bool(false),
		OpeningDate: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		BestEmployeeIds: []string{"emp1", "emp2", "emp3"},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Foo:     proto.Uint64(12345),
		Morefoo: []uint64{1, 2, 3, 4, 5},
	}.Build()

	// Perform a full round trip using our compression layer
	mOpts := protozstd.NewMarshalOptions()
	mOpts.MinimumSizeToCompress = 100 // Ensure compression happens

	// Apply our compression
	compressed, err := mOpts.Marshal(store)
	if err != nil {
		t.Fatalf("Failed to compress: %v", err)
	}

	// Verify it was actually compressed
	uOpts := protozstd.NewUnmarshalOptions()
	if !uOpts.IsCompressed(compressed) {
		t.Errorf("Data not correctly compressed")
	}

	// Decompress
	newStore := &v1.Store{}
	err = protozstd.Unmarshal(compressed, newStore)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	// Verify message equality
	if !proto.Equal(store, newStore) {
		t.Errorf("Round-tripped protobuf message doesn't match original")
	}

	// Now test using the full Marshal/Unmarshal methods
	t.Run("Full Marshal-Unmarshal", func(t *testing.T) {
		// Marshal with our wrapper
		data, err := mOpts.Marshal(store)
		if err != nil {
			t.Fatalf("mOpts.Marshal failed: %v", err)
		}

		// Verify it's compressed
		if !uOpts.IsCompressed(data) {
			t.Errorf("Data not compressed with full Marshal method")
		}

		// Unmarshal with our wrapper
		newStore2 := &v1.Store{}
		err = protozstd.Unmarshal(data, newStore2)
		if err != nil {
			t.Fatalf("protozstd.Unmarshal failed: %v", err)
		}

		// Verify message equality
		if !proto.Equal(store, newStore2) {
			t.Errorf("Messages don't match after full Marshal/Unmarshal cycle")
		}
	})

	// Test with a large message to ensure compression happens
	t.Run("Large message", func(t *testing.T) {
		// Create a large message by adding lots of employee IDs
		largeStore := proto.Clone(store).(*v1.Store)
		arr := largeStore.GetBestEmployeeIds()
		for i := 0; i < 1000; i++ {
			arr = append(arr, "employee"+strconv.Itoa(i))
		}
		largeStore.SetBestEmployeeIds(arr)

		// Marshal with our wrapper
		data, err := mOpts.Marshal(largeStore)
		if err != nil {
			t.Fatalf("mOpts.Marshal failed for large message: %v", err)
		}

		// Verify it's compressed
		if !uOpts.IsCompressed(data) {
			t.Errorf("Large message data not compressed")
		}

		// Calculate compression ratio
		protoBytes, _ := proto.Marshal(largeStore)
		compressionRatio := float64(len(data)) / float64(len(protoBytes))
		t.Logf("Compression ratio: %.2f (compressed: %d bytes, original: %d bytes)",
			compressionRatio, len(data), len(protoBytes))

		// Unmarshal with our wrapper
		newLargeStore := &v1.Store{}
		err = protozstd.Unmarshal(data, newLargeStore)
		if err != nil {
			t.Fatalf("protozstd.Unmarshal failed for large message: %v", err)
		}

		// Verify message equality
		if !proto.Equal(largeStore, newLargeStore) {
			t.Errorf("Large messages don't match after Marshal/Unmarshal cycle")
		}
	})
}

// Benchmark with real protobuf messages
func BenchmarkProtoMarshalUnmarshal(b *testing.B) {
	// Create a Store message similar to our test
	now := time.Now()
	store := v1.Store_builder{
		Id:      proto.String("store123"),
		Country: proto.String("USA"),
		Region:  proto.String("West"),
		State:   proto.String("California"),
		City:    proto.String("San Francisco"),
		Closed:  proto.Bool(false),
		OpeningDate: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		BestEmployeeIds: []string{"emp1", "emp2", "emp3"},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Foo:     proto.Uint64(12345),
		Morefoo: []uint64{1, 2, 3, 4, 5},
	}.Build()

	// Create options
	mOpts := protozstd.NewMarshalOptions()
	uOpts := protozstd.NewUnmarshalOptions()

	b.Run("Small Message", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data, err := mOpts.Marshal(store)
			if err != nil {
				b.Fatalf("Marshal failed: %v", err)
			}

			newStore := &v1.Store{}
			err = uOpts.Unmarshal(data, newStore)
			if err != nil {
				b.Fatalf("Unmarshal failed: %v", err)
			}
		}
	})

	// Create a large message
	largeStore := proto.Clone(store).(*v1.Store)
	arr := largeStore.GetBestEmployeeIds()
	for i := 0; i < 1000; i++ {
		arr = append(arr, "employee"+strconv.Itoa(i))
	}
	largeStore.SetBestEmployeeIds(arr)

	b.Run("Large Message", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data, err := mOpts.Marshal(largeStore)
			if err != nil {
				b.Fatalf("Marshal failed: %v", err)
			}

			newStore := &v1.Store{}
			err = uOpts.Unmarshal(data, newStore)
			if err != nil {
				b.Fatalf("Unmarshal failed: %v", err)
			}
		}
	})

	// Compare with standard proto marshaling for large message
	b.Run("Standard Proto (Large)", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data, err := proto.Marshal(largeStore)
			if err != nil {
				b.Fatalf("Proto Marshal failed: %v", err)
			}

			newStore := &v1.Store{}
			err = proto.Unmarshal(data, newStore)
			if err != nil {
				b.Fatalf("Proto Unmarshal failed: %v", err)
			}
		}
	})
}
