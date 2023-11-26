package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
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
