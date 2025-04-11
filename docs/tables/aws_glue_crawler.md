---
title: "Steampipe Table: aws_glue_crawler - Query AWS Glue Crawlers using SQL"
description: "Allows users to query AWS Glue Crawlers and retrieve essential information about the crawler's configuration, status, and associated metadata."
folder: "Glue"
---

# Table: aws_glue_crawler - Query AWS Glue Crawlers using SQL

The AWS Glue Crawler is a component of AWS Glue service that automates the extraction, transformation, and loading (ETL) process. It traverses your data stores, identifies data formats, and suggests schemas and transformations. This enables you to categorize, search, and query metadata across your AWS environment.

## Table Usage Guide

The `aws_glue_crawler` table in Steampipe provides you with information about crawlers within AWS Glue. This table allows you, as a DevOps engineer, to query crawler-specific details, including its role, database, schedule, classifiers, and associated metadata. You can utilize this table to gather insights on crawlers, such as their run frequency, the database they are associated with, their status, and more. The schema outlines the various attributes of the Glue crawler for you, including the crawler ARN, creation date, last run time, and associated tags.

## Examples

### Basic info
Determine the status and creation details of your AWS Glue crawlers to better understand their function and manage them effectively. This can be particularly useful for identifying any crawlers that may require attention or modification.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are currently operational within your AWS Glue Crawlers to understand which tasks are active and could be consuming resources. This could be useful for resource management and troubleshooting ongoing tasks.

```sql+postgres
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

```sql+sqlite
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