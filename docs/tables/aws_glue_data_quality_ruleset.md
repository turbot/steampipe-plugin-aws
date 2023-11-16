---
title: "Table: aws_glue_data_quality_ruleset - Query AWS Glue Data Quality Ruleset using SQL"
description: "Allows users to query AWS Glue Data Quality Ruleset to obtain information about the rulesets used for data quality checks in AWS Glue."
---

# Table: aws_glue_data_quality_ruleset - Query AWS Glue Data Quality Ruleset using SQL

The `aws_glue_data_quality_ruleset` table in Steampipe provides information about the rulesets used for data quality checks in AWS Glue. This table allows data engineers and developers to query ruleset-specific details, including the ruleset name, status, related applications, and associated metadata. Users can utilize this table to gather insights on rulesets, such as ruleset usage, associated applications, status of rulesets, and more. The schema outlines the various attributes of the data quality ruleset, including the ruleset ARN, creation date, last modified date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_data_quality_ruleset` table, you can use the `.inspect aws_glue_data_quality_ruleset` command in Steampipe.

### Key columns:

- `name`: The name of the ruleset. This can be used to join this table with others that also contain ruleset names.
- `arn`: The Amazon Resource Number (ARN) of the ruleset. This is a unique identifier that can be used to join this table with others that contain ARN information.
- `status`: The status of the ruleset. This can be useful for joining this table with others to get a comprehensive view of active and inactive rulesets.

## Examples

### Basic info

```sql
select
  name,
  database_name,
  table_name,
  created_on,
  description,
  rule_set,
  recommendation_run_id
from
  aws_glue_data_quality_ruleset;
```

### List rulesets created in the last 30 days

```sql
select
  name,
  database_name,
  table_name,
  created_on,
  description,
  rule_set,
  recommendation_run_id
from
  aws_glue_data_quality_ruleset
where
  created_on >= now() - interval '30' day;
```

### Count ruleset by database

```sql
select
  database_name,
  count("name") as rulset_count
from
  aws_glue_data_quality_ruleset
group by
  database_name;
```

### Get Glue database details for a ruleset

```sql
select
  r.name,
  r.database_name,
  d.catalog_id,
  d.create_time as databse_create_time,
  d.location_uri
from
  aws_glue_data_quality_ruleset as r,
  aws_glue_catalog_database as d
where
  r.database_name = d.name
and
  r.name = 'ruleset1';
```


### Count rules per data quality ruleset

```sql
select
  name,
  rule_count
from
  aws_glue_data_quality_ruleset;
```