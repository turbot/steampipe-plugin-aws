---
title: "Steampipe Table: aws_vpc_block_public_access_options - Query AWS VPC Block Public Access Options using SQL"
description: "Allows users to query AWS VPC Block Public Access Options to retrieve details about the block public access configurations for VPCs."
folder: "VPC"
---

# Table: aws_vpc_block_public_access_options - Query AWS VPC Block Public Access Options using SQL

The AWS VPC Block Public Access Options is a security feature that helps protect your VPC resources from public accessibility. It allows you to control the inbound network connections to your Amazon VPC, preventing unauthorized access. This configuration is crucial for ensuring the security and privacy of your network resources on AWS.

## Table Usage Guide

The `aws_vpc_block_public_access_options` table in Steampipe provides you with information about the block public access configurations for Amazon VPCs. This table allows you, as a DevOps engineer, to query configuration-specific details, including the VPC BPA mode, exclusions allowed status, and management information. You can utilize this table to gather insights on configurations, such as which VPCs have public access blocked, the current state of VPC BPA, and who manages the configuration. The schema outlines the various attributes of the block public access configuration for you, including the internet gateway block mode, exclusions allowed status, and management details.

## Examples
### Basic info
Determine the areas in which public access to your AWS VPCs is blocked, and gain insights into the security configurations and their settings. This can help enhance your understanding of the access control measures in place for your VPC resources.

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

### List VPCs with internet gateway blocking enabled
Identify VPCs that have internet gateway blocking enabled, allowing you to understand which VPCs are configured to block public access through internet gateways. This can be useful in strengthening your security measures and preventing unauthorized access.

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

### List VPCs managed by declarative policies
Discover VPCs that are managed by declarative policies, helping you understand your security setup and identify which VPCs are centrally managed.

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

### List VPCs with bidirectional internet gateway blocking
Explore VPCs that have bidirectional internet gateway blocking enabled, providing maximum security against public access.

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
