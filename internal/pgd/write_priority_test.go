package pgd

import (
	"strings"
	"testing"

	dynamopb "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
	"github.com/pquerna/protoc-gen-dynamo/pkg/dynamoattr"
)

func defaultPriority(p dynamopb.WritePriority) *dynamopb.DynamoMessageOptions_WritePriority {
	return &dynamopb.DynamoMessageOptions_WritePriority{WritePriority: p}
}

func pinnedPriority(p dynamopb.WritePriority) *dynamopb.DynamoMessageOptions_PinnedWritePriority {
	return &dynamopb.DynamoMessageOptions_PinnedWritePriority{PinnedWritePriority: p}
}

func TestWritePriorityAttrValue(t *testing.T) {
	cases := []struct {
		name       string
		mext       *dynamopb.DynamoMessageOptions
		want       string
		wantPinned bool
	}{
		{name: "neither set resolves to normal", mext: &dynamopb.DynamoMessageOptions{}, want: "normal"},
		{name: "unspecified resolves to normal", mext: &dynamopb.DynamoMessageOptions{Priority: defaultPriority(dynamopb.WritePriority_WRITE_PRIORITY_UNSPECIFIED)}, want: "normal"},
		{name: "normal", mext: &dynamopb.DynamoMessageOptions{Priority: defaultPriority(dynamopb.WritePriority_WRITE_PRIORITY_NORMAL)}, want: "normal"},
		{name: "background", mext: &dynamopb.DynamoMessageOptions{Priority: defaultPriority(dynamopb.WritePriority_WRITE_PRIORITY_BACKGROUND)}, want: "background"},
		{name: "pinned normal", mext: &dynamopb.DynamoMessageOptions{Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_NORMAL)}, want: "normal", wantPinned: true},
		{name: "pinned background", mext: &dynamopb.DynamoMessageOptions{Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_BACKGROUND)}, want: "background", wantPinned: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := writePriorityAttrValue(tc.mext)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
			if pinned := writePriorityPinned(tc.mext); pinned != tc.wantPinned {
				t.Fatalf("pinned = %v, want %v", pinned, tc.wantPinned)
			}
		})
	}

	if _, err := writePriorityAttrValue(&dynamopb.DynamoMessageOptions{Priority: defaultPriority(999)}); err == nil {
		t.Fatal("expected an error for an unknown write_priority")
	}
	if _, err := writePriorityAttrValue(&dynamopb.DynamoMessageOptions{Priority: pinnedPriority(999)}); err == nil {
		t.Fatal("expected an error for an unknown pinned_write_priority")
	}
}

