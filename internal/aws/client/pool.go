package client

import "github.com/bondyra/swamp/internal/reader"

type Pool interface {
	GetResource(profiles []string, id string, typeName string) (map[string]*reader.ItemData, error)
	ListResources(profiles []string, typeName string) (map[string][]*reader.ItemData, error)
}

type PoolFactory struct{}

func (df DefaultFactory) NewPool(profiles []string) (Pool, error) {
	clients := make(map[string]AwsClientInterface, len(profiles))
	for _, p := range profiles {
		clients[p] = nil
	}
	return LazyPool{clients}, nil
}

type LazyPool struct {
	clients map[string]AwsClientInterface
}

func (lp LazyPool) GetResource(profiles []string, id string, typeName string) (map[string]*reader.ItemData, error) {
	return nil, nil
}

func (lp LazyPool) ListResources(profiles []string, typeName string) (map[string][]*reader.ItemData, error) {
	return nil, nil
}
