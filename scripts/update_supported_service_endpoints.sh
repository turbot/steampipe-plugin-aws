#!/bin/bash



# Define the target file
MAIN_FILE="main.go"

# Define the new content for the main.go file
NEW_CONTENT='package main

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-aws/awsEndpointGenerator"
)

// "github.com/turbot/steampipe-plugin-aws/aws"
// "github.com/turbot/steampipe-plugin-sdk/v5/plugin"

func main() {
	// plugin.Serve(&plugin.ServeOpts{
	// 	PluginFunc: aws.Plugin})

	if err := awsEndpointGenerator.Generate(); err != nil {
		fmt.Printf("Error generating Service supported endpoint file: %v\n", err)
	}

	if err := awsEndpointGenerator.GenerateServiceID(); err != nil {
		fmt.Printf("Error generating Service IDs file: %v\n", err)
	}
}
'

# Replace the content of main.go
echo "Replacing content of $MAIN_FILE..."
echo "$NEW_CONTENT" > "$MAIN_FILE"

# Check if replacement was successful
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to replace content in $MAIN_FILE."
  exit 1
fi

echo "Content replaced successfully in $MAIN_FILE."

# Run the Go file
echo "Running $MAIN_FILE..."
go run "$MAIN_FILE"

# Check if the Go command was successful
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to run $MAIN_FILE."
  exit 1
fi

# Revert the main.go file changes

git checkout -- main.go

echo "Script executed successfully."
