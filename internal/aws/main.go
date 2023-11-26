package aws

import (
	"github.com/bondyra/wtf/internal/reader"
)

type AwsReader struct {
}

func (r *AwsReader) QueryAllProfiles() ([]reader.Credentials, error) {
	return []reader.Credentials{}, nil
}

func (r *AwsReader) QueryProfiles(profiles []string) ([]reader.Credentials, error) {
	return []reader.Credentials{}, nil
}
