package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/bondyra/swamp/internal/aws/common"
)

type Factory interface {
	NewClient(string) (AwsClientInterface, error)
}

type DefaultFactory struct{}

func (df DefaultFactory) NewClient(profile string) (AwsClientInterface, error) {
	context := context.TODO()
	cfg, err := config.LoadDefaultConfig(context, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	return AwsClient{ccClient: cloudcontrol.NewFromConfig(cfg)}, nil
}

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
	output, err := common.Unmarshal[map[string]any]([]byte(response))
	if err != nil {
		return nil, err
	}
	processedOutput := map[string]string{}
	for k := range *output {
		processedOutput[k] = fmt.Sprintf("%v", (*output)[k])
	}
	return processedOutput, nil
}
