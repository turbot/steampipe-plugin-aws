---
title: "Steampipe Table: aws_glue_data_quality_ruleset - Query AWS Glue Data Quality Ruleset using SQL"
description: "Allows users to query AWS Glue Data Quality Ruleset to obtain information about the rulesets used for data quality checks in AWS Glue."
folder: "Glue"
---

# Table: aws_glue_data_quality_ruleset - Query AWS Glue Data Quality Ruleset using SQL

The AWS Glue Data Quality Ruleset is a feature of AWS Glue that enables you to enforce quality rules on your data sources. It allows you to define, manage, and run data quality rules on your AWS Glue Data Catalog tables. This feature helps ensure that your data is accurate, consistent, and reliable, thereby improving the overall quality of your data.

## Table Usage Guide

The `aws_glue_data_quality_ruleset` table in Steampipe provides you with information about the rulesets used for data quality checks in AWS Glue. This table allows you as a data engineer or developer to query ruleset-specific details, including the ruleset name, status, related applications, and associated metadata. You can utilize this table to gather insights on rulesets, such as ruleset usage, associated applications, status of rulesets, and more. The schema outlines the various attributes of the data quality ruleset for you, including the ruleset ARN, creation date, last modified date, and associated tags.

## Examples

### Basic info
Explore the creation dates and descriptions of various data quality rulesets in AWS Glue. This can help in understanding the evolution of data quality standards and guidelines over time in your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which rulesets have been created in the past month, providing a recent history of data quality ruleset generation. This can be useful for monitoring the frequency and timing of new ruleset creation.

```sql+postgres
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

```sql+sqlite
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
  created_on >= datetime('now','-30 day');
```

### Count ruleset by database
Explore which databases have the most rulesets in order to optimize data quality checks. This insight can help prioritize which databases need more attention or resources for maintaining data quality.

```sql+postgres
select
  database_name,
  count("name") as rulset_count
from
  aws_glue_data_quality_ruleset
group by
  database_name;
```

```sql+sqlite
select
  database_name,
  count("name") as rulset_count
from
  aws_glue_data_quality_ruleset
group by
  database_name;
```

### Get Glue database details for a ruleset
Analyze the settings to understand the specific details of a Glue database associated with a certain data quality ruleset. This can be particularly useful for auditing or troubleshooting purposes, allowing you to pinpoint specific locations and creation times of the database.

```sql+postgres
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

```sql+sqlite
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
Determine the number of rules within each data quality ruleset to assess the complexity and thoroughness of your data validation process.

```sql+postgres
select
  name,
  rule_count
from
  aws_glue_data_quality_ruleset;
```

```sql+sqlite
select
  name,
  rule_count
from
  aws_glue_data_quality_ruleset;
```