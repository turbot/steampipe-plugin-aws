---
title: "Table: aws_glue_catalog_database - Query AWS Glue Databases using SQL"
description: "Allows users to query AWS Glue Databases for detailed information about their Glue Catalog Databases."
---

# Table: aws_glue_catalog_database - Query AWS Glue Databases using SQL

The `aws_glue_catalog_database` table in Steampipe provides information about databases within AWS Glue Catalog. This table allows DevOps engineers, data scientists, and database administrators to query database-specific details, including the catalog ID, database name, description, location URI, and associated metadata. Users can utilize this table to gather insights on databases, such as the creation time, last modified time, and the number of tables in each database. The schema outlines the various attributes of the AWS Glue Catalog Database, including the create time, compatibility, data location, parameters, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_catalog_database` table, you can use the `.inspect aws_glue_catalog_database` command in Steampipe.

**Key columns**:

- `name`: The name of the database. This can be used to join with other tables that reference the database name.
- `catalog_id`: The ID of the data catalog in which the database resides. This can be used to join with other tables that reference the catalog ID.
- `create_time`: The time at which the database was created. This can be useful for understanding the lifecycle and age of your databases.

## Examples

### Basic info

```sql
select
  name,
  catalog_id,
  create_time,
  description,
  location_uri,
  create_table_default_permissions
from
  aws_glue_catalog_database;
```

### Count the number of databases per catalog

```sql
select
  catalog_id,
  count(name) as database_count
from
  aws_glue_catalog_database
group by
  catalog_id;
```
