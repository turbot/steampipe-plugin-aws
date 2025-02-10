package main

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-aws/internal/awsEndpointGenerator"
)

func main() {
	// Generate the AWS Service supported Endpoints
	if err := awsEndpointGenerator.Generate(); err != nil {
		fmt.Printf("Error generating Service supported endpoint file: %v\n", err)
	}

	if err := awsEndpointGenerator.GenerateServiceID(); err != nil {
		fmt.Printf("Error generating Service IDs file: %v\n", err)
	}
}