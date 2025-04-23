---
title: "Steampipe Table: aws_inspector2_coverage_statistics - Query AWS Inspector2 Coverage Statistics using SQL"
description: "Allows users to query AWS Inspector2 Coverage Statistics to obtain detailed information about the assessment targets and the number of instances they cover."
folder: "Inspector2"
---

# Table: aws_inspector2_coverage_statistics - Query AWS Inspector2 Coverage Statistics using SQL

The AWS Inspector2 Coverage Statistics is a feature of the AWS Inspector service. It provides a detailed overview of your AWS resources and helps identify potential security issues. With this service, you can evaluate the security state of your applications deployed on AWS and improve their security and compliance.

## Table Usage Guide

The `aws_inspector2_coverage_statistics` table in Steampipe provides you with information about AWS Inspector2's coverage statistics. This table allows you as a DevOps engineer, security analyst, or other technical professional to query detailed information about the assessment targets and the number of instances they cover. You can utilize this table to gather insights on assessment targets, including their ARNs, the number of instances they cover, and other associated metadata. The schema outlines the various attributes of the coverage statistics for you, including the assessment target ARN, the instance count, and the agent ID.

## Examples

### Basic info
Determine the areas in which your AWS Inspector service's coverage statistics are distributed. This query can help you understand how your resources are allocated, aiding in efficient resource management.

```sql+postgres
select
  total_counts,
  counts_by_group
from
  aws_inspector2_coverage_statistics;
```

```sql+sqlite
select
  total_counts,
  counts_by_group
from
  aws_inspector2_coverage_statistics;
```

### Get the count of resources within a group
Determine the number of resources within a specific group in AWS Inspector to understand resource distribution and manage resource allocation more effectively.

```sql+postgres
select
  g ->> 'Count' as count,
  g ->> 'GroupKey' as group_key
from
  aws_inspector2_coverage_statistics,
  jsonb_array_elements(counts_by_group) as g;
```

```sql+sqlite
select
  json_extract(g.value, '$.Count') as count,
  json_extract(g.value, '$.GroupKey') as group_key
from
  aws_inspector2_coverage_statistics,
  json_each(counts_by_group) as g;
```