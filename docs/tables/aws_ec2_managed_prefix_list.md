# Table: aws_ec2_managed_prefix_list

A prefix list is a set of one or more CIDR blocks. You can use prefix lists to make it easier to configure and maintain your security groups and route tables. You can create a prefix list from the IP addresses that you frequently use, and reference them as a set in security group rules and routes instead of referencing them individually.

There are two types of prefix lists:

* **Customer-managed prefix lists** — Sets of IP address ranges that you define and manage. You can share your prefix list with other AWS accounts, enabling those accounts to reference the prefix list in their own resources.
* **AWS-managed prefix lists** — Sets of IP address ranges for AWS services. You cannot create, modify, share, or delete an AWS-managed prefix list.

## Examples

### Basic Info

```sql
select
  prefix_list_name,
  prefix_list_id,
  prefix_list_arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list;
```

### List customer-managed prefix lists

```sql
select
  prefix_list_name,
  prefix_list_id,
  prefix_list_arn,
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
  prefix_list_name,
  prefix_list_id,
  address_family
from
  aws_ec2_managed_prefix_list
where
  address_family = 'IPv6';
```

### List prefix lists by specific IDs

```sql
select
  prefix_list_name,
  prefix_list_id,
  prefix_list_arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  prefix_list_id in ('pl-0a81b1c04d3d1485f', 'pl-4ca54025');
```

### List prefix lists by specific names

```sql
select
  prefix_list_name,
  prefix_list_id,
  prefix_list_arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  prefix_list_name in ('testPrefix', 'com.amazonaws.us-east-2.dynamodb');
```

### List prefix lists by a specific owner ID

```sql
select
  prefix_list_name,
  prefix_list_id,
  prefix_list_arn,
  state,
  owner_id
from
  aws_ec2_managed_prefix_list
where
  owner_id = '632901234528';
```
