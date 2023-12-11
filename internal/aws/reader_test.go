package aws

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/common"
	"github.com/bondyra/swamp/internal/reader"
	"golang.org/x/exp/slices"
)

func TestNewReader(t *testing.T) {
	type args struct {
		profiles   []string
		createPool client.CreatePool
	}
	tests := []struct {
		name    string
		args    args
		want    *AwsReader
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				profiles:   []string{"p1", "p2"},
				createPool: client.NewLazyPool,
			},
			want: &AwsReader{
				profiles: []string{"p1", "p2"},
				pool:     client.NewLazyPool([]string{"p1", "p2"}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReader(tt.args.profiles, tt.args.createPool)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.profiles, tt.want.profiles) { // pool structs cannot be compared
				t.Errorf("NewReader() = %v, want %v", got.profiles, tt.want.profiles)
			}
		})
	}
}

func TestAwsReader_GetSupportedProfiles(t *testing.T) {
	tests := []struct {
		name     string
		profiles []string
		want     []string
	}{
		{
			name:     "test",
			profiles: []string{"p1", "p2"},
			want:     []string{"p1", "p2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{profiles: tt.profiles}
			if got := ar.GetSupportedProfiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwsReader.GetSupportedProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	ID_1_1 = iota
	ID_1_2
	ID_1_3
	ID_2_1
	ID_2_2
	ID_3_1
	ID_3_ERR
)

var (
	baseItemData []reader.ItemData = []reader.ItemData{
		{Identifier: fmt.Sprint(ID_1_1)},
		{Identifier: fmt.Sprint(ID_1_2)},
		{Identifier: fmt.Sprint(ID_1_3)},
		{Identifier: fmt.Sprint(ID_2_1)},
		{Identifier: fmt.Sprint(ID_2_2)},
		{Identifier: fmt.Sprint(ID_3_1)},
		{Identifier: fmt.Sprint(ID_3_ERR)},
	}
	itemData []reader.ItemData = []reader.ItemData{
		{Identifier: fmt.Sprint(ID_1_1), Properties: &reader.Properties{"a": "a_1_1", "b": "11", "c": "true"}},
		{Identifier: fmt.Sprint(ID_1_2), Properties: &reader.Properties{"a": "a_1_2", "b": "12", "c": "false"}},
		{Identifier: fmt.Sprint(ID_1_3), Properties: &reader.Properties{"a": "a_1_3", "b": "13", "c": "true"}},
		{Identifier: fmt.Sprint(ID_2_1), Properties: &reader.Properties{"a": "a_2_1", "b": "21", "c": "true"}},
		{Identifier: fmt.Sprint(ID_2_2), Properties: &reader.Properties{"a": "a_2_2", "b": "22", "c": "false"}},
		{Identifier: fmt.Sprint(ID_3_1), Properties: &reader.Properties{"a": "a_3_1", "b": "31", "c": "true"}},
		{Identifier: fmt.Sprint(ID_3_ERR), Properties: &reader.Properties{"a": "a_3_2", "b": "32", "c": "false"}},
	}
)

const (
	PROFILE_1                     = "p_1"
	PROFILE_2                     = "p_2"
	TYPE_1                        = "type_1"
	TYPE_1_ID_FIELD               = "type_1_id_field"
	TYPE_2                        = "type_2"
	TYPE_2_ID_FIELD               = "type_2_id_field"
	TYPE_3                        = "type_3"
	TYPE_3_ID_FIELD               = "type_3_id_field"
	TYPE_THAT_CAUSES_ERR          = "type_err"
	TYPE_THAT_CAUSES_ERR_ID_FIELD = "type_err_id_field"
	TYPE_NOT_FOUND_IN_DEFINITION  = "type_def_err"
)

type mockPool struct {
}

func (mp mockPool) ListResources(profile string, typeName string) ([]*reader.Item, error) {
	if typeName == TYPE_THAT_CAUSES_ERR {
		return nil, errors.New("some error")
	}
	switch typeName {
	case TYPE_1:
		return []*reader.Item{{Profile: profile, Data: &baseItemData[ID_1_1]}, {Profile: profile, Data: &baseItemData[ID_1_2]}, {Profile: profile, Data: &baseItemData[ID_1_3]}}, nil
	case TYPE_2:
		return []*reader.Item{{Profile: profile, Data: &baseItemData[ID_2_1]}, {Profile: profile, Data: &baseItemData[ID_2_2]}}, nil
	case TYPE_3:
		return []*reader.Item{{Profile: profile, Data: &baseItemData[ID_3_1]}, {Profile: profile, Data: &baseItemData[ID_3_ERR]}}, nil
	default:
		return []*reader.Item{}, nil
	}
}

func (mp mockPool) GetResource(profile string, id string, typeName string) (*reader.Item, error) {
	if typeName == TYPE_THAT_CAUSES_ERR {
		return nil, errors.New("some error")
	}
	switch typeName {
	case TYPE_1, TYPE_2, TYPE_3:
		intId, err := strconv.Atoi(id)
		if err != nil {
			panic("test setup error, non number id in mock " + id)
		}
		switch intId {
		case ID_3_ERR:
			return nil, errors.New("some error")
		case ID_1_1, ID_1_2, ID_1_3, ID_2_1, ID_2_2, ID_3_1:
			return &reader.Item{Profile: profile, Data: &itemData[intId]}, nil
		default:
			panic("test setup error, unexpected id in mock " + id)
		}
	default:
		panic("test setup error, unexpected type in mock " + typeName)
	}
}

func TestAwsReader_GetItems(t *testing.T) {
	keyFunc := func(a *reader.Item) string { return fmt.Sprintf("%v+%v", a.Profile, a.Data.Identifier) }
	sortFunc := func(a *reader.Item, b *reader.Item) int { return strings.Compare(keyFunc(a), keyFunc(b)) }
	type args struct {
		itemType   string
		profiles   []string
		attrs      []string
		conditions []reader.Condition
	}
	tests := []struct {
		name    string
		args    args
		want    []*reader.Item
		wantErr bool
	}{
		{
			name:    "test no profiles returns nothing",
			args:    args{itemType: TYPE_1, profiles: []string{}, attrs: nil, conditions: nil},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name: "test one profile no id condition",
			args: args{itemType: TYPE_1, profiles: []string{PROFILE_1}, attrs: nil, conditions: nil},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_3]},
			},
			wantErr: false,
		},
		{
			name: "test multiple profiles no id condition",
			args: args{itemType: TYPE_1, profiles: []string{PROFILE_1, PROFILE_2}, attrs: nil, conditions: nil},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_1, Data: &itemData[ID_1_3]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_2]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_3]},
			},
			wantErr: false,
		},
		{
			name: "test multiple profiles with one id condition",
			args: args{
				itemType:   TYPE_1,
				profiles:   []string{PROFILE_1, PROFILE_2},
				attrs:      nil,
				conditions: []reader.Condition{{Attr: TYPE_1_ID_FIELD, Op: common.EqualsTo, Value: fmt.Sprint(ID_1_1)}},
			},
			want: []*reader.Item{
				{Profile: PROFILE_1, Data: &itemData[ID_1_1]},
				{Profile: PROFILE_2, Data: &itemData[ID_1_1]},
			},
			wantErr: false,
		},
		{
			name:    "test error when type is not found in definition",
			args:    args{itemType: TYPE_NOT_FOUND_IN_DEFINITION, profiles: []string{PROFILE_1}, attrs: nil, conditions: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test error on list resources error",
			args:    args{itemType: TYPE_THAT_CAUSES_ERR, profiles: []string{PROFILE_1}, attrs: nil, conditions: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test error on get resource error",
			args:    args{itemType: TYPE_3, profiles: []string{PROFILE_1}, attrs: nil, conditions: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test error on not supported condition op",
			args: args{
				itemType:   TYPE_1,
				profiles:   []string{PROFILE_1, PROFILE_2},
				attrs:      nil,
				conditions: []reader.Condition{{Attr: TYPE_1_ID_FIELD, Op: -1, Value: fmt.Sprint(ID_1_1)}},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test error on not supported condition",
			args: args{
				itemType:   TYPE_1,
				profiles:   []string{PROFILE_1, PROFILE_2},
				attrs:      nil,
				conditions: []reader.Condition{{Attr: "unknown", Op: common.EqualsTo, Value: fmt.Sprint(ID_1_1)}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := AwsReader{
				pool: mockPool{},
			}
			got, err := ar.GetItems(tt.args.itemType, tt.args.profiles, tt.args.attrs, tt.args.conditions)
			slices.SortFunc(got, sortFunc)
			slices.SortFunc(tt.want, sortFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("AwsReader.GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwsReader.GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
