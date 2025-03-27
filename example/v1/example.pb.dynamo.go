// Code generated by protoc-gen-dynamo v0.1.0. DO NOT EDIT.
// source: v1/example.proto

package v1

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pquerna/protoc-gen-dynamo/pkg/protozstd"
	"strconv"
	"strings"
	"time"
)

func (p *Store) MarshalDynamo() (types.AttributeValue, error) {
	return p.MarshalDynamoDBAttributeValue()
}

func (p *Store) MarshalDynamoItem() (map[string]types.AttributeValue, error) {
	av, err := p.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}
	avm, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("unable to marshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	return avm.Value, nil
}

func (p *Store) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var sb strings.Builder
	var err error
	nullBoolTrue := true
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	v1 := &types.AttributeValueMemberS{Value: sb.String()}
	v2 := &types.AttributeValueMemberS{Value: "example"}
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	v3 := &types.AttributeValueMemberS{Value: sb.String()}
	v4 := &types.AttributeValueMemberS{Value: "dummyvalue"}
	v5, err := p.Version()
	if err != nil {
		return nil, err
	}
	v6buf, err := protozstd.Marshal(p)
	if err != nil {
		return nil, err
	}
	v6 := &types.AttributeValueMemberB{Value: v6buf}
	v7 := &types.AttributeValueMemberS{Value: "examplepb.v1.Store"}
	var v8 types.AttributeValue
	if len(p.Id) != 0 {
		v8 = &types.AttributeValueMemberS{Value: p.Id}
	} else {
		v8 = &types.AttributeValueMemberNULL{Value: nullBoolTrue}
	}
	var av types.AttributeValue
	av = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"expires_at":    &types.AttributeValueMemberN{Value: (strconv.FormatInt(int64(p.ExpiresAt.AsTime().Round(time.Second).Unix()), 10))},
		"expires_at_ms": &types.AttributeValueMemberN{Value: (strconv.FormatInt(int64(p.ExpiresAtMs.AsTime().Round(time.Millisecond).UnixNano()/int64(time.Millisecond)), 10))},
		"expires_at_ns": &types.AttributeValueMemberN{Value: (strconv.FormatInt(int64(p.ExpiresAtNs.AsTime().UnixNano()), 10))},
		"gsi1pk":        v3,
		"gsi1sk":        v4,
		"pk":            v1,
		"sk":            v2,
		"store_id":      v8,
		"typ":           v7,
		"value":         v6,
		"version":       &types.AttributeValueMemberN{Value: (strconv.FormatInt(v5, 10))},
	}}
	return av, nil
}

func (p *User) MarshalDynamo() (types.AttributeValue, error) {
	return p.MarshalDynamoDBAttributeValue()
}

func (p *User) MarshalDynamoItem() (map[string]types.AttributeValue, error) {
	av, err := p.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}
	avm, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("unable to marshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	return avm.Value, nil
}

func (p *User) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var sb strings.Builder
	var err error
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user:")
	_, _ = sb.WriteString(p.TenantId)
	v1 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString(p.Id)
	v2 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user:")
	_, _ = sb.WriteString(p.TenantId)
	v3 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	v4 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user:")
	_, _ = sb.WriteString(p.TenantId)
	v5 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatInt(int64(p.AnEnum), 10))
	v6 := &types.AttributeValueMemberS{Value: sb.String()}
	v7, err := p.Version()
	if err != nil {
		return nil, err
	}
	v8buf, err := protozstd.Marshal(p)
	if err != nil {
		return nil, err
	}
	v8 := &types.AttributeValueMemberB{Value: v8buf}
	v9 := &types.AttributeValueMemberS{Value: "examplepb.v1.User"}
	v10 := &types.AttributeValueMemberBOOL{Value: p.DeletedAt.IsValid()}
	var av types.AttributeValue
	av = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"deleted": v10,
		"gsi1pk":  v3,
		"gsi1sk":  v4,
		"gsi2pk":  v5,
		"gsi2sk":  v6,
		"pk":      v1,
		"sk":      v2,
		"typ":     v9,
		"value":   v8,
		"version": &types.AttributeValueMemberN{Value: (strconv.FormatInt(v7, 10))},
	}}
	return av, nil
}

func (p *StoreV2) MarshalDynamo() (types.AttributeValue, error) {
	return p.MarshalDynamoDBAttributeValue()
}

