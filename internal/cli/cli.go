package cli

import (
	"log"
	"os"

	"github.com/bondyra/swamp/internal/aws"
	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/common"
	"github.com/bondyra/swamp/internal/engine"
	"github.com/bondyra/swamp/internal/language"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/bondyra/swamp/internal/topology"
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
	r, err := aws.NewReader(profiles, client.NewLazyPool)
	if err != nil {
		log.Fatal(err)
	}
	err = engine.NewEngine(
		topology.ReaderTopologyLoader([]reader.Reader{r}),
		engine.DefaultExecutionPlanner,
		engine.ParallelExecutionRunner,
		engine.PrintYamlPlotter,
	).Run(ast, map[string]reader.Reader{"aws": r}, common.DebugVerbosity)
	if err != nil {
		log.Fatal(err)
	}
}
