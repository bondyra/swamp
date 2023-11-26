package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
)

type ProfileFactory interface {
	NewProfileProvider() ProfileProvider
}

type DefaultProfileFactory struct{}

func (dpf DefaultProfileFactory) NewProfileProvider() ProfileProvider {
	return &DefaultProfileProvider{AwsConfigReader{}}
}

type AwsFactory interface {
	NewClient(string) (AwsClientInterface, error)
}

type DefaultAwsFactory struct{}

func (daf DefaultAwsFactory) NewClient(profile string) (AwsClientInterface, error) {
	context := context.TODO()
	cfg, err := config.LoadDefaultConfig(context, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	return AwsClient{cloudcontrol.NewFromConfig(cfg)}, nil
}