func (p *StoreV2) MarshalDynamoItem() (map[string]types.AttributeValue, error) {
	av, err := p.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}
	avm, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("unable to marshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	return avm.Value, nil
}

func (p *StoreV2) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var sb strings.Builder
	var err error
	nullBoolTrue := true
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store_v_2:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	v1 := &types.AttributeValueMemberS{Value: sb.String()}
	v2 := &types.AttributeValueMemberS{Value: "example"}
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store_v_2:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	v3 := &types.AttributeValueMemberS{Value: sb.String()}
	v4 := &types.AttributeValueMemberS{Value: "dummyvalue"}
	v5, err := p.Version()
	if err != nil {
		return nil, err
	}
	v6buf, err := protozstd.Marshal(p)
	if err != nil {
		return nil, err
	}
	v6 := &types.AttributeValueMemberB{Value: v6buf}
	v7 := &types.AttributeValueMemberS{Value: "examplepb.v1.StoreV2"}
	var v8 types.AttributeValue
	if len(p.Id) != 0 {
		v8 = &types.AttributeValueMemberS{Value: p.Id}
	} else {
		v8 = &types.AttributeValueMemberNULL{Value: nullBoolTrue}
	}
	var av types.AttributeValue
	av = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"expires_at":    &types.AttributeValueMemberN{Value: (strconv.FormatInt(int64(p.ExpiresAt.AsTime().Round(time.Second).Unix()), 10))},
		"expires_at_ms": &types.AttributeValueMemberN{Value: (strconv.FormatInt(int64(p.ExpiresAtMs.AsTime().Round(time.Millisecond).UnixNano()/int64(time.Millisecond)), 10))},
		"expires_at_ns": &types.AttributeValueMemberN{Value: (strconv.FormatInt(int64(p.ExpiresAtNs.AsTime().UnixNano()), 10))},
		"gsi1pk":        v3,
		"gsi1sk":        v4,
		"pk":            v1,
		"sk":            v2,
		"store_id":      v8,
		"typ":           v7,
		"value":         v6,
		"version":       &types.AttributeValueMemberN{Value: (strconv.FormatInt(v5, 10))},
	}}
	return av, nil
}

func (p *UserV2) MarshalDynamo() (types.AttributeValue, error) {
	return p.MarshalDynamoDBAttributeValue()
}

func (p *UserV2) MarshalDynamoItem() (map[string]types.AttributeValue, error) {
	av, err := p.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}
	avm, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("unable to marshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	return avm.Value, nil
}

func (p *UserV2) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var sb strings.Builder
	var err error
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user_v_2:")
	_, _ = sb.WriteString(p.TenantId)
	v1 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString(p.Id)
	v2 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user_v_2:")
	_, _ = sb.WriteString(p.TenantId)
	v3 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	v4 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user_v_2:")
	_, _ = sb.WriteString(p.TenantId)
	v5 := &types.AttributeValueMemberS{Value: sb.String()}
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatInt(int64(p.AnEnum), 10))
	v6 := &types.AttributeValueMemberS{Value: sb.String()}
	v7, err := p.Version()
	if err != nil {
		return nil, err
	}
	v8buf, err := protozstd.Marshal(p)
	if err != nil {
		return nil, err
	}
	v8 := &types.AttributeValueMemberB{Value: v8buf}
	v9 := &types.AttributeValueMemberS{Value: "examplepb.v1.UserV2"}
	v10 := &types.AttributeValueMemberBOOL{Value: p.DeletedAt.IsValid()}
	var av types.AttributeValue
	av = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"deleted": v10,
		"gsi1pk":  v3,
		"gsi1sk":  v4,
		"gsi2pk":  v5,
		"gsi2sk":  v6,
		"pk":      v1,
		"sk":      v2,
		"typ":     v9,
		"value":   v8,
		"version": &types.AttributeValueMemberN{Value: (strconv.FormatInt(v7, 10))},
	}}
	return av, nil
}

func (p *Store) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	typ, ok := m.Value["typ"]
	if !ok {
		return fmt.Errorf("dynamo: typ missing")
	}
	t, ok := typ.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberS, got %T", typ)
	}
	if t.Value != "examplepb.v1.Store" {
		return fmt.Errorf("dynamo: _type mismatch: examplepb.v1.Store expected, got: '%s'", typ)
	}
	value, ok := m.Value["value"]
	if !ok {
		return fmt.Errorf("dynamo: value missing")
	}
	v, ok := value.(*types.AttributeValueMemberB)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberB, got %T", value)
	}
	var data []byte
	data = v.Value
	return protozstd.Unmarshal(data, p)
}

