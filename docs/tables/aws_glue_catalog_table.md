---
title: "Table: aws_glue_catalog_table - Query AWS Glue Catalog Tables using SQL"
description: "Allows users to query AWS Glue Catalog Tables for a comprehensive overview of table metadata, including table names, database names, owner information, and more."
---

# Table: aws_glue_catalog_table - Query AWS Glue Catalog Tables using SQL

The `aws_glue_catalog_table` table in Steampipe provides information about AWS Glue Catalog Tables. It allows DevOps engineers, data engineers, and other technical professionals to query table-specific details, including table names, database names, owner information, creation time, and associated metadata. Users can utilize this table to gather insights on tables, such as their storage descriptors, partition keys, and table parameters. The schema outlines the various attributes of the AWS Glue Catalog Table, including the catalog ID, database name, table type, storage descriptor, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_catalog_table` table, you can use the `.inspect aws_glue_catalog_table` command in Steampipe.

**Key columns**:

- `name`: This is the name of the table. It is a crucial column as it can be used to join this table with others that also contain table name information.
- `database_name`: This column contains the name of the database where the table is located. This is important for joining with other tables that have database name information.
- `catalog_id`: This column holds the AWS Account ID of the data catalog. It is useful for joining with other tables that contain catalog ID information.

## Examples

### Basic info

```sql
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

```sql
select
  catalog_id,
  count(name) as table_count
from
  aws_glue_catalog_table
group by
  catalog_id;
```

### List tables with retention period less than 30 days

```sql
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
