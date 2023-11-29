package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bondyra/swamp/internal/aws"
	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/engine"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/language"
)

func loadConfigPaths() []string {
	awsCredentialsPath := os.Getenv("SWAMP_AWS_CREDENTIALS_PATH")
	if awsCredentialsPath == "" {
		home := os.Getenv("HOME")
		awsCredentialsPath = home + "/.aws/credentials"
	}
	return []string{awsCredentialsPath}
}

type Cli struct{}

func (c Cli) Run(query string) {
	ast, err := language.ParseString(query)
	if err != nil {
		log.Fatal(err)
	}

	profileProvider := profile.NewProvider()
	poolFactory := client.LazyPoolFactory{}
	defFactory := definition.DefaultFactory{}
	configPaths := loadConfigPaths()
	reader, err := aws.NewReader(profileProvider, poolFactory, defFactory, configPaths)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.Run(reader, ast)
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output) + "\n")
}
