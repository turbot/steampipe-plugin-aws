---
title: "Steampipe Table: aws_vpc_block_public_access_options - Query AWS VPC Block Public Access Options using SQL"
description: "Allows users to query AWS VPC Block Public Access Options to retrieve details about the block public access configurations for VPCs."
folder: "VPC"
---

# Table: aws_vpc_block_public_access_options - Query AWS VPC Block Public Access Options using SQL

The AWS VPC Block Public Access Options is a regional security feature that helps protect your VPC resources from public accessibility. It allows you to control the inbound network connections to Amazon VPCs at the regional level, preventing unauthorized access. This regional configuration is crucial for ensuring the security and privacy of your network resources on AWS across entire regions.

## Table Usage Guide

The `aws_vpc_block_public_access_options` table in Steampipe provides you with information about the regional VPC Block Public Access (BPA) configurations. This table allows you, as a cloud administrator, security team member, or DevOps engineer, to query regional configuration details, including the VPC BPA mode, exclusions allowed status, and management information. You can utilize this table to gather insights on regional configurations, such as which regions have public access blocked, the current state of VPC BPA, and who manages the configuration. The schema outlines the various attributes of the regional block public access configuration for you, including the internet gateway block mode, exclusions allowed status, and management details.

## Examples

### Basic settings info
Analyze the VPC Block Public Access settings to understand the security configuration across your AWS regions. This is useful for ensuring your VPC resources are properly protected and for maintaining compliance with your organization's security policies.

```sql+postgres
select
  exclusions_allowed,
  internet_gateway_block_mode,
  last_update_timestamp,
  managed_by,
  reason,
  state
from
  aws_vpc_block_public_access_options;
```

```sql+sqlite
select
  exclusions_allowed,
  internet_gateway_block_mode,
  last_update_timestamp,
  managed_by,
  reason,
  state
from
  aws_vpc_block_public_access_options;
```

### List configurations with internet gateway blocking enabled
Identify VPC Block Public Access configurations that have internet gateway blocking enabled. This is useful for maintaining network security by ensuring that public access through internet gateways is properly blocked.

```sql+postgres
select
  internet_gateway_block_mode,
  state,
  managed_by
from
  aws_vpc_block_public_access_options
where
  internet_gateway_block_mode != 'off';
```

```sql+sqlite
select
  internet_gateway_block_mode,
  state,
  managed_by
from
  aws_vpc_block_public_access_options
where
  internet_gateway_block_mode != 'off';
```

### List configurations managed by declarative policies
Discover VPC Block Public Access configurations that are managed by declarative policies. This is useful for understanding your security setup and identifying which configurations are centrally managed through policy-based controls.

```sql+postgres
select
  managed_by,
  state,
  reason
from
  aws_vpc_block_public_access_options
where
  managed_by = 'declarative-policy';
```

```sql+sqlite
select
  managed_by,
  state,
  reason
from
  aws_vpc_block_public_access_options
where
  managed_by = 'declarative-policy';
```

### List configurations with bidirectional internet gateway blocking
Explore VPC Block Public Access configurations that have bidirectional internet gateway blocking enabled. This provides maximum security against public access and is useful for environments requiring the highest level of network protection.

```sql+postgres
select
  internet_gateway_block_mode,
  state,
  exclusions_allowed
from
  aws_vpc_block_public_access_options
where
  internet_gateway_block_mode = 'block-bidirectional';
```

```sql+sqlite
select
  internet_gateway_block_mode,
  state,
  exclusions_allowed
from
  aws_vpc_block_public_access_options
where
  internet_gateway_block_mode = 'block-bidirectional';
```

### Check BPA configuration status
Analyze the current status of VPC Block Public Access configurations across all regions. This is useful for understanding your security posture and ensuring consistent protection across your AWS environment.

```sql+postgres
select
  region,
  state,
  internet_gateway_block_mode,
  managed_by,
  last_update_timestamp
from
  aws_vpc_block_public_access_options
order by
  region;
```

```sql+sqlite
select
  region,
  state,
  internet_gateway_block_mode,
  managed_by,
  last_update_timestamp
from
  aws_vpc_block_public_access_options
order by
  region;
```

### Find configurations with exclusions not allowed
Identify VPC Block Public Access configurations where exclusions are not allowed. This indicates strict security policies and is useful for compliance auditing and security assessment.

```sql+postgres
select
  region,
  exclusions_allowed,
  internet_gateway_block_mode,
  managed_by
from
  aws_vpc_block_public_access_options
where
  exclusions_allowed = 'not-allowed';
```

```sql+sqlite
select
  region,
  exclusions_allowed,
  internet_gateway_block_mode,
  managed_by
from
  aws_vpc_block_public_access_options
where
  exclusions_allowed = 'not-allowed';
```

### Monitor recent BPA configuration changes
Track recent changes to VPC Block Public Access configurations across regions. This is useful for security auditing and maintaining an audit trail of configuration modifications.

```sql+postgres
select
  region,
  last_update_timestamp,
  internet_gateway_block_mode,
  state,
  reason
from
  aws_vpc_block_public_access_options
where
  last_update_timestamp >= now() - interval '30 days'
order by
  last_update_timestamp desc;
```

```sql+sqlite
select
  region,
  last_update_timestamp,
  internet_gateway_block_mode,
  state,
  reason
from
  aws_vpc_block_public_access_options
where
  last_update_timestamp >= datetime('now', '-30 days')
order by
  last_update_timestamp desc;
```
