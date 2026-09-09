package pgd

import (
	"strings"
	"testing"

	dynamopb "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
)

func TestValidateShardConfig(t *testing.T) {
	const msgName = "test.v1.Widget"

	cases := []struct {
		name    string
		key     *dynamopb.Key
		wantErr string
	}{
		{
			name: "nil key",
			key:  nil,
		},
		{
			name: "no shard config is not sharded",
			key:  &dynamopb.Key{PkFields: []string{"tenant_id"}, SkConst: "widget"},
		},
		{
			name: "valid sharded key",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				SkFields: []string{"id"},
				Shard:    &dynamopb.ShardOptions{ShardCount: 32},
			},
		},
		{
			name: "shard_count below minimum",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				SkFields: []string{"id"},
				Shard:    &dynamopb.ShardOptions{ShardCount: 1},
			},
			wantErr: "shard_count must be >=",
		},
		{
			name: "shard_count above maximum",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				SkFields: []string{"id"},
				Shard:    &dynamopb.ShardOptions{ShardCount: 2000},
			},
			wantErr: "shard_count must be <=",
		},
		{
			name: "shard_count not a power of two",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				SkFields: []string{"id"},
				Shard:    &dynamopb.ShardOptions{ShardCount: 5},
			},
			wantErr: "shard_count must be a power of 2",
		},
		{
			name: "sk_const cannot be sharded",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				SkConst:  "widget",
				Shard:    &dynamopb.ShardOptions{ShardCount: 32},
			},
			wantErr: "sharded key cannot use sk_const",
		},
		{
			name: "sk_const wins over sk_fields when both are set",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				SkFields: []string{"id"},
				SkConst:  "widget",
				Shard:    &dynamopb.ShardOptions{ShardCount: 32},
			},
			wantErr: "sharded key cannot use sk_const",
		},
		{
			name: "sharded key without any sort key",
			key: &dynamopb.Key{
				PkFields: []string{"tenant_id"},
				Shard:    &dynamopb.ShardOptions{ShardCount: 32},
			},
			wantErr: "sharded key must have sk_fields configured",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateShardConfig(msgName, tc.key)
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
