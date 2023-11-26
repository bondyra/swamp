package client

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/bondyra/swamp/internal/reader"
)

func TestLazyPoolFactory_NewPool(t *testing.T) {
	type args struct {
		profiles []string
	}
	tests := []struct {
		name string
		lpf  LazyPoolFactory
		args args
		want Pool
	}{
		{
			name: "test",
			lpf:  LazyPoolFactory{},
			args: args{profiles: []string{"p1", "p2"}},
			want: LazyPool{clients: map[string]AwsClientInterface{"p1": nil, "p2": nil}, factory: DefaultClientFactory{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lpf := LazyPoolFactory{}
			if got := lpf.NewPool(tt.args.profiles...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyPoolFactory.NewPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	PROFILE_OK_1               = "p_ok_1"
	PROFILE_OK_2               = "p_ok_2"
	PROFILE_ERR_CLIENT         = "p_err_1"
	PROFILE_UNKNOWN            = "p_err_2"
	ID_1                       = "id_ok_1"
	ID_2                       = "id_ok_2"
	TYPE_OK                    = "type_ok_1"
	TYPE_THAT_CAUSES_NOT_FOUND = "type_not_found"
	TYPE_THAT_CAUSES_ERR       = "type_err"
)

var dummyProps1 map[string]string = map[string]string{"a": "1", "b": "1"}
var dummyProps2 map[string]string = map[string]string{"a": "2", "b": "2"}

type mockClient struct {
}

type mockAwsError struct {
	errorCode string
}

func (mockAwsError) Error() string {
	return "mock aws error"
}

func (mae mockAwsError) ErrorCode() string {
	return mae.errorCode
}
func (mae mockAwsError) ErrorMessage() string {
	return "test"
}

func (mae mockAwsError) ErrorFault() smithy.ErrorFault {
	return smithy.FaultUnknown
}

func (mc mockClient) GetResource(id, typeName string) (*reader.ItemData, error) {
	switch typeName {
	case TYPE_OK:
		return &reader.ItemData{Identifier: id, Properties: &dummyProps1}, nil
	case TYPE_THAT_CAUSES_NOT_FOUND:
		return nil, mockAwsError{"ResourceNotFoundError"}
	default:
		return nil, mockAwsError{"AnyOtherErrorCode"}
	}
}

func (mc mockClient) ListResources(typeName string) ([]*reader.ItemData, error) {
	switch typeName {
	case TYPE_OK:
		return []*reader.ItemData{
			{Identifier: ID_1, Properties: &dummyProps1},
			{Identifier: ID_2, Properties: &dummyProps2},
		}, nil
	case TYPE_THAT_CAUSES_NOT_FOUND:
		return nil, mockAwsError{"ResourceNotFoundError"}
	default:
		return nil, mockAwsError{"AnyOtherErrorCode"}
	}
}

type MockClientFactory struct {
}

func (mcf MockClientFactory) NewClient(profile string) (AwsClientInterface, error) {
	switch profile {
	case PROFILE_OK_1, PROFILE_OK_2:
		return &mockClient{}, nil
	default:
		return nil, fmt.Errorf("unexpected test profile: %v", profile)
	}
}

func TestLazyPool_GetResource(t *testing.T) {
	clients := map[string]AwsClientInterface{PROFILE_OK_1: nil, PROFILE_OK_2: nil, PROFILE_ERR_CLIENT: nil}
	type args struct {
		profile  string
		id       string
		typeName string
	}
	tests := []struct {
		name    string
		args    args
		want    *reader.Item
		wantErr bool
	}{
		{
			name:    "test profile",
			args:    args{profile: PROFILE_OK_1, id: ID_1, typeName: TYPE_OK},
			want:    &reader.Item{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_1, Properties: &dummyProps1}},
			wantErr: false,
		},
		{
			name:    "test ignore unknown profiles",
			args:    args{profile: PROFILE_UNKNOWN, id: ID_1, typeName: TYPE_OK},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "test no inputs on not found aws error",
			args:    args{profile: PROFILE_OK_1, id: ID_1, typeName: TYPE_THAT_CAUSES_NOT_FOUND},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "test error on any other aws error",
			args:    args{profile: PROFILE_OK_1, id: ID_1, typeName: TYPE_THAT_CAUSES_ERR},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test error on client creation failure",
			args:    args{profile: PROFILE_ERR_CLIENT, id: ID_1, typeName: TYPE_OK},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := LazyPool{
				clients: clients,
				factory: MockClientFactory{},
			}
			got, err := lp.GetResource(tt.args.profile, tt.args.id, tt.args.typeName)
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
		profile  string
		typeName string
	}
	tests := []struct {
		name    string
		args    args
		want    []*reader.Item
		wantErr bool
	}{
		{
			name: "test one profile",
			args: args{profile: PROFILE_OK_1, typeName: TYPE_OK},
			want: []*reader.Item{
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_1, Properties: &dummyProps1}},
				{Profile: PROFILE_OK_1, Data: &reader.ItemData{Identifier: ID_2, Properties: &dummyProps2}},
			},
			wantErr: false,
		},
		{
			name:    "test ignore unknown profiles",
			args:    args{profile: PROFILE_UNKNOWN, typeName: TYPE_OK},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "test no inputs on not found aws error",
			args:    args{profile: PROFILE_OK_1, typeName: TYPE_THAT_CAUSES_NOT_FOUND},
			want:    []*reader.Item{},
			wantErr: false,
		},
		{
			name:    "test error on any other aws error",
			args:    args{profile: PROFILE_OK_1, typeName: TYPE_THAT_CAUSES_ERR},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test error on client creation failure",
			args:    args{profile: PROFILE_ERR_CLIENT, typeName: TYPE_OK},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := LazyPool{
				clients: clients,
				factory: MockClientFactory{},
			}
			got, err := lp.ListResources(tt.args.profile, tt.args.typeName)
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
