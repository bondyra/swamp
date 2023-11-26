package client

import (
	"fmt"

	"github.com/bondyra/swamp/internal/reader"
)

type Pool interface {
	GetResource(profiles []string, id string, typeName string) ([]*reader.Item, error)
	ListResources(profiles []string, typeName string) ([]*reader.Item, error)
}

type PoolFactory interface {
}

type LazyPool struct {
	clients map[string]AwsClientInterface
	factory ClientFactory
}

type LazyPoolFactory struct {
}

func (lpf LazyPoolFactory) NewPool(profiles []string, factory ClientFactory) (Pool, error) {
	clients := make(map[string]AwsClientInterface, len(profiles))
	for _, p := range profiles {
		clients[p] = nil
	}
	return LazyPool{clients, factory}, nil
}

func (lp LazyPool) GetResource(profiles []string, id string, typeName string) ([]*reader.Item, error) {
	results := make([]*reader.Item, 0)
	for _, p := range profiles {
		it, err := lp.getResourceSingle(p, id, typeName)
		if err != nil {
			return nil, err
		}
		results = append(results, it)
	}
	return results, nil
}

func (lp LazyPool) ListResources(profiles []string, typeName string) ([]*reader.Item, error) {
	results := make([]*reader.Item, 0)
	for _, p := range profiles {
		items, err := lp.listResourcesSingle(p, typeName)
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}
	return results, nil
}

func (lp LazyPool) getResourceSingle(profile string, id string, typeName string) (*reader.Item, error) {
	client, err := lp.getClient(profile)
	if err != nil {
		return nil, err
	}
	resp, err := client.GetResource(id, typeName)
	if err != nil {
		return nil, err
	}
	return &reader.Item{Profile: profile, Data: resp}, nil
}

func (lp LazyPool) listResourcesSingle(profile string, typeName string) ([]*reader.Item, error) {
	client, err := lp.getClient(profile)
	if err != nil {
		return nil, err
	}
	resp, err := client.ListResources(typeName)
	if err != nil {
		return nil, err
	}
	results := make([]*reader.Item, len(resp))
	for _, r := range resp {
		results = append(results, &reader.Item{Profile: profile, Data: r})
	}
	return results, nil
}

func (lp LazyPool) getClient(profile string) (AwsClientInterface, error) {
	client, profileValid := lp.clients[profile]
	if !profileValid {
		return nil, fmt.Errorf("aws profile is not defined: \"%v\"", profile)
	}
	if client == nil {
		newClient, err := lp.factory.NewClient(profile)
		if err != nil {
			return nil, err
		}
		lp.clients[profile] = newClient
	}
	return client, nil
}
