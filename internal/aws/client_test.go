package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/google/go-cmp/cmp"
)

type mockClient struct {
	getResourceOutput   cloudcontrol.GetResourceOutput
	getResourceError    error
	listResourcesOutput cloudcontrol.ListResourcesOutput
	listResourcesError  error
}

func (mc *mockClient) GetResource(ctx context.Context, input *cloudcontrol.GetResourceInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceOutput, error) {
	return &mc.getResourceOutput, mc.getResourceError
}

func (mc *mockClient) ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error) {
	return &mc.listResourcesOutput, mc.listResourcesError
}

func TestGetItem(t *testing.T) {
	tests := []struct {
		name           string
		mockClient     string
		expectedOutput map[string]string
		expectedError  error
		returnsErr     bool
	}{
		{
			name: "",
			mockClient: mockClient{}, 
			expectedOutput: cloudcontrol.GetResourceOutput(),
			expectedError: ,
			returnsErr: ,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			dapp := DefaultProfileProvider{configReader: MockConfigReader{content: test.configContent}}

			profiles, err := dapp.ProvideProfiles("path")

			if !cmp.Equal(test.expectedProfiles, profiles) {
				t.Errorf("%s expected:\n%v\ngot:\n%v", test.name, test.expectedProfiles, profiles)
			}
			if err != nil {
				t.Errorf("%s error occured: %v", test.name, err)
			}
		})
	}
}
