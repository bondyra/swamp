package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/reader"
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
	GetResource(string, string) (*reader.ItemData, error)
	ListResources(string) ([]*reader.ItemData, error)
}

type ccInterface interface {
	GetResource(ctx context.Context, input *cloudcontrol.GetResourceInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.GetResourceOutput, error)
	ListResources(ctx context.Context, input *cloudcontrol.ListResourcesInput, optFns ...func(*cloudcontrol.Options)) (*cloudcontrol.ListResourcesOutput, error)
}

type AwsClient struct {
	ccClient ccInterface
}

func (ac AwsClient) GetResource(id string, typeName string) (*reader.ItemData, error) {
	input := &cloudcontrol.GetResourceInput{
		Identifier: &id,
		TypeName:   &typeName,
	}

	resp, err := ac.ccClient.GetResource(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("unexpected null cc ListResources response")
	}

	props, err := ac.parseProperties(*resp.ResourceDescription.Properties)
	if err != nil {
		return nil, err
	}
	return &reader.ItemData{Identifier: *resp.ResourceDescription.Identifier, Properties: props}, err
}

func (ac AwsClient) ListResources(typeName string) ([]*reader.ItemData, error) {
	input := &cloudcontrol.ListResourcesInput{
		TypeName: &typeName,
	}

	resp, err := ac.ccClient.ListResources(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("unexpected null cc ListResources response")
	}

	outputs := make([]*reader.ItemData, 0)
	for _, rd := range resp.ResourceDescriptions {
		props, err := ac.parseProperties(*rd.Properties)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, &reader.ItemData{Identifier: *rd.Identifier, Properties: props})
	}
	return outputs, nil
}

func (ac AwsClient) parseProperties(response string) (*map[string]string, error) {
	output, err := common.Unmarshal[map[string]any]([]byte(response))
	if err != nil {
		return nil, err
	}
	processedOutput := map[string]string{}
	for k := range *output {
		processedOutput[k] = fmt.Sprintf("%v", (*output)[k])
	}
	return &processedOutput, nil
}
