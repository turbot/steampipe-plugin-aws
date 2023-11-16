---
title: "Table: aws_emr_block_public_access_configuration - Query AWS EMR Block Public Access Configuration using SQL"
description: "Allows users to query AWS EMR Block Public Access Configuration to retrieve details about the block public access configurations for EMR clusters."
---

# Table: aws_emr_block_public_access_configuration - Query AWS EMR Block Public Access Configuration using SQL

The `aws_emr_block_public_access_configuration` table in Steampipe provides information about the block public access configurations for Amazon EMR clusters. This table allows DevOps engineers to query configuration-specific details, including the block public access status, permitted public security group rules, and associated metadata. Users can utilize this table to gather insights on configurations, such as the number of permitted public security group rules, the block public access status, and more. The schema outlines the various attributes of the block public access configuration, including the block public access status, the permitted public security group rules, and the metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_emr_block_public_access_configuration` table, you can use the `.inspect aws_emr_block_public_access_configuration` command in Steampipe.

### Key columns:

- `region`: The AWS Region of the resource. This column can be used to join with other tables to get region-specific information.
- `block_public_access_enabled`: Indicates whether block public access is enabled. This can be useful in understanding the security posture of your EMR clusters.
- `permitted_public_security_group_rule_count`: The count of rules that are allowed for the public security groups. This can be important in managing and controlling access to your EMR clusters.

## Examples

### Basic info

```sql
select
  created_by_arn,
  block_public_security_group_rules,
  creation_date,
  classification,
  configurations,
  permitted_public_security_group_rule_ranges
from
  aws_emr_block_public_access_configuration
order by
  created_by_arn,
  creation_date;
```

### List block public access configurations that block public security group rules

```sql
select
  created_by_arn,
  creation_date,
  configurations
from
  aws_emr_block_public_access_configuration
where
  block_public_security_group_rules;
```

### List permitted public security group rule maximum and minimum port ranges

```sql
select
  created_by_arn,
  creation_date,
  rules ->> 'MaxRange' as max_range,
  rules ->> 'MinRange' as min_range
from
  aws_emr_block_public_access_configuration
  cross join jsonb_array_elements(permitted_public_security_group_rule_ranges) as rules;
```

### List block public access configurations created in last 90 days

```sql
select
  created_by_arn,
  creation_date,
  configurations
from
  aws_emr_block_public_access_configuration
where
  date_part('day', now() - creation_date) < 90;
```
