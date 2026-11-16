package aws

import (
	"context"
	"testing"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/rate_limiter"
)

// matchingLimiters returns the plugin's default limiters that apply to a call
// with these scope values. It picks them as the SDK does in
// Plugin.getRateLimitersForScopeValues (steampipe-plugin-sdk/v6
// plugin/plugin_rate_limiter.go): the call must have a value for each scope
// key, and the values must satisfy the limiter's where clause.
func matchingLimiters(t *testing.T, p *plugin.Plugin, scopeValues map[string]string) []*rate_limiter.Definition {
	t.Helper()
	var res []*rate_limiter.Definition
	for _, def := range p.RateLimiters {
		if err := def.Initialise(); err != nil {
			t.Fatalf("limiter %s: %v", def.Name, err)
		}
		if len(helpers.FilterMap(scopeValues, def.Scope)) < len(def.Scope) {
			continue
		}
		if def.SatisfiesFilters(scopeValues) {
			res = append(res, def)
		}
	}
	return res
}

// callScope returns the scope values the SDK resolves for a call with these
// tags, in a scan of one region: the scan's values, then the table's tags,
// then the call's tags, each taking precedence over the one before.
func callScope(tbl *plugin.Table, tags map[string]string, region string) map[string]string {
	values := map[string]string{"connection": "aws", "region": region}
	for k, v := range tbl.Tags {
		values[k] = v
	}
	for k, v := range tags {
		values[k] = v
	}
	return values
}

// hydrateScope returns the scope values for a call to hydrate in table.
func hydrateScope(t *testing.T, p *plugin.Plugin, table string, hydrate plugin.HydrateFunc, region string) map[string]string {
	t.Helper()
	tbl := p.TableMap[table]
	name := helpers.GetFunctionName(hydrate)
	for _, c := range tbl.HydrateConfig {
		if helpers.GetFunctionName(c.Func) == name {
			return callScope(tbl, c.Tags, region)
		}
	}
	t.Fatalf("%s has no hydrate config for %s, so no limiter can match its calls", table, name)
	return nil
}

func ec2AttributeLimiter(t *testing.T, p *plugin.Plugin) *rate_limiter.Definition {
	t.Helper()
	for _, def := range p.RateLimiters {
		if def.Name == "aws_ec2_describe_attribute" {
			if err := def.Initialise(); err != nil {
				t.Fatal(err)
			}
			return def
		}
	}
	t.Fatal("plugin has no aws_ec2_describe_attribute limiter")
	return nil
}

// Each per-row EC2 attribute call must be rate limited. Without a limiter the
// SDK runs up to 500 rows at once, EC2 throttles the call, and a table with
// thousands of rows fails when a call runs out of retries (#2819).
func TestEC2AttributeCallsAreRateLimited(t *testing.T) {
	p := Plugin(context.Background())
	cases := []struct {
		table   string
		hydrate plugin.HydrateFunc
	}{
		{"aws_ebs_snapshot", getAwsEBSSnapshotCreateVolumePermissions},
		{"aws_ebs_volume", getVolumeAutoEnableIOData},
		{"aws_ebs_volume", getVolumeProductCodes},
		{"aws_ec2_ami", getAwsEc2AmiLaunchPermissionData},
		{"aws_ec2_instance", getInstanceDisableAPITerminationData},
		{"aws_ec2_instance", getInstanceInitiatedShutdownBehavior},
		{"aws_ec2_instance", getInstanceKernelID},
		{"aws_ec2_instance", getInstanceRAMDiskID},
		{"aws_ec2_instance", getInstanceSriovNetSupport},
		{"aws_ec2_instance", getInstanceUserData},
	}
	for _, c := range cases {
		name := helpers.GetFunctionName(c.hydrate)
		t.Run(c.table+"/"+name, func(t *testing.T) {
			got := matchingLimiters(t, p, hydrateScope(t, p, c.table, c.hydrate, "us-east-1"))
			if len(got) != 1 || got[0].Name != "aws_ec2_describe_attribute" {
				var names []string
				for _, d := range got {
					names = append(names, d.Name)
				}
				t.Fatalf("limiters for %s = %v, want [aws_ec2_describe_attribute]", name, names)
			}
		})
	}
}

// EC2 throttles each action on its own, per account and Region. So calls share
// a bucket only when both the action and the region are the same.
func TestEC2AttributeLimiterHasOneBucketPerActionAndRegion(t *testing.T) {
	p := Plugin(context.Background())
	def := ec2AttributeLimiter(t, p)
	buckets := rate_limiter.NewLimiterMap()
	bucket := func(table string, hydrate plugin.HydrateFunc, region string) *rate_limiter.HydrateLimiter {
		t.Helper()
		l, err := buckets.GetOrCreate(def, helpers.FilterMap(hydrateScope(t, p, table, hydrate, region), def.Scope))
		if err != nil {
			t.Fatal(err)
		}
		return l
	}

	snapshotEast1 := bucket("aws_ebs_snapshot", getAwsEBSSnapshotCreateVolumePermissions, "us-east-1")
	userDataEast1 := bucket("aws_ec2_instance", getInstanceUserData, "us-east-1")

	if bucket("aws_ec2_instance", getInstanceKernelID, "us-east-1") != userDataEast1 {
		t.Error("two DescribeInstanceAttribute hydrates in one region use different buckets; EC2 counts them as one action")
	}
	if snapshotEast1 == userDataEast1 {
		t.Error("DescribeSnapshotAttribute and DescribeInstanceAttribute share a bucket; EC2 throttles each action on its own")
	}
	if bucket("aws_ebs_snapshot", getAwsEBSSnapshotCreateVolumePermissions, "us-east-2") == snapshotEast1 {
		t.Error("us-east-1 and us-east-2 share a bucket; EC2 throttles each Region on its own")
	}
}

// The limiter is for the per-row attribute calls only. List calls use other
// EC2 actions, which have their own buckets in EC2, so it must not slow them.
func TestEC2AttributeLimiterSkipsListCalls(t *testing.T) {
	p := Plugin(context.Background())
	for _, table := range []string{"aws_ebs_snapshot", "aws_ebs_volume", "aws_ec2_ami", "aws_ec2_instance"} {
		tbl := p.TableMap[table]
		for _, def := range matchingLimiters(t, p, callScope(tbl, tbl.List.Tags, "us-east-1")) {
			if def.Name == "aws_ec2_describe_attribute" {
				t.Errorf("%s list call (%s) is limited by aws_ec2_describe_attribute", table, tbl.List.Tags["action"])
			}
		}
	}
}
