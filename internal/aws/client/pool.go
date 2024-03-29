package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/smithy-go"
	"github.com/bondyra/swamp/internal/reader"
)

type newClient func(string) (AwsClientInterface, error)

func newDefaultClient(profile string) (AwsClientInterface, error) {
	// todo test coverage for aws client
	context := context.TODO()
	cfg, err := config.LoadDefaultConfig(context, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	return &AwsClient{ccClient: cloudcontrol.NewFromConfig(cfg)}, nil
}

type Pool interface {
	GetResource(profile string, id string, typeName string) (*reader.Item, error)
	ListResources(profile string, typeName string) ([]*reader.Item, error)
}

type CreatePool func(profiles []string) Pool

type LazyPool struct {
	clients      map[string]AwsClientInterface
	createClient newClient
}

func NewLazyPool(profiles []string) Pool {
	clients := make(map[string]AwsClientInterface, len(profiles))
	for _, p := range profiles {
		clients[p] = nil
	}
	return LazyPool{clients, newDefaultClient}
}

func (lp LazyPool) GetResource(profile string, id string, typeName string) (*reader.Item, error) {
	client, err := lp.client(profile)
	if err != nil {
		return nil, fmt.Errorf("GetResource: %w", err)
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
		return nil, fmt.Errorf("GetResource: %w", err)
	}
	return &reader.Item{Profile: profile, Data: resp}, nil
}

func (lp LazyPool) ListResources(profile string, typeName string) ([]*reader.Item, error) {
	client, err := lp.client(profile)
	if err != nil {
		return nil, fmt.Errorf("ListResources: %w", err)
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
		return nil, fmt.Errorf("ListResources: %w", err)
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
		newClient, err := lp.createClient(profile)
		if err != nil {
			return nil, fmt.Errorf("client: %w", err)
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
