---
title: "Steampipe Table: aws_codebuild_fleet - Query AWS CodeBuild Fleet using SQL"
description: "Allows users to query AWS CodeBuild Fleet resources to obtain details about compute fleets used for CodeBuild projects, including capacity, status, and configuration."
folder: "CodeBuild"
---

# Table: aws_codebuild_fleet - Query AWS CodeBuild Fleets using SQL

AWS CodeBuild Fleet is a feature that allows you to create and manage dedicated compute resources for your CodeBuild projects. Compute fleets enable you to provision capacity ahead of time, reducing wait times for builds and allowing for more predictable build performance. Fleets help optimize build costs and can be configured with various compute types and capacity settings.

## Table Usage Guide

The `aws_codebuild_fleet` table in Steampipe provides you with information about compute fleets within AWS CodeBuild service. This table allows you, as a DevOps engineer or administrator, to query fleet-specific details, including capacity configurations, status, compute types, and associated metadata. You can utilize this table to gather insights on fleets, such as their current status, capacity utilization, VPC configuration, and more. The schema outlines the various attributes of the CodeBuild fleet for you, including the fleet name, ARN, creation time, capacities, and associated tags.

## Examples

### Basic info
Retrieve fundamental information about all compute fleets in your AWS environment. This query helps you get a quick overview of your fleet configurations, including their names, ARNs, status, and compute specifications.

```sql+postgres
select
  name,
  arn,
  status,
  compute_type,
  environment_type,
  region
from
  aws_codebuild_fleet;
```

```sql+sqlite
select
  name,
  arn,
  status,
  compute_type,
  environment_type,
  region
from
  aws_codebuild_fleet;
```

### List compute fleets by status
Monitor active compute fleets and their creation times. This query is useful for identifying when fleets were created and their current operational status, helping you track fleet lifecycle and troubleshoot any status-related issues.

```sql+postgres
select
  name,
  status,
  status_reason,
  created,
  region
from
  aws_codebuild_fleet
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  name,
  status,
  status_reason,
  created,
  region
from
  aws_codebuild_fleet
where
  status = 'ACTIVE';
```

### Get fleet capacity details
Analyze the capacity configurations of your compute fleets. This query helps you understand your fleet's scaling capabilities by showing current, desired, minimum, maximum, and base capacities, which is essential for capacity planning and optimization.

```sql+postgres
select
  name,
  current_capacity,
  desired_capacity,
  min_capacity,
  max_capacity,
  base_capacity,
  region
from
  aws_codebuild_fleet
order by
  current_capacity desc;
```

```sql+sqlite
select
  name,
  current_capacity,
  desired_capacity,
  min_capacity,
  max_capacity,
  base_capacity,
  region
from
  aws_codebuild_fleet
order by
  current_capacity desc;
```

### Find fleets with VPC configuration
Identify compute fleets that are configured to run within a VPC. This query is valuable for security and networking teams to ensure proper network isolation and access control for build environments.

```sql+postgres
select
  name,
  vpc_config ->> 'VpcId' as vpc_id,
  vpc_config ->> 'Subnets' as subnets,
  vpc_config ->> 'SecurityGroupIds' as security_group_ids,
  region
from
  aws_codebuild_fleet
where
  vpc_config is not null;
```

```sql+sqlite
select
  name,
  json_extract(vpc_config, '$.VpcId') as vpc_id,
  json_extract(vpc_config, '$.Subnets') as subnets,
  json_extract(vpc_config, '$.SecurityGroupIds') as security_group_ids,
  region
from
  aws_codebuild_fleet
where
  vpc_config is not null;
```

### Get fleets by compute type
Analyze the distribution of compute types across your fleets and their total capacity. This query helps in understanding resource allocation and identifying potential areas for optimization or consolidation.

```sql+postgres
select
  compute_type,
  count(*) as fleet_count,
  sum(current_capacity) as total_capacity
from
  aws_codebuild_fleet
group by
  compute_type
order by
  total_capacity desc;
```

```sql+sqlite
select
  compute_type,
  count(*) as fleet_count,
  sum(current_capacity) as total_capacity
from
  aws_codebuild_fleet
group by
  compute_type
order by
  total_capacity desc;
```

### Find recently modified fleets
Track recent changes to your compute fleets. This query helps in change management and auditing by showing fleets that have been modified in the last week, along with their current status and any status-related messages.

```sql+postgres
select
  name,
  last_modified,
  status,
  status_reason,
  region
from
  aws_codebuild_fleet
where
  last_modified > now() - interval '7 days'
order by
  last_modified desc;
```

```sql+sqlite
select
  name,
  last_modified,
  status,
  status_reason,
  region
from
  aws_codebuild_fleet
where
  last_modified > datetime('now', '-7 days')
order by
  last_modified desc;
```

### List fleets by environment type
Group and analyze fleets based on their environment types. This query helps in understanding the distribution of different build environments across your fleets, which is useful for environment standardization and management.

```sql+postgres
select
  environment_type,
  count(*) as fleet_count,
  array_agg(name) as fleet_names
from
  aws_codebuild_fleet
group by
  environment_type;
```

```sql+sqlite
select
  environment_type,
  count(*) as fleet_count,
  group_concat(name) as fleet_names
from
  aws_codebuild_fleet
group by
  environment_type;
```

### Find fleets without proper tagging
Identify fleets that may not comply with tagging standards. This query helps maintain consistent resource tagging by finding fleets that are missing required tags like 'Environment' and 'Project', which are important for resource organization and cost allocation.

```sql+postgres
select
  name,
  region,
  tags
from
  aws_codebuild_fleet
where
  tags is null
  or not tags ? 'Environment'
  or not tags ? 'Project';
```

```sql+sqlite
select
  name,
  region,
  tags
from
  aws_codebuild_fleet
where
  tags is null
  or json_extract(tags, '$.Environment') is null
  or json_extract(tags, '$.Project') is null;
```
