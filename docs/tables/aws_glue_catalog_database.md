---
title: "Steampipe Table: aws_glue_catalog_database - Query AWS Glue Databases using SQL"
description: "Allows users to query AWS Glue Databases for detailed information about their Glue Catalog Databases."
folder: "Glue"
---

# Table: aws_glue_catalog_database - Query AWS Glue Databases using SQL

The AWS Glue Catalog Database is a managed service that serves as your integrated, centralized data catalog. It organizes, locates, moves, controls, and cleans data across various data stores. It also stores metadata about databases, tables, or other data catalog objects created by AWS Glue.

## Table Usage Guide

The `aws_glue_catalog_database` table in Steampipe provides you with information about databases within AWS Glue Catalog. This table allows you, as a DevOps engineer, data scientist, or database administrator, to query database-specific details, including the catalog ID, database name, description, location URI, and associated metadata. You can utilize this table to gather insights on databases, such as the creation time, last modified time, and the number of tables in each database. The schema outlines for you the various attributes of the AWS Glue Catalog Database, including the create time, compatibility, data location, parameters, and associated tags.

## Examples

### Basic info
Determine the areas in which AWS Glue has cataloged databases, including when they were created and their default permissions. This can be useful for understanding your data landscape and ensuring appropriate access controls are in place.

```sql+postgres
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

```sql+sqlite
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
Determine the volume of databases associated with each catalog in your AWS Glue service. This can help you understand how your databases are distributed across different catalogs, aiding in efficient resource allocation and management.

```sql+postgres
select
  catalog_id,
  count(name) as database_count
from
  aws_glue_catalog_database
group by
  catalog_id;
```

```sql+sqlite
select
  catalog_id,
  count(name) as database_count
from
  aws_glue_catalog_database
group by
  catalog_id;
```