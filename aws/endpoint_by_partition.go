package aws

import (
	"encoding/json"
	"fmt"
	"regexp"

	emb "github.com/turbot/steampipe-plugin-aws/internal/aws_endpoint_generator"
)

// Partition represents an AWS partition in the JSON data.
type Partition struct {
	ID          string             `json:"partition"`
	Name        string             `json:"partitionName"`
	DNSSuffix   string             `json:"dnsSuffix"`
	RegionRegex *regexp.Regexp     `json:"regionRegex"`
	Regions     map[string]Region  `json:"regions"`
	Services    map[string]Service `json:"services"`
}

// Region represents an AWS region.
type Region struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Service represents an AWS service with its endpoints.
type Service struct {
	Endpoints map[string]Endpoint `json:"endpoints"`
}

// Endpoint represents an individual endpoint for a service.
type Endpoint struct {
	CredentialScope *CredentialScope `json:"credentialScope"`
	Hostname        string           `json:"hostname"`
	Deprecated      bool             `json:"deprecated,omitempty"`
	Variants        []Variant        `json:"variants,omitempty"`
}

// CredentialScope defines the credential scope for an endpoint.
type CredentialScope struct {
	Region string `json:"region"`
}

// Variant represents a variant of an endpoint with additional metadata.
type Variant struct {
	Hostname string   `json:"hostname"`
	Tags     []string `json:"tags"`
}

// EndpointsData is the root structure of the JSON data.
type EndpointInfo struct {
	Partitions []Partition `json:"partitions"`
}

func getPartitionValueByPartitionName(partitionId string) (*Partition, error) {

	data, err := emb.GetEmbedEndpointJSONfileContent()
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint json file: %v", err)
	}

	var endpoints EndpointInfo
	if err := json.Unmarshal(data, &endpoints); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, partition := range endpoints.Partitions {
		if partition.ID == partitionId {
			return &partition, nil
		}
	}

	return nil, nil
}
