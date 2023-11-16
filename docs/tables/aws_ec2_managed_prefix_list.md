---
title: "Table: aws_ec2_managed_prefix_list - Query AWS EC2 Managed Prefix Lists using SQL"
description: "Allows users to query AWS EC2 Managed Prefix Lists, providing information about IP address ranges (CIDRs), permissions, and associated metadata."
---

# Table: aws_ec2_managed_prefix_list - Query AWS EC2 Managed Prefix Lists using SQL

The `aws_ec2_managed_prefix_list` table in Steampipe provides information about Managed Prefix Lists within AWS EC2. This table allows DevOps engineers to query details about IP address ranges, permissions, and associated metadata. Users can utilize this table to gather insights on IP address ranges, such as which IP addresses are allowed or denied access to a VPC, the maximum number of entries that a prefix list can have, and more. The schema outlines the various attributes of the Managed Prefix List, including the prefix list id, name, owner id, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_managed_prefix_list` table, you can use the `.inspect aws_ec2_managed_prefix_list` command in Steampipe.

**Key columns**:

- `prefix_list_id`: The ID of the managed prefix list. This can be used to join with other tables that reference prefix lists.
- `owner_id`: The AWS account ID of the owner of the managed prefix list. This can be used to join with other tables that reference AWS account ownership.
- `entries`: The entries (CIDR blocks) for the prefix list. This can be used to join with other tables that reference IP address ranges.

## Examples

### Basic Info

```sql
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list;
```

### List customer-managed prefix lists

```sql
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  owner_id <> 'AWS';
```

### List prefix lists with IPv6 as IP address version

```sql
select
  name,
  id,
  address_family
from
  aws_ec2_managed_prefix_list
where
  address_family = 'IPv6';
```

### List prefix lists by specific IDs

```sql
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  id in ('pl-03a3e735e3467c0c4', 'pl-4ca54025');
```

### List prefix lists by specific names

```sql
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  name in ('testPrefix', 'com.amazonaws.us-east-2.dynamodb');
```

### List prefix lists by a specific owner ID

```sql
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  owner_id = '632901234528';
```
