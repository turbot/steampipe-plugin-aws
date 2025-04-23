---
title: "Steampipe Table: aws_glue_catalog_table - Query AWS Glue Catalog Tables using SQL"
description: "Allows users to query AWS Glue Catalog Tables for a comprehensive overview of table metadata, including table names, database names, owner information, and more."
folder: "Glue"
---

# Table: aws_glue_catalog_table - Query AWS Glue Catalog Tables using SQL

The AWS Glue Catalog Tables are a part of AWS Glue, a fully managed extract, transform, and load (ETL) service that makes it easy for users to prepare and load their data for analytics. AWS Glue Catalog Tables store metadata related to data sources, transformations, and targets, allowing users to discover and manage their data. AWS Glue automatically generates the schema for your data and stores it as tables in its data catalog, which you can use for ETL jobs.

## Table Usage Guide

The `aws_glue_catalog_table` table in Steampipe provides you with information about AWS Glue Catalog Tables. It allows you, whether you're a DevOps engineer, data engineer, or another technical professional, to query table-specific details, including table names, database names, owner information, creation time, and associated metadata. You can utilize this table to gather insights on tables, such as their storage descriptors, partition keys, and table parameters. The schema outlines the various attributes of the AWS Glue Catalog Table for you, including the catalog ID, database name, table type, storage descriptor, and associated tags.

## Examples

### Basic info
Explore the basic information of your AWS Glue Catalog Tables to understand when and why they were created. This can help in managing your resources and planning future database schemas.

```sql+postgres
select
  name,
  catalog_id,
  create_time,
  description,
  database_name
from
  aws_glue_catalog_table;
```

```sql+sqlite
select
  name,
  catalog_id,
  create_time,
  description,
  database_name
from
  aws_glue_catalog_table;
```

### Count the number of tables per catalog
Analyze the distribution of tables across different catalogs in AWS Glue service. This helps in understanding the organization of your data assets and can assist in optimizing data management strategies.

```sql+postgres
select
  catalog_id,
  count(name) as table_count
from
  aws_glue_catalog_table
group by
  catalog_id;
```

```sql+sqlite
select
  catalog_id,
  count(name) as table_count
from
  aws_glue_catalog_table
group by
  catalog_id;
```

### List tables with retention period less than 30 days
Explore which AWS Glue tables have a retention period of less than 30 days. This can be useful in identifying tables that may require a more extended retention period for data backup and disaster recovery purposes.

```sql+postgres
select
  name,
  catalog_id,
  create_time,
  description,
  retention
from
  aws_glue_catalog_table
where
  retention < 30;
```

```sql+sqlite
select
  name,
  catalog_id,
  create_time,
  description,
  retention
from
  aws_glue_catalog_table
where
  retention < 30;
```