package aws_endpoint_generator

import (
	"embed"
	"fmt"
)

//go:embed endpoints.json
var endpointJson embed.FS

func GetEmbedEndpointJSONfileContent() ([]byte, error) {
	// Read the embedded JSON file
	jsonData, err := endpointJson.ReadFile("endpoints.json")
	if err != nil {
		return jsonData, fmt.Errorf("error reading embedded JSON file: %v", err)
	}

	return jsonData, nil
}
