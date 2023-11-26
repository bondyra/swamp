package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
)

type AwsClientInterface interface {
	GetItem(string, string) (map[string]string, error)
	ListItems(string) ([]map[string]string, error)
}

type ccInterface interface {
	GetResource(ctx context.Context, input *cloudcontrol.GetResourceInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceOutput, error)
	ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error)
}

type AwsClient struct {
	ccClient ccInterface
}

func (ac AwsClient) GetItem(id string, typeName string) (map[string]string, error) {
	input := &cloudcontrol.GetResourceInput{
		Identifier: &id,
		TypeName:   &typeName,
	}

	response, err := ac.ccClient.GetResource(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	output, err := ac.processResponse(*response.ResourceDescription.Properties)
	if err != nil {
		return nil, err
	}
	return output, err
}

func (ac AwsClient) ListItems(typeName string) ([]map[string]string, error) {
	input := &cloudcontrol.ListResourcesInput{
		TypeName: &typeName,
	}

	resp, err := ac.ccClient.ListResources(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	outputs := make([]map[string]string, 0)
	for _, rd := range resp.ResourceDescriptions {
		output, err := ac.processResponse(*rd.Properties)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (ac AwsClient) processResponse(response string) (map[string]string, error) {
	output := map[string]string{}
	err := json.Unmarshal([]byte(response), &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
