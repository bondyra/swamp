package aws

type AwsConnection struct{}

type AwsConnectionPool struct {
	connections map[string]AwsConnection
}

func NewConnectionPool(profiles []string) (AwsConnectionPool, error) {
	connections := make(map[string]AwsConnection, 0)
	for _, profile := range profiles {
		connections[profile] = AwsConnection{}
	}
	return AwsConnectionPool{connections: connections}, nil
}
