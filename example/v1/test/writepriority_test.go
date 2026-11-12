package v1_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	dynamov1 "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
	v1 "github.com/pquerna/protoc-gen-dynamo/example/v1"
	"github.com/pquerna/protoc-gen-dynamo/pkg/dynamoattr"
)

func newBulkRecord() *v1.BulkRecord {
	return v1.BulkRecord_builder{
		TenantId: proto.String("acme"),
		Id:       proto.String("rec-1"),
		Payload:  proto.String("hello"),
	}.Build()
}

func TestWritePriorityMarshalled(t *testing.T) {
	user := v1.UserV2_builder{
		TenantId:  proto.String("acme"),
		Id:        proto.String("user-1"),
		IdpId:     proto.String("okta-42"),
		Email:     proto.String("someone@example.com"),
		UpdatedAt: timestamppb.New(time.Unix(1700000001, 0)),
	}.Build()
	bulk := newBulkRecord()

	cases := []struct {
		name     string
		marshal  func() (map[string]types.AttributeValue, error)
		accessor dynamov1.WritePriority
		static   dynamov1.WritePriority
		want     string
	}{
		{
			name:     "default message is normal",
			marshal:  user.MarshalDynamoItem,
			accessor: user.DynamoDefaultWritePriority(),
			static:   v1.UserV2DefaultWritePriority(),
			want:     "normal",
		},
		{
			name:     "background message is background",
			marshal:  bulk.MarshalDynamoItem,
			accessor: bulk.DynamoDefaultWritePriority(),
			static:   v1.BulkRecordDefaultWritePriority(),
			want:     "background",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item, err := tc.marshal()
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			got, ok := item[dynamoattr.WritePriorityAttr].(*types.AttributeValueMemberS)
			if !ok {
				t.Fatalf("item[%q]: expected *types.AttributeValueMemberS, got %T", dynamoattr.WritePriorityAttr, item[dynamoattr.WritePriorityAttr])
			}
			if got.Value != tc.want {
				t.Fatalf("item[%q] = %q, want %q", dynamoattr.WritePriorityAttr, got.Value, tc.want)
			}
			if tc.accessor == dynamov1.WritePriority_WRITE_PRIORITY_UNSPECIFIED {
				t.Fatal("DynamoDefaultWritePriority returned UNSPECIFIED")
			}
			if tc.accessor != tc.static {
				t.Fatalf("instance accessor %s != static accessor %s", tc.accessor, tc.static)
			}
			if s := dynamoattr.WritePriorityAttrValue(tc.accessor); s != got.Value {
				t.Fatalf("accessor %s stores as %q, item has %q", tc.accessor, s, got.Value)
			}
		})
	}
}

func TestWritePriorityKeylessOmitted(t *testing.T) {
	ref := v1.StoreRef_builder{Id: proto.String("store-1")}.Build()
	item, err := ref.MarshalDynamoItem()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if _, ok := item[dynamoattr.WritePriorityAttr]; ok {
		t.Fatalf("keyless message should not carry %q", dynamoattr.WritePriorityAttr)
	}
}

// The priority describes how a write is processed, not the object, so it must
// never affect what unmarshals.
func TestWritePriorityIgnoredOnUnmarshal(t *testing.T) {
	orig := newBulkRecord()

	cases := map[string]func(map[string]types.AttributeValue){
		"missing": func(item map[string]types.AttributeValue) {
			delete(item, dynamoattr.WritePriorityAttr)
		},
		"unknown value": func(item map[string]types.AttributeValue) {
			item[dynamoattr.WritePriorityAttr] = &types.AttributeValueMemberS{Value: "bogus"}
		},
		"overridden": func(item map[string]types.AttributeValue) {
			dynamoattr.SetWritePriority(item, dynamov1.WritePriority_WRITE_PRIORITY_NORMAL)
		},
	}

	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			item, err := orig.MarshalDynamoItem()
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			mutate(item)

			got := &v1.BulkRecord{}
			if err := got.UnmarshalDynamoItem(item); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if !proto.Equal(orig, got) {
				t.Fatalf("round trip mismatch: got %v, want %v", got, orig)
			}
		})
	}
}

func TestWritePriorityOverride(t *testing.T) {
	item, err := newBulkRecord().MarshalDynamoItem()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if p, ok := dynamoattr.GetWritePriority(item); !ok || p != dynamov1.WritePriority_WRITE_PRIORITY_BACKGROUND {
		t.Fatalf("default: got %s, %v; want BACKGROUND", p, ok)
	}
	dynamoattr.SetWritePriority(item, dynamov1.WritePriority_WRITE_PRIORITY_NORMAL)
	if p, ok := dynamoattr.GetWritePriority(item); !ok || p != dynamov1.WritePriority_WRITE_PRIORITY_NORMAL {
		t.Fatalf("override: got %s, %v; want NORMAL", p, ok)
	}
}
