package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/smithy-go"
	"github.com/bondyra/swamp/internal/aws/client"
	"github.com/bondyra/swamp/internal/parser"
)

func main() {
	query := strings.Join(os.Args[1:], " ")
	_, err := parser.ParseString(query)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%v", ast)
	pf := client.LazyPoolFactory{}
	p := pf.NewPool("default")
	fmt.Printf("%v - %v", p, err)
	if err == nil {
		r, err := p.GetResource("default", "vpc-c50797ae1", "AWS::EC2::VPC")
		var ae smithy.APIError
		if errors.As(err, &ae) {
			fmt.Printf("%v - %v", r, ae.ErrorCode())
		}
	}
}
