---
title: "Steampipe Table: aws_ec2_managed_prefix_list_entry - Query AWS EC2 Managed Prefix List Entry using SQL"
description: "Allows users to query AWS EC2 Managed Prefix List Entries, providing details such as the CIDR block, description, and the prefix list ID. This table is useful for understanding the IP address ranges included in a managed prefix list."
folder: "EC2"
---

# Table: aws_ec2_managed_prefix_list_entry - Query AWS EC2 Managed Prefix List Entry using SQL

The AWS EC2 Managed Prefix List Entry is a part of Amazon Elastic Compute Cloud (EC2) service. It helps you to manage IP address ranges, allowing you to create lists of IP address ranges, known as prefix lists, and use them to simplify the configuration of security groups and route tables. This makes it easier to set up, secure, and manage the network access to your Amazon EC2 instances.

## Table Usage Guide

The `aws_ec2_managed_prefix_list_entry` table in Steampipe provides you with information about the IP address ranges, or prefixes, that AWS has added to a managed prefix list. This table allows you, as a DevOps engineer, to query prefix-specific details, including the CIDR block, description, and the prefix list ID. You can utilize this table to gather insights on the managed prefix lists, such as the IP address ranges included in a managed prefix list, and more. The schema outlines for you the various attributes of the managed prefix list entry, including the CIDR, description, and prefix list ID.

## Examples

### Basic Info
Explore which AWS EC2 managed prefix list entries exist in your environment. This can help you determine if there are any unexpected or unnecessary entries that may need to be addressed for security or efficiency reasons.

```sql+postgres
select
  prefix_list_id,
  cidr,
  description
from
  aws_ec2_managed_prefix_list_entry;
```

```sql+sqlite
select
  prefix_list_id,
  cidr,
  description
from
  aws_ec2_managed_prefix_list_entry;
```

### List customer-managed prefix lists entries
Explore which customer-managed prefix lists entries are owned by entities other than AWS. This can be useful to understand the distribution and ownership of these resources, helping you to manage and control access to your network resources.

```sql+postgres
select
  l.name,
  l.id,
  e.cidr,
  e.description,
  l.state,
  l.owner_id
from
  aws_ec2_managed_prefix_list_entry as e,
  aws_ec2_managed_prefix_list as l
where
  l.owner_id <> 'AWS';
```

```sql+sqlite
select
  l.name,
  l.id,
  e.cidr,
  e.description,
  l.state,
  l.owner_id
from
  aws_ec2_managed_prefix_list_entry as e,
  aws_ec2_managed_prefix_list as l
where
  l.owner_id <> 'AWS';
```

### Count prefix list entries by prefix list
Discover the segments that have varying numbers of entries in AWS EC2 managed prefix lists, providing a useful summary of the distribution of entries across different lists. This can assist in identifying any disproportionate allocation of entries which may require rebalancing.

```sql+postgres
select
  prefix_list_id,
  count(cidr) as numbers_of_entries
from
  aws_ec2_managed_prefix_list_entry
group by
  prefix_list_id;
```

```sql+sqlite
select
  prefix_list_id,
  count(cidr) as numbers_of_entries
from
  aws_ec2_managed_prefix_list_entry
group by
  prefix_list_id;
```