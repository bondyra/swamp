package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/bondyra/swamp/internal/aws"
	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/engine"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
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
	configPaths := loadConfigPaths()
	profiles, err := profile.FromConfigFiles(configPaths...)
	if err != nil {
		log.Fatal(err)
	}
	_, filename, _, _ := runtime.Caller(0)
	definition, err := definition.FromFile(path.Dir(filename) + "/definition.json")
	if err != nil {
		log.Fatal(err)
	}
	r := aws.NewReader(profiles, client.NewLazyPool, definition)

	result, err := engine.Run(ast, []reader.Reader{r})
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output) + "\n")
}
