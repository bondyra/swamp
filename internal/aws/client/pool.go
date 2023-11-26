package client

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/smithy-go"
	"github.com/bondyra/swamp/internal/reader"
)

type ClientFactory interface {
	NewClient(string) (AwsClientInterface, error)
}

type DefaultClientFactory struct {
}

func (dcf DefaultClientFactory) NewClient(profile string) (AwsClientInterface, error) {
	context := context.TODO()
	cfg, err := config.LoadDefaultConfig(context, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	return &AwsClient{ccClient: cloudcontrol.NewFromConfig(cfg)}, nil
}

type Pool interface {
	GetResource(profile string, id string, typeName string) (*reader.Item, error)
	ListResources(profile string, typeName string) ([]*reader.Item, error)
}

type LazyPool struct {
	clients map[string]AwsClientInterface
	factory ClientFactory
}

type PoolFactory interface {
	NewPool(profiles []string) Pool
}

type LazyPoolFactory struct {
}

func (lpf LazyPoolFactory) NewPool(profiles ...string) Pool {
	clients := make(map[string]AwsClientInterface, len(profiles))
	for _, p := range profiles {
		clients[p] = nil
	}
	return LazyPool{clients, DefaultClientFactory{}}
}

func (lp LazyPool) GetResource(profile string, id string, typeName string) (*reader.Item, error) {
	client, err := lp.client(profile)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, nil
	}
	resp, err := client.GetResource(id, typeName)
	if err != nil {
		if lp.shouldSuppressError(err) {
			// todo: add debug logging
			return nil, nil
		}
		return nil, err
	}
	return &reader.Item{Profile: profile, Data: resp}, nil
}

func (lp LazyPool) ListResources(profile string, typeName string) ([]*reader.Item, error) {
	client, err := lp.client(profile)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, nil
	}
	resp, err := client.ListResources(typeName)
	if err != nil {
		if lp.shouldSuppressError(err) {
			// todo: add debug logging
			return []*reader.Item{}, nil
		}
		return nil, err
	}
	results := make([]*reader.Item, 0)
	for _, r := range resp {
		results = append(results, &reader.Item{Profile: profile, Data: r})
	}
	return results, nil
}

func (lp LazyPool) client(profile string) (AwsClientInterface, error) {
	client, profileValid := lp.clients[profile]
	if !profileValid {
		// ignoring profiles that are unknown
		// todo add debug logging
		return nil, nil
	}
	if client == nil {
		newClient, err := lp.factory.NewClient(profile)
		if err != nil {
			return nil, err
		}
		lp.clients[profile] = newClient
	}
	return lp.clients[profile], nil
}

func (lp LazyPool) shouldSuppressError(err error) bool {
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorCode() == "ResourceNotFoundError" {
		return true
	}
	return false
}
