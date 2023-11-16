---
title: "Table: aws_glue_crawler - Query AWS Glue Crawlers using SQL"
description: "Allows users to query AWS Glue Crawlers and retrieve essential information about the crawler's configuration, status, and associated metadata."
---

# Table: aws_glue_crawler - Query AWS Glue Crawlers using SQL

The `aws_glue_crawler` table in Steampipe provides information about crawlers within AWS Glue. This table allows DevOps engineers to query crawler-specific details, including its role, database, schedule, classifiers, and associated metadata. Users can utilize this table to gather insights on crawlers, such as their run frequency, the database they are associated with, their status, and more. The schema outlines the various attributes of the Glue crawler, including the crawler ARN, creation date, last run time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_crawler` table, you can use the `.inspect aws_glue_crawler` command in Steampipe.

**Key columns**:

- `crawler_name`: The name of the crawler. This is a unique identifier for the crawler and can be used to join with other tables.
- `role`: The IAM role (or Amazon Resource Name - ARN) associated with the crawler. This provides information about the permissions assigned to the crawler.
- `database_name`: The name of the database in which the crawler metadata resides. This can be used to join with other tables that reference the same database.

## Examples

### Basic info

```sql
select
  name,
  state,
  database_name,
  creation_time,
  description,
  recrawl_behavior
from
  aws_glue_crawler;
```

### List running crawlers

```sql
select
  name,
  state,
  database_name,
  creation_time,
  description,
  recrawl_behavior
from
  aws_glue_crawler
where
  state = 'RUNNING'; 
```