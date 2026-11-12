package dynamoattr

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"

	dynamov1 "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
)

// Every WritePriority value except UNSPECIFIED must pin a distinct stored
// string. This guards future enum edits.
func TestWritePriorityEnumComplete(t *testing.T) {
	seen := map[string]string{}
	values := dynamov1.WritePriority(0).Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		ev := values.Get(i)
		s, _ := proto.GetExtension(ev.Options(), dynamov1.E_AttrValue).(string)
		if ev.Number() == 0 {
			if s != "" {
				t.Errorf("%s must not have an attr_value, got %q", ev.Name(), s)
			}
			continue
		}
		if s == "" {
			t.Errorf("%s has no (dynamo.v1.attr_value)", ev.Name())
			continue
		}
		if prev, dup := seen[s]; dup {
			t.Errorf("%s and %s share attr_value %q", prev, ev.Name(), s)
		}
		seen[s] = string(ev.Name())
	}
}

func TestWritePriorityRoundTrip(t *testing.T) {
	values := dynamov1.WritePriority(0).Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		p := dynamov1.WritePriority(values.Get(i).Number())
		want := p
		if p == dynamov1.WritePriority_WRITE_PRIORITY_UNSPECIFIED {
			want = dynamov1.WritePriority_WRITE_PRIORITY_NORMAL
		}
		got, ok := ParseWritePriority(WritePriorityAttrValue(p))
		if !ok || got != want {
			t.Errorf("round trip %s: got %s, %v; want %s", p, got, ok, want)
		}
	}
}

func TestWritePriorityAttrValues(t *testing.T) {
	cases := map[dynamov1.WritePriority]string{
		dynamov1.WritePriority_WRITE_PRIORITY_UNSPECIFIED: "normal",
		dynamov1.WritePriority_WRITE_PRIORITY_NORMAL:      "normal",
		dynamov1.WritePriority_WRITE_PRIORITY_BACKGROUND:  "background",
		dynamov1.WritePriority(999):                       "normal",
	}
	for p, want := range cases {
		if got := WritePriorityAttrValue(p); got != want {
			t.Errorf("WritePriorityAttrValue(%s) = %q, want %q", p, got, want)
		}
	}
}

func TestParseWritePriorityUnknown(t *testing.T) {
	for _, s := range []string{"", "bogus", "NORMAL"} {
		if p, ok := ParseWritePriority(s); ok {
			t.Errorf("ParseWritePriority(%q) = %s, true; want ok=false", s, p)
		}
	}
}

func TestSetGetWritePriority(t *testing.T) {
	item := map[string]types.AttributeValue{}
	if _, ok := GetWritePriority(item); ok {
		t.Fatal("GetWritePriority on empty item: want ok=false")
	}

	SetWritePriority(item, dynamov1.WritePriority_WRITE_PRIORITY_BACKGROUND)
	if p, ok := GetWritePriority(item); !ok || p != dynamov1.WritePriority_WRITE_PRIORITY_BACKGROUND {
		t.Fatalf("after set BACKGROUND: got %s, %v", p, ok)
	}

	SetWritePriority(item, dynamov1.WritePriority_WRITE_PRIORITY_NORMAL)
	if p, ok := GetWritePriority(item); !ok || p != dynamov1.WritePriority_WRITE_PRIORITY_NORMAL {
		t.Fatalf("after overwrite NORMAL: got %s, %v", p, ok)
	}
	if s := item[WritePriorityAttr].(*types.AttributeValueMemberS).Value; s != "normal" {
		t.Fatalf("stored value = %q, want %q", s, "normal")
	}

	item[WritePriorityAttr] = &types.AttributeValueMemberN{Value: "1"}
	if _, ok := GetWritePriority(item); ok {
		t.Fatal("GetWritePriority on non-S attribute: want ok=false")
	}

	item[WritePriorityAttr] = &types.AttributeValueMemberS{Value: "bogus"}
	if _, ok := GetWritePriority(item); ok {
		t.Fatal("GetWritePriority on unknown value: want ok=false")
	}
}
