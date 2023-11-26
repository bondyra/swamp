package profile

type Factory interface {
	NewProvider() Provider
}

type DefaultFactory struct{}

func (df DefaultFactory) NewProvider() Provider {
	return &DefaultProvider{AwsConfigReader{}}
}
