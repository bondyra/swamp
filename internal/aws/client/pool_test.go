package client

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bondyra/swamp/internal/reader"
)

func TestLazyPoolFactory_NewPool(t *testing.T) {
	type args struct {
		profiles []string
		factory  ClientFactory
	}
	tests := []struct {
		name    string
		lpf     LazyPoolFactory
		args    args
		want    Pool
		wantErr bool
	}{
		{
			name:    "test",
			lpf:     LazyPoolFactory{},
			args:    args{profiles: []string{"p1", "p2"}},
			want:    LazyPool{clients: map[string]AwsClientInterface{"p1": nil, "p2": nil}, factory: DefaultClientFactory{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lpf := LazyPoolFactory{}
			got, err := lpf.NewPool(tt.args.profiles, tt.args.factory)
			if (err != nil) != tt.wantErr {
				t.Errorf("LazyPoolFactory.NewPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyPoolFactory.NewPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	PROFILE_OK_1       string = "p_ok_1"
	PROFILE_OK_2              = "p_ok_2"
	PROFILE_ERR_CLIENT        = "p_err_1"
	PROFILE_UNKNOWN           = "p_err_2"
	ID_OK_1                   = "id_ok_1"
	ID_OK_2                   = "id_ok_2"
	ID_ERR                    = "id_err"
	TYPE_OK                   = "type_ok_1"
	TYPE_ERR                  = "type_err_1"
)

var dummyProps1 map[string]string = map[string]string{"a": "1", "b": "1"}
var dummyProps2 map[string]string = map[string]string{"a": "2", "b": "2"}

type mockClient struct{}

func (mc mockClient) GetResource(id, typeName string) (*reader.ItemData, error) {
	if typeName != TYPE_OK {
		return nil, errors.New("invalid type")
	}
	switch id {
	case ID_OK_1:
		return &reader.ItemData{Identifier: id, Properties: &dummyProps1}, nil
	case ID_OK_2:
		return &reader.ItemData{Identifier: id, Properties: &dummyProps2}, nil
	default:
		return nil, errors.New("some error")
	}
}

func (mc mockClient) ListResources(typeName string) ([]*reader.ItemData, error) {
	if typeName != TYPE_OK {
		return nil, errors.New("invalid type")
	}
	return []*reader.ItemData{
		{Identifier: ID_OK_1, Properties: &dummyProps1},
		{Identifier: ID_OK_2, Properties: &dummyProps2},
	}, nil
}

type mockClientFactory struct {
}

func (mcf mockClientFactory) NewClient(profile string) (AwsClientInterface, error) {
	switch profile {
	case PROFILE_OK_1, PROFILE_OK_2:
		return &mockClient{}, nil
	case PROFILE_ERR_CLIENT:
		return nil, errors.New("client creation error")
	default:
		return nil, errors.New("aws profile is not defined")
	}
}

func TestLazyPool_GetResource(t *testing.T) {
	clients := map[string]AwsClientInterface{PROFILE_OK_1: nil, PROFILE_OK_2: nil, PROFILE_ERR_CLIENT: nil}
	type args struct {
		profiles []string
		id       string
		typeName string
	}
	tests := []struct {
		name    string
		args    args
		want    []*reader.Item
		wantErr bool
	}{
		{
			name:    "test empty profiles",
			args:    args{profiles: []string{}, id: ID_OK_1, typeName: TYPE_OK},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name:    "test one profile",
			args:    args{profiles: []string{PROFILE_OK_1}, id: ID_OK_1, typeName: TYPE_OK},
			want:    []*reader.Item{{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_OK_1, Properties: &dummyProps1}}},
			wantErr: false,
		},
		{
			name: "test two profiles",
			args: args{profiles: []string{PROFILE_OK_1, PROFILE_OK_2}, id: ID_OK_1, typeName: TYPE_OK},
			want: []*reader.Item{
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_OK_1, Properties: &dummyProps1}},
				{Profile: PROFILE_OK_2, Data: &reader.ItemData{Identifier: ID_OK_1, Properties: &dummyProps1}},
			},
			wantErr: false,
		},
		{
			name:    "test no inputs on not found id",
			args:    args{profiles: []string{PROFILE_OK_1}, id: ID_ERR, typeName: TYPE_OK},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name:    "test no inputs on not found type",
			args:    args{profiles: []string{PROFILE_OK_1}, id: ID_OK_1, typeName: TYPE_ERR},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name:    "test error on client creation failure",
			args:    args{profiles: []string{PROFILE_OK_1, PROFILE_ERR_CLIENT}, id: ID_OK_1, typeName: TYPE_OK},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test unknown profile",
			args:    args{profiles: []string{PROFILE_OK_1, PROFILE_UNKNOWN}, id: ID_OK_1, typeName: TYPE_OK},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := LazyPool{
				clients: clients,
				factory: mockClientFactory{},
			}
			got, err := lp.GetResource(tt.args.profiles, tt.args.id, tt.args.typeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LazyPool.GetResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyPool.GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyPool_ListResources(t *testing.T) {
	clients := map[string]AwsClientInterface{PROFILE_OK_1: nil, PROFILE_OK_2: nil, PROFILE_ERR_CLIENT: nil}
	type args struct {
		profiles []string
		typeName string
	}
	tests := []struct {
		name    string
		args    args
		want    []*reader.Item
		wantErr bool
	}{
		{
			name:    "test empty profiles",
			args:    args{profiles: []string{}, typeName: TYPE_OK},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name: "test one profile",
			args: args{profiles: []string{PROFILE_OK_1}, typeName: TYPE_OK},
			want: []*reader.Item{
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_OK_1, Properties: &dummyProps1}},
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_OK_2, Properties: &dummyProps2}},
			},
			wantErr: false,
		},
		{
			name: "test two profiles",
			args: args{profiles: []string{PROFILE_OK_1, PROFILE_OK_2}, typeName: TYPE_OK},
			want: []*reader.Item{
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_OK_1, Properties: &dummyProps1}},
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_OK_2, Properties: &dummyProps2}},
				{Profile: PROFILE_OK_2, Data: &reader.ItemData{Identifier: ID_OK_1, Properties: &dummyProps1}},
				{Profile: PROFILE_OK_2, Data: &reader.ItemData{Identifier: ID_OK_2, Properties: &dummyProps2}},
			},
			wantErr: false,
		},
		{
			name:    "test no inputs on not found type",
			args:    args{profiles: []string{PROFILE_OK_1}, typeName: TYPE_ERR},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name:    "test error on client creation failure",
			args:    args{profiles: []string{PROFILE_OK_1, PROFILE_ERR_CLIENT}, typeName: TYPE_OK},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test unknown profile",
			args:    args{profiles: []string{PROFILE_OK_1, PROFILE_UNKNOWN}, typeName: TYPE_OK},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := LazyPool{
				clients: clients,
				factory: mockClientFactory{},
			}
			got, err := lp.ListResources(tt.args.profiles, tt.args.typeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LazyPool.ListResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyPool.ListResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
