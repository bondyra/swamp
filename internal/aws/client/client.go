package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/reader"
)

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

func (ac *AwsClient) GetResource(id string, typeName string) (*reader.ItemData, error) {
	input := &cloudcontrol.GetResourceInput{
		Identifier: &id,
		TypeName:   &typeName,
	}

	resp, err := ac.ccClient.GetResource(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("GetResource: %w", err)
	}
	if resp == nil {
		return nil, errors.New("unexpected null AWS GetResource response")
	}

	props, err := ac.parseProperties(*resp.ResourceDescription.Properties)
	if err != nil {
		return nil, fmt.Errorf("GetResource: %w", err)
	}
	return &reader.ItemData{Identifier: *resp.ResourceDescription.Identifier, Properties: props}, err
}

func (ac *AwsClient) ListResources(typeName string) ([]*reader.ItemData, error) {
	input := &cloudcontrol.ListResourcesInput{
		TypeName: &typeName,
	}

	resp, err := ac.ccClient.ListResources(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("ListResources: %w", err)
	}
	if resp == nil {
		return nil, errors.New("unexpected null AWS ListResources response")
	}

	outputs := make([]*reader.ItemData, 0)
	for _, rd := range resp.ResourceDescriptions {
		props, err := ac.parseProperties(*rd.Properties)
		if err != nil {
			return nil, fmt.Errorf("ListResources: %w", err)
		}
		outputs = append(outputs, &reader.ItemData{Identifier: *rd.Identifier, Properties: props})
	}
	return outputs, nil
}

func (ac *AwsClient) parseProperties(response string) (*reader.Properties, error) {
	output, err := common.Unmarshal[map[string]any]([]byte(response))
	if err != nil {
		return nil, fmt.Errorf("parseProperties: %w", err)
	}
	processedOutput := reader.Properties{}
	for k := range *output {
		processedOutput[k] = fmt.Sprintf("%v", (*output)[k])
	}
	return &processedOutput, nil
}
