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
	getResourceProperties   *string
	getResourceError        error
	listResourcesProperties []*string
	listResourcesError      error
}

func (mc mockClient) GetResource(ctx context.Context, input *cloudcontrol.GetResourceInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceOutput, error) {
	output := cloudcontrol.GetResourceOutput{
		ResourceDescription: &cctypes.ResourceDescription{Properties: mc.getResourceProperties},
	}
	return &output, mc.getResourceError
}

func (mc mockClient) ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error) {
	output := cloudcontrol.ListResourcesOutput{}
	for _, p := range mc.listResourcesProperties {
		output.ResourceDescriptions = append(output.ResourceDescriptions, cctypes.ResourceDescription{Properties: p})
	}
	return &output, mc.listResourcesError
}

func TestGetResource(t *testing.T) {

	tests := []struct {
		name               string
		mockClient         ccInterface
		expectedProperties *reader.ItemData
		returnsErr         bool
	}{
		{
			name:               "test empty response",
			mockClient:         mockClient{getResourceProperties: &emptyProperties},
			expectedProperties: &reader.ItemData{Properties: &map[string]string{}},
			returnsErr:         false,
		},
		{
			name:               "test error",
			mockClient:         mockClient{getResourceError: errors.New("some error")},
			expectedProperties: nil,
			returnsErr:         true,
		},
		{
			name:               "test valid response",
			mockClient:         mockClient{getResourceProperties: &someProperties},
			expectedProperties: &reader.ItemData{Properties: &map[string]string{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
			returnsErr:         false,
		},
		{
			name:               "test invalid response",
			mockClient:         mockClient{getResourceProperties: &invalidProperties},
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
			name:           "test empty response",
			mockClient:     mockClient{},
			expectedOutput: []*reader.ItemData{},
			returnsErr:     false,
		},
		{
			name:           "test response with empty properties",
			mockClient:     mockClient{listResourcesProperties: []*string{&emptyProperties, &emptyProperties}},
			expectedOutput: []*reader.ItemData{{Properties: &map[string]string{}}, {Properties: &map[string]string{}}},
			returnsErr:     false,
		},
		{
			name:           "test response with empty properties and valid properties",
			mockClient:     mockClient{listResourcesProperties: []*string{&someProperties, &emptyProperties}},
			expectedOutput: []*reader.ItemData{{Properties: &map[string]string{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}}, &reader.ItemData{&map[string]string{}}},
			returnsErr:     false,
		},
		{
			name:           "test error",
			mockClient:     mockClient{listResourcesError: errors.New("some error")},
			expectedOutput: nil,
			returnsErr:     true,
		},
		{
			name:           "test response with invalid properties",
			mockClient:     mockClient{listResourcesProperties: []*string{&someProperties, &invalidProperties}},
			expectedOutput: nil,
			returnsErr:     true,
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
