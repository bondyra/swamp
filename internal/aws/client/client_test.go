package client

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	types "github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/google/go-cmp/cmp"
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
		ResourceDescription: &types.ResourceDescription{Properties: mc.getResourceProperties},
	}
	return &output, mc.getResourceError
}

func (mc mockClient) ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error) {
	output := cloudcontrol.ListResourcesOutput{}
	for _, p := range mc.listResourcesProperties {
		output.ResourceDescriptions = append(output.ResourceDescriptions, types.ResourceDescription{Properties: p})
	}
	return &output, mc.listResourcesError
}

func TestGetItem(t *testing.T) {

	tests := []struct {
		name               string
		mockClient         ccInterface
		expectedProperties map[string]string
		returnsErr         bool
	}{
		{
			name:               "test empty response",
			mockClient:         mockClient{getResourceProperties: &emptyProperties},
			expectedProperties: map[string]string{},
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
			expectedProperties: map[string]string{"str": "abc", "int": "1", "float": "1.23", "bool": "true"},
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

			actualProperties, err := a.GetItem("id", "type")

			if !cmp.Equal(actualProperties, test.expectedProperties) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProperties, actualProperties)
			}
			if test.returnsErr {
				if err == nil {
					t.Errorf("expected:\nerror\ngot:\n%v", err)
				}
			} else {
				if err != nil {
					t.Errorf("%s error occured: %v", test.name, err)
				}
			}
		})
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name           string
		mockClient     ccInterface
		expectedOutput []map[string]string
		returnsErr     bool
	}{
		{
			name:           "test empty response",
			mockClient:     mockClient{},
			expectedOutput: []map[string]string{},
			returnsErr:     false,
		},
		{
			name:           "test response with empty properties",
			mockClient:     mockClient{listResourcesProperties: []*string{&emptyProperties, &emptyProperties}},
			expectedOutput: []map[string]string{{}, {}},
			returnsErr:     false,
		},
		{
			name:           "test response with valid properties",
			mockClient:     mockClient{listResourcesProperties: []*string{&someProperties, &someProperties}},
			expectedOutput: []map[string]string{{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}, {"str": "abc", "int": "1", "float": "1.23", "bool": "true"}},
			returnsErr:     false,
		},
		{
			name:           "test response with empty properties and valid properties",
			mockClient:     mockClient{listResourcesProperties: []*string{&someProperties, &emptyProperties}},
			expectedOutput: []map[string]string{{"str": "abc", "int": "1", "float": "1.23", "bool": "true"}, {}},
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

			actualProperties, err := a.ListItems("type")

			if !cmp.Equal(actualProperties, test.expectedOutput) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedOutput, actualProperties)
			}
			if test.returnsErr {
				if err == nil {
					t.Errorf("expected:\nerror\ngot:\n%v", err)
				}
			} else {
				if err != nil {
					t.Errorf("%s error occured: %v", test.name, err)
				}
			}
		})
	}
}
