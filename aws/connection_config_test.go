package aws

import (
	"sync"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func TestGetConfigNormalizesRegionsWithoutMutatingConnection(t *testing.T) {
	regions := []string{"US-EAST-1", "eu-West-1"}
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(awsConfig{Regions: regions})

	config := GetConfig(conn)

	want := []string{"us-east-1", "eu-west-1"}
	for i := range want {
		if config.Regions[i] != want[i] {
			t.Fatalf("GetConfig().Regions[%d] = %q, want %q", i, config.Regions[i], want[i])
		}
	}
	if regions[0] != "US-EAST-1" || regions[1] != "eu-West-1" {
		t.Fatalf("GetConfig changed the shared connection regions to %v", regions)
	}
}

func TestGetConfigConcurrentRegions(t *testing.T) {
	conn := &plugin.Connection{Name: "race"}
	conn.SetConfig(awsConfig{Regions: []string{"us-east-1", "eu-west-1"}})

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = GetConfig(conn)
		}()
	}
	wg.Wait()
}
