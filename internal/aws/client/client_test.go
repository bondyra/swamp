package client

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	cctypes "github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/bondyra/swamp/internal/reader"
)

var (
	dummyId           string = "id"
	emptyProperties   string = ""
	someProperties    string = "{\"str\":\"abc\", \"int\":1, \"float\": 1.23, \"bool\": true}"
	invalidProperties string = "{invalid"
)

type mockccClient struct {
	getResourceOutput    *cloudcontrol.GetResourceOutput
	getResourceError     error
	listResourcesOutputs *cloudcontrol.ListResourcesOutput
	listResourcesError   error
}

func (mc mockccClient) GetResource(ctx context.Context, input *cloudcontrol.GetResourceInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceOutput, error) {
	return mc.getResourceOutput, mc.getResourceError
}

func (mc mockccClient) ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error) {
	return mc.listResourcesOutputs, mc.listResourcesError
}

func TestGetResource(t *testing.T) {
	tests := []struct {
		name               string
		mockccClient       ccInterface
		expectedProperties *reader.ItemData
		returnsErr         bool
	}{
		{
			name:               "test nil response",
			mockccClient:       mockccClient{},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name:               "test error",
			mockccClient:       mockccClient{getResourceError: errors.New("some error")},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name: "test valid response",
			mockccClient: mockccClient{getResourceOutput: &cloudcontrol.GetResourceOutput{
				ResourceDescription: &cctypes.ResourceDescription{Identifier: &dummyId, Properties: &someProperties},
			}},
			expectedProperties: &reader.ItemData{Identifier: dummyId, Properties: &reader.Properties{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
			returnsErr:         false,
		},
		{
			name: "test empty response",
			mockccClient: mockccClient{getResourceOutput: &cloudcontrol.GetResourceOutput{
				ResourceDescription: &cctypes.ResourceDescription{Identifier: &dummyId, Properties: &emptyProperties},
			}},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name: "test invalid response",
			mockccClient: mockccClient{getResourceOutput: &cloudcontrol.GetResourceOutput{
				ResourceDescription: &cctypes.ResourceDescription{Identifier: &dummyId, Properties: &invalidProperties},
			}},
			expectedProperties: nil,
			returnsErr:         true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := AwsClient{test.mockccClient}

			actualProperties, err := a.GetResource("id", "type")

			if test.returnsErr {
				if err == nil {
					t.Errorf("expected:\nerror\ngot:\n%v", err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("%s error occured: %v", test.name, err)
				}
			}
			if !reflect.DeepEqual(*actualProperties.Properties, *test.expectedProperties.Properties) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProperties.Properties, actualProperties.Properties)
			}
		})
	}
}

func TestListResources(t *testing.T) {
	tests := []struct {
		name           string
		mockccClient   ccInterface
		expectedOutput []*reader.ItemData
		returnsErr     bool
	}{
		{
			name:           "test nil reponse",
			mockccClient:   mockccClient{listResourcesOutputs: nil},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with empty properties",
			mockccClient: mockccClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Identifier: &dummyId, Properties: &emptyProperties},
					{Identifier: &dummyId, Properties: &emptyProperties},
				},
			}},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with empty properties and valid properties",
			mockccClient: mockccClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Identifier: &dummyId, Properties: &someProperties},
					{Identifier: &dummyId, Properties: &emptyProperties},
				},
			}},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name:           "test error",
			mockccClient:   mockccClient{listResourcesError: errors.New("some error")},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with invalid properties",
			mockccClient: mockccClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Identifier: &dummyId, Properties: &someProperties},
					{Identifier: &dummyId, Properties: &invalidProperties},
				},
			}},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with valid properties",
			mockccClient: mockccClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Identifier: &dummyId, Properties: &someProperties},
					{Identifier: &dummyId, Properties: &someProperties},
				},
			}},
			expectedOutput: []*reader.ItemData{
				{Identifier: dummyId, Properties: &reader.Properties{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
				{Identifier: dummyId, Properties: &reader.Properties{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
			},
			returnsErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := AwsClient{test.mockccClient}

			actualProperties, err := a.ListResources("type")

			if test.returnsErr {
				if err == nil {
					t.Errorf("expected:\nerror\ngot:\n%v", err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("%s error occured: %v", test.name, err)
				}
			}
			if !reflect.DeepEqual(actualProperties, test.expectedOutput) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedOutput, actualProperties)
			}
		})
	}
}
