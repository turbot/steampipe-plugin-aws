---
title: "Steampipe Table: aws_ec2_managed_prefix_list - Query AWS EC2 Managed Prefix Lists using SQL"
description: "Allows users to query AWS EC2 Managed Prefix Lists, providing information about IP address ranges (CIDRs), permissions, and associated metadata."
folder: "EC2"
---

# Table: aws_ec2_managed_prefix_list - Query AWS EC2 Managed Prefix Lists using SQL

The AWS EC2 Managed Prefix List is a resource that allows you to create and manage prefix lists for your AWS account. These prefix lists are used to group IP address ranges and simplify the configuration of security group rules and route table entries. They are especially useful in managing large IP address ranges and maintaining security in your AWS environment.

There are two types of prefix lists:

* **Customer-managed prefix lists** - Sets of IP address ranges that you define and manage. You can share your prefix list with other AWS accounts, enabling those accounts to reference the prefix list in their own resources.
* **AWS-managed prefix lists** - Sets of IP address ranges for AWS services. You cannot create, modify, share, or delete an AWS-managed prefix list.

## Table Usage Guide

The `aws_ec2_managed_prefix_list` table in Steampipe provides you with information about Managed Prefix Lists within AWS EC2. This table allows you as a DevOps engineer to query details about IP address ranges, permissions, and associated metadata. You can utilize this table to gather insights on IP address ranges, such as which IP addresses are allowed or denied access to a VPC, the maximum number of entries that a prefix list can have, and more. The schema outlines the various attributes of the Managed Prefix List for you, including the prefix list id, name, owner id, and associated tags.

## Examples

### Basic Info
Explore the ownership and status of your managed prefix lists in AWS EC2. This can help you understand who controls these resources and their current operational state.

```sql+postgres
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list;
```

```sql+sqlite
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
Explore which customer-managed prefix lists are in use to gain insights into your AWS EC2 configurations. This helps identify any potential security risks or configuration issues.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  owner_id != 'AWS';
```

### List prefix lists with IPv6 as IP address version
Determine the areas in which IPv6 is used as the IP address version within your managed prefix lists. This is useful for understanding your network's IPv6 usage and ensuring compatibility with IPv6-only systems.

```sql+postgres
select
  name,
  id,
  address_family
from
  aws_ec2_managed_prefix_list
where
  address_family = 'IPv6';
```

```sql+sqlite
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
Determine the areas in which specific AWS EC2 managed prefix lists are being used by identifying them through their unique IDs. This query is beneficial in managing and tracking the usage of prefix lists in your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which specific managed prefix lists are used within the AWS EC2 service. This can be beneficial for understanding the configuration and usage of these lists in your cloud environment.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which specific AWS EC2 managed prefix lists are owned by a particular user. This is useful for understanding the distribution and ownership of these resources, helping to manage and organize your AWS environment effectively.

```sql+postgres
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

```sql+sqlite
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