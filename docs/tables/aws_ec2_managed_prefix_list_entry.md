---
title: "Table: aws_ec2_managed_prefix_list_entry - Query AWS EC2 Managed Prefix List Entry using SQL"
description: "Allows users to query AWS EC2 Managed Prefix List Entries, providing details such as the CIDR block, description, and the prefix list ID. This table is useful for understanding the IP address ranges included in a managed prefix list."
---

# Table: aws_ec2_managed_prefix_list_entry - Query AWS EC2 Managed Prefix List Entry using SQL

The `aws_ec2_managed_prefix_list_entry` table in Steampipe provides information about the IP address ranges, or prefixes, that AWS has added to a managed prefix list. This table allows DevOps engineers to query prefix-specific details, including the CIDR block, description, and the prefix list ID. Users can utilize this table to gather insights on the managed prefix lists, such as the IP address ranges included in a managed prefix list, and more. The schema outlines the various attributes of the managed prefix list entry, including the CIDR, description, and prefix list ID.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_managed_prefix_list_entry` table, you can use the `.inspect aws_ec2_managed_prefix_list_entry` command in Steampipe.

### Key columns:

- `cidr`: This column contains the IP address range, or prefix, in CIDR notation. It is important as it provides the IP address range included in a managed prefix list.
- `description`: This column provides a description of the managed prefix list entry. It is useful for understanding the purpose or use of the IP address range.
- `prefix_list_id`: This column contains the ID of the managed prefix list. It can be used to join this table with other tables that contain details about the managed prefix list.

## Examples

### Basic Info

```sql
select
  prefix_list_id,
  cidr,
  description
from
  aws_ec2_managed_prefix_list_entry;
```

### List customer-managed prefix lists entries

```sql
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

```sql
select
  prefix_list_id,
  count(cidr) as numbers_of_entries
from
  aws_ec2_managed_prefix_list_entry
group by
  prefix_list_id;
```