func (p *Store) UnmarshalDynamo(av types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(av)
}

func (p *Store) UnmarshalDynamoItem(av map[string]types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberM{Value: av})
}

func (p *User) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	typ, ok := m.Value["typ"]
	if !ok {
		return fmt.Errorf("dynamo: typ missing")
	}
	t, ok := typ.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberS, got %T", typ)
	}
	if t.Value != "examplepb.v1.User" {
		return fmt.Errorf("dynamo: _type mismatch: examplepb.v1.User expected, got: '%s'", typ)
	}
	value, ok := m.Value["value"]
	if !ok {
		return fmt.Errorf("dynamo: value missing")
	}
	v, ok := value.(*types.AttributeValueMemberB)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberB, got %T", value)
	}
	var data []byte
	data = v.Value
	return protozstd.Unmarshal(data, p)
}

func (p *User) UnmarshalDynamo(av types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(av)
}

func (p *User) UnmarshalDynamoItem(av map[string]types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberM{Value: av})
}

func (p *StoreV2) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	typ, ok := m.Value["typ"]
	if !ok {
		return fmt.Errorf("dynamo: typ missing")
	}
	t, ok := typ.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberS, got %T", typ)
	}
	if t.Value != "examplepb.v1.StoreV2" {
		return fmt.Errorf("dynamo: _type mismatch: examplepb.v1.StoreV2 expected, got: '%s'", typ)
	}
	value, ok := m.Value["value"]
	if !ok {
		return fmt.Errorf("dynamo: value missing")
	}
	v, ok := value.(*types.AttributeValueMemberB)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberB, got %T", value)
	}
	var data []byte
	data = v.Value
	return protozstd.Unmarshal(data, p)
}

func (p *StoreV2) UnmarshalDynamo(av types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(av)
}

func (p *StoreV2) UnmarshalDynamoItem(av map[string]types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberM{Value: av})
}

func (p *UserV2) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberM, got %T", av)
	}
	typ, ok := m.Value["typ"]
	if !ok {
		return fmt.Errorf("dynamo: typ missing")
	}
	t, ok := typ.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberS, got %T", typ)
	}
	if t.Value != "examplepb.v1.UserV2" {
		return fmt.Errorf("dynamo: _type mismatch: examplepb.v1.UserV2 expected, got: '%s'", typ)
	}
	value, ok := m.Value["value"]
	if !ok {
		return fmt.Errorf("dynamo: value missing")
	}
	v, ok := value.(*types.AttributeValueMemberB)
	if !ok {
		return fmt.Errorf("unable to unmarshal: expected type *types.AttributeValueMemberB, got %T", value)
	}
	var data []byte
	data = v.Value
	return protozstd.Unmarshal(data, p)
}

func (p *UserV2) UnmarshalDynamo(av types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(av)
}

func (p *UserV2) UnmarshalDynamoItem(av map[string]types.AttributeValue) error {
	return p.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberM{Value: av})
}

func (p *Store) Version() (int64, error) {
	err := p.UpdatedAt.CheckValid()
	if err != nil {
		return 0, err
	}
	t := p.UpdatedAt.AsTime()
	return t.UnixNano(), nil
}

func (p *Store) PartitionKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	return sb.String()
}

func StorePartitionKey(id string, country string, foo uint64) string {
	return (&Store{
		Country: country,
		Foo:     foo,
		Id:      id,
	}).PartitionKey()
}

func (p *Store) SortKey() string {
	return "example"
}

func StoreSortKey() string {
	return (&Store{}).SortKey()
}

func (p *Store) Gsi1PkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	return sb.String()
}

func StoreGsi1PkKey(id string, country string, foo uint64) string {
	return (&Store{
		Country: country,
		Foo:     foo,
		Id:      id,
	}).Gsi1PkKey()
}

func (p *Store) Gsi1SkKey() string {
	return "dummyvalue"
}

func StoreGsi1SkKey() string {
	return (&Store{}).Gsi1SkKey()
}

func (p *User) Version() (int64, error) {
	err := p.UpdatedAt.CheckValid()
	if err != nil {
		return 0, err
	}
	t := p.UpdatedAt.AsTime()
	return t.UnixNano(), nil
}

func (p *User) PartitionKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user:")
	_, _ = sb.WriteString(p.TenantId)
	return sb.String()
}

