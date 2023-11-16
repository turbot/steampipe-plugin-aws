---
title: "Table: aws_accessanalyzer_analyzer - Query AWS Access Analyzer using SQL"
description: "Allows users to query Access Analyzer Analyzer in AWS IAM to retrieve information about analyzers."
---

# Table: aws_accessanalyzer_analyzer - Query AWS Access Analyzer using SQL

The `aws_accessanalyzer_analyzer` table in Steampipe provides information about analyzers within AWS IAM Access Analyzer. This table allows DevOps engineers to query analyzer-specific details, including the analyzer ARN, type, status, and associated metadata. Users can utilize this table to gather insights on analyzers, such as the status of each analyzer, the type of analyzer, and the resource that was analyzed. The schema outlines the various attributes of the Access Analyzer, including the analyzer ARN, creation time, last resource scanned, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_accessanalyzer_analyzer` table, you can use the `.inspect aws_accessanalyzer_analyzer` command in Steampipe.

### Key columns:

- `arn`: The ARN (Amazon Resource Name) of the analyzer. This can be used to join with other tables where the analyzer ARN is required.
- `name`: The name of the analyzer. This can be used to join with other tables where the analyzer name is required.
- `type`: The type of analyzer. This can be useful for filtering the data based on the analyzer type.

## Examples

### Basic info

```sql
select
  name,
  last_resource_analyzed,
  last_resource_analyzed_at,
  status,
  type
from
  aws_accessanalyzer_analyzer;
```

### List analyzers which are enabled

```sql
select
  name,
  status
  last_resource_analyzed,
  last_resource_analyzed_at,
  tags
from
  aws_accessanalyzer_analyzer
where
  status = 'ACTIVE';
```


### List analyzers with findings that need to be resolved

```sql
select
  name,
  status,
  type,
  last_resource_analyzed
from
  aws_accessanalyzer_analyzer
where
  status = 'ACTIVE'
  and findings is not null;
```
