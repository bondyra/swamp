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
	emptyProperties   string = ""
	someProperties    string = "{\"str\":\"abc\", \"int\":1, \"float\": 1.23, \"bool\": true}"
	invalidProperties string = "{invalid"
)

type mockClient struct {
	getResourceOutput    *cloudcontrol.GetResourceOutput
	getResourceError     error
	listResourcesOutputs *cloudcontrol.ListResourcesOutput
	listResourcesError   error
}

func (mc mockClient) GetResource(ctx context.Context, input *cloudcontrol.GetResourceInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceOutput, error) {
	return mc.getResourceOutput, mc.getResourceError
}

func (mc mockClient) ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error) {
	return mc.listResourcesOutputs, mc.listResourcesError
}

func TestGetResource(t *testing.T) {
	tests := []struct {
		name               string
		mockClient         ccInterface
		expectedProperties *reader.ItemData
		returnsErr         bool
	}{
		{
			name:               "test nil response",
			mockClient:         mockClient{},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name:               "test error",
			mockClient:         mockClient{getResourceError: errors.New("some error")},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name: "test valid response",
			mockClient: mockClient{getResourceOutput: &cloudcontrol.GetResourceOutput{
				ResourceDescription: &cctypes.ResourceDescription{Properties: &someProperties},
			}},
			expectedProperties: &reader.ItemData{Properties: &map[string]string{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
			returnsErr:         false,
		},
		{
			name: "test empty response",
			mockClient: mockClient{getResourceOutput: &cloudcontrol.GetResourceOutput{
				ResourceDescription: &cctypes.ResourceDescription{Properties: &emptyProperties},
			}},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name: "test invalid response",
			mockClient: mockClient{getResourceOutput: &cloudcontrol.GetResourceOutput{
				ResourceDescription: &cctypes.ResourceDescription{Properties: &invalidProperties},
			}},
			expectedProperties: nil,
			returnsErr:         true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := AwsClient{test.mockClient}

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
		mockClient     ccInterface
		expectedOutput []*reader.ItemData
		returnsErr     bool
	}{
		{
			name:           "test nil reponse",
			mockClient:     mockClient{listResourcesOutputs: nil},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with empty properties",
			mockClient: mockClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Properties: &emptyProperties},
					{Properties: &emptyProperties},
				},
			}},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with empty properties and valid properties",
			mockClient: mockClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Properties: &someProperties},
					{Properties: &emptyProperties},
				},
			}},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name:           "test error",
			mockClient:     mockClient{listResourcesError: errors.New("some error")},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with invalid properties",
			mockClient: mockClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Properties: &someProperties},
					{Properties: &invalidProperties},
				},
			}},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name: "test response with valid properties",
			mockClient: mockClient{listResourcesOutputs: &cloudcontrol.ListResourcesOutput{
				ResourceDescriptions: []cctypes.ResourceDescription{
					{Properties: &someProperties},
					{Properties: &someProperties},
				},
			}},
			expectedOutput: []*reader.ItemData{
				{Properties: &map[string]string{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
				{Properties: &map[string]string{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
			},
			returnsErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := AwsClient{test.mockClient}

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