func UserPartitionKey(tenantId string) string {
	return (&User{TenantId: tenantId}).PartitionKey()
}

func (p *User) SortKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString(p.Id)
	return sb.String()
}

func UserSortKey(id string) string {
	return (&User{Id: id}).SortKey()
}

func (p *User) Gsi1PkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user:")
	_, _ = sb.WriteString(p.TenantId)
	return sb.String()
}

func UserGsi1PkKey(tenantId string) string {
	return (&User{TenantId: tenantId}).Gsi1PkKey()
}

func (p *User) Gsi1SkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	return sb.String()
}

func UserGsi1SkKey(idpId string) string {
	return (&User{IdpId: idpId}).Gsi1SkKey()
}

func (p *User) Gsi2PkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user:")
	_, _ = sb.WriteString(p.TenantId)
	return sb.String()
}

func UserGsi2PkKey(tenantId string) string {
	return (&User{TenantId: tenantId}).Gsi2PkKey()
}

func (p *User) Gsi2SkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatInt(int64(p.AnEnum), 10))
	return sb.String()
}

func UserGsi2SkKey(idpId string, anEnum BasicEnum) string {
	return (&User{
		AnEnum: anEnum,
		IdpId:  idpId,
	}).Gsi2SkKey()
}

func (p *StoreV2) Version() (int64, error) {
	err := p.UpdatedAt.CheckValid()
	if err != nil {
		return 0, err
	}
	t := p.UpdatedAt.AsTime()
	return t.UnixNano(), nil
}

func (p *StoreV2) PartitionKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store_v_2:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	return sb.String()
}

func StoreV2PartitionKey(id string, country string, foo uint64) string {
	return (&StoreV2{
		Country: country,
		Foo:     foo,
		Id:      id,
	}).PartitionKey()
}

func (p *StoreV2) SortKey() string {
	return "example"
}

func StoreV2SortKey() string {
	return (&StoreV2{}).SortKey()
}

func (p *StoreV2) Gsi1PkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_store_v_2:")
	_, _ = sb.WriteString(p.Id)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p.Country)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatUint(uint64(p.Foo), 10))
	return sb.String()
}

func StoreV2Gsi1PkKey(id string, country string, foo uint64) string {
	return (&StoreV2{
		Country: country,
		Foo:     foo,
		Id:      id,
	}).Gsi1PkKey()
}

func (p *StoreV2) Gsi1SkKey() string {
	return "dummyvalue"
}

func StoreV2Gsi1SkKey() string {
	return (&StoreV2{}).Gsi1SkKey()
}

func (p *UserV2) Version() (int64, error) {
	err := p.UpdatedAt.CheckValid()
	if err != nil {
		return 0, err
	}
	t := p.UpdatedAt.AsTime()
	return t.UnixNano(), nil
}

func (p *UserV2) PartitionKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user_v_2:")
	_, _ = sb.WriteString(p.TenantId)
	return sb.String()
}

func UserV2PartitionKey(tenantId string) string {
	return (&UserV2{TenantId: tenantId}).PartitionKey()
}

func (p *UserV2) SortKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString(p.Id)
	return sb.String()
}

func UserV2SortKey(id string) string {
	return (&UserV2{Id: id}).SortKey()
}

func (p *UserV2) Gsi1PkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user_v_2:")
	_, _ = sb.WriteString(p.TenantId)
	return sb.String()
}

func UserV2Gsi1PkKey(tenantId string) string {
	return (&UserV2{TenantId: tenantId}).Gsi1PkKey()
}

func (p *UserV2) Gsi1SkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	return sb.String()
}

func UserV2Gsi1SkKey(idpId string) string {
	return (&UserV2{IdpId: idpId}).Gsi1SkKey()
}

func (p *UserV2) Gsi2PkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString("examplepb_v1_user_v_2:")
	_, _ = sb.WriteString(p.TenantId)
	return sb.String()
}

func UserV2Gsi2PkKey(tenantId string) string {
	return (&UserV2{TenantId: tenantId}).Gsi2PkKey()
}

func (p *UserV2) Gsi2SkKey() string {
	var sb strings.Builder
	sb.Reset()
	_, _ = sb.WriteString(p.IdpId)
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatInt(int64(p.AnEnum), 10))
	return sb.String()
}

func UserV2Gsi2SkKey(idpId string, anEnum BasicEnum) string {
	return (&UserV2{
		AnEnum: anEnum,
		IdpId:  idpId,
	}).Gsi2SkKey()
}
