package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bondyra/swamp/internal/aws"
	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/aws/definition"
	"github.com/bondyra/swamp/internal/aws/profile"
	"github.com/bondyra/swamp/internal/parser"
	"github.com/bondyra/swamp/internal/reader"
)

func main() {
	query := strings.Join(os.Args[1:], " ")
	ast, err := parser.ParseString(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", ast)
	var reader reader.Reader
	profileProvider := profile.DefaultProvider{}
	awsFactory := client.DefaultFactory{}
	defFactory := definition.DefaultFactory{}
	reader, _ = aws.NewReader(profileProvider, awsFactory, defFactory, []string{})
	fmt.Println(reader.Name())
	d, err3 := awsFactory.NewClient("default")
	fmt.Println(err3)
	p, err4 := d.ListResources("AWS::EC2::VPC")
	fmt.Println(err4)
	fmt.Println(fmt.Printf("%v", p))
}