func TestValidateAttributeNames(t *testing.T) {
	const msgName = "test.v1.Widget"

	twoKeys := &dynamopb.DynamoMessageOptions{Key: []*dynamopb.Key{
		{PkFields: []string{"tenant_id"}, SkFields: []string{"id"}},
		{PkFields: []string{"tenant_id"}, SkFields: []string{"email"}},
	}}
	keyless := &dynamopb.DynamoMessageOptions{}

	cases := []struct {
		name       string
		mext       *dynamopb.DynamoMessageOptions
		hasVersion bool
		hasDeleted bool
		exposed    map[string]string
		wantErr    string
	}{
		{
			name: "no exposed fields",
			mext: twoKeys,
		},
		{
			name:    "non-colliding exposed field",
			mext:    twoKeys,
			exposed: map[string]string{"expires_at": "test.v1.Widget.expires_at"},
		},
		{
			name:    "collides with pk",
			mext:    twoKeys,
			exposed: map[string]string{"pk": "test.v1.Widget.pk"},
			wantErr: `"pk"`,
		},
		{
			name:    "collides with gsi1sk",
			mext:    twoKeys,
			exposed: map[string]string{"gsi1sk": "test.v1.Widget.gsi1sk"},
			wantErr: `"gsi1sk"`,
		},
		{
			name:    "collides with typ",
			mext:    twoKeys,
			exposed: map[string]string{"typ": "test.v1.Widget.typ"},
			wantErr: `"typ"`,
		},
		{
			name:    "collides with value",
			mext:    keyless,
			exposed: map[string]string{"value": "test.v1.Widget.value"},
			wantErr: `"value"`,
		},
		{
			name:    "collides with wp on a keyed message",
			mext:    twoKeys,
			exposed: map[string]string{"wp": "test.v1.Widget.wp"},
			wantErr: `"wp"`,
		},
		{
			// The map key is the rendered attribute name, so a field renamed via
			// (dynamo.v1.field).name collides under its new name.
			name:    "collides through (dynamo.v1.field).name",
			mext:    twoKeys,
			exposed: map[string]string{"sk": "test.v1.Widget.sort_key"},
			wantErr: "test.v1.Widget.sort_key",
		},
		{
			name:    "wp is free on a keyless message",
			mext:    keyless,
			exposed: map[string]string{"wp": "test.v1.Widget.wp"},
		},
		{
			name:    "gsi2pk is free with only one gsi",
			mext:    twoKeys,
			exposed: map[string]string{"gsi2pk": "test.v1.Widget.gsi2pk"},
		},
		{
			name:    "version is free without updated_at",
			mext:    twoKeys,
			exposed: map[string]string{"version": "test.v1.Widget.version"},
		},
		{
			name:       "version collides with updated_at",
			mext:       twoKeys,
			hasVersion: true,
			exposed:    map[string]string{"version": "test.v1.Widget.version"},
			wantErr:    `"version"`,
		},
		{
			name:    "deleted is free without deleted_at",
			mext:    twoKeys,
			exposed: map[string]string{"deleted": "test.v1.Widget.deleted"},
		},
		{
			name:       "deleted collides with deleted_at",
			mext:       twoKeys,
			hasDeleted: true,
			exposed:    map[string]string{"deleted": "test.v1.Widget.deleted"},
			wantErr:    `"deleted"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			emitted := emittedAttributeNames(tc.mext, tc.hasVersion, tc.hasDeleted)
			err := validateAttributeNames(msgName, emitted, tc.exposed)
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("expected no error, got %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			case tc.wantErr != "":
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}
				if !strings.Contains(err.Error(), msgName) {
					t.Fatalf("error should name the message %q, got %v", msgName, err)
				}
			}
		})
	}
}

// The generator and pkg/dynamoattr each define the attribute name, to keep the
// generator free of runtime imports. They must agree.
func TestWritePriorityFieldMatchesRuntime(t *testing.T) {
	if writePriorityField != dynamoattr.WritePriorityAttr {
		t.Fatalf("generator writes %q, dynamoattr reads %q", writePriorityField, dynamoattr.WritePriorityAttr)
	}
}

func TestValidateWritePriority(t *testing.T) {
	const msgName = "test.v1.Widget"
	keyed := []*dynamopb.Key{{PkFields: []string{"tenant_id"}, SkFields: []string{"id"}}}

	cases := []struct {
		name    string
		mext    *dynamopb.DynamoMessageOptions
		wantErr string
	}{
		{name: "keyless, unset", mext: &dynamopb.DynamoMessageOptions{}},
		{name: "keyed, unset", mext: &dynamopb.DynamoMessageOptions{Key: keyed}},
		{name: "keyed, background", mext: &dynamopb.DynamoMessageOptions{Key: keyed, Priority: defaultPriority(dynamopb.WritePriority_WRITE_PRIORITY_BACKGROUND)}},
		{name: "keyless, background", mext: &dynamopb.DynamoMessageOptions{Priority: defaultPriority(dynamopb.WritePriority_WRITE_PRIORITY_BACKGROUND)}, wantErr: "write_priority is set"},
		{name: "keyless, normal", mext: &dynamopb.DynamoMessageOptions{Priority: defaultPriority(dynamopb.WritePriority_WRITE_PRIORITY_NORMAL)}, wantErr: "write_priority is set"},
		{name: "keyed, pinned background", mext: &dynamopb.DynamoMessageOptions{Key: keyed, Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_BACKGROUND)}},
		{name: "keyed, pinned normal", mext: &dynamopb.DynamoMessageOptions{Key: keyed, Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_NORMAL)}},
		{name: "keyed, pinned unspecified", mext: &dynamopb.DynamoMessageOptions{Key: keyed, Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_UNSPECIFIED)}, wantErr: "pinned_write_priority is WRITE_PRIORITY_UNSPECIFIED"},
		{name: "keyless, pinned background", mext: &dynamopb.DynamoMessageOptions{Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_BACKGROUND)}, wantErr: "pinned_write_priority is set"},
		{name: "keyless, pinned unspecified", mext: &dynamopb.DynamoMessageOptions{Priority: pinnedPriority(dynamopb.WritePriority_WRITE_PRIORITY_UNSPECIFIED)}, wantErr: "pinned_write_priority is set"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateWritePriority(msgName, tc.mext)
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("expected no error, got %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			case tc.wantErr != "":
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}
				if !strings.Contains(err.Error(), msgName) {
					t.Fatalf("error should name the message %q, got %v", msgName, err)
				}
			}
		})
	}
}
