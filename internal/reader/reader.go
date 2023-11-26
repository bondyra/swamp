package reader

type Credentials struct {
}

type Reader interface {
	QueryAllProfiles() ([]Credentials, error)
	QueryProfiles(profileNames []string) ([]Credentials, error)
}
