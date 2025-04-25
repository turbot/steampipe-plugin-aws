---
title: "Steampipe Table: aws_emr_block_public_access_configuration - Query AWS EMR Block Public Access Configuration using SQL"
description: "Allows users to query AWS EMR Block Public Access Configuration to retrieve details about the block public access configurations for EMR clusters."
folder: "Config"
---

# Table: aws_emr_block_public_access_configuration - Query AWS EMR Block Public Access Configuration using SQL

The AWS EMR Block Public Access Configuration is a security feature that helps protect your EMR resources from public accessibility. It allows you to control the inbound network connections to your Amazon EMR clusters, preventing unauthorized access. This configuration is crucial for ensuring the security and privacy of your data processing and analytical tasks on AWS EMR.

## Table Usage Guide

The `aws_emr_block_public_access_configuration` table in Steampipe provides you with information about the block public access configurations for Amazon EMR clusters. This table allows you, as a DevOps engineer, to query configuration-specific details, including the block public access status, permitted public security group rules, and associated metadata. You can utilize this table to gather insights on configurations, such as the number of permitted public security group rules, the block public access status, and more. The schema outlines the various attributes of the block public access configuration for you, including the block public access status, the permitted public security group rules, and the metadata.

## Examples

### Basic info
Determine the areas in which public access to your AWS Elastic MapReduce (EMR) clusters is blocked, and gain insights into the security group rules and their creation dates. This can help enhance your understanding of the access control measures in place for your EMR resources.

```sql+postgres
select
  created_by_arn,
  block_public_security_group_rules,
  creation_date,
  classification,
  permitted_public_security_group_rule_ranges
from
  aws_emr_block_public_access_configuration
order by
  created_by_arn,
  creation_date;
```

```sql+sqlite
select
  created_by_arn,
  block_public_security_group_rules,
  creation_date,
  classification,
  permitted_public_security_group_rule_ranges
from
  aws_emr_block_public_access_configuration
order by
  created_by_arn,
  creation_date;
```

### List block public access configurations that block public security group rules
Identify configurations that are set to block public security group rules, allowing you to understand which elements in your AWS EMR block public access settings are preventing public access. This can be useful in strengthening your security measures and preventing unauthorized access.

```sql+postgres
select
  created_by_arn,
  creation_date
from
  aws_emr_block_public_access_configuration
where
  block_public_security_group_rules;
```

```sql+sqlite
select
  created_by_arn,
  creation_date
from
  aws_emr_block_public_access_configuration
where
  block_public_security_group_rules = 1;
```

### List permitted public security group rule maximum and minimum port ranges
Discover the segments that have the maximum and minimum port ranges allowed by public security group rules. This can be useful for understanding your security setup and identifying potential vulnerabilities.

```sql+postgres
select
  created_by_arn,
  creation_date,
  rules ->> 'MaxRange' as max_range,
  rules ->> 'MinRange' as min_range
from
  aws_emr_block_public_access_configuration
  cross join jsonb_array_elements(permitted_public_security_group_rule_ranges) as rules;
```

```sql+sqlite
select
  created_by_arn,
  creation_date,
  json_extract(rules.value, '$.MaxRange') as max_range,
  json_extract(rules.value, '$.MinRange') as min_range
from
  aws_emr_block_public_access_configuration,
  json_each(permitted_public_security_group_rule_ranges) as rules;
```

### List block public access configurations created in last 90 days
Explore the recent configurations that block public access, specifically those set up within the last 90 days. This can help maintain security by ensuring that public access restrictions are up-to-date and relevant.

```sql+postgres
select
  created_by_arn,
  creation_date
from
  aws_emr_block_public_access_configuration
where
  date_part('day', now() - creation_date) < 90;
```

```sql+sqlite
select
  created_by_arn,
  creation_date
from
  aws_emr_block_public_access_configuration
where
  julianday('now') - julianday(creation_date) < 90;
```