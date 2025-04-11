---
title: "Steampipe Table: aws_lakeformation_permission - Query AWS Lake Formation Permissions Using SQL"
description: "Query AWS Lake Formation permissions, including granted permissions, principals, and associated AWS resources."
folder: "Lake Formation"
---

# Table: aws_lakeformation_permission - Query AWS Lake Formation Permissions Using SQL

The `aws_lakeformation_permission` table allows you to query **AWS Lake Formation permissions**, providing insights into **who has access** to data lake resources, what **permissions** they hold, and **which AWS resources** those permissions apply to. This table helps **data governance teams and security administrators** monitor **data access controls** in AWS Lake Formation.

## Table Usage Guide

The `aws_lakeformation_permission` table provides detailed information about **permissions granted in AWS Lake Formation**, including the **principals** (IAM users, roles, or AWS services) who have permissions, the **resource types** such as databases, tables, columns, LF-tags, and data locations, and the **types of permissions** granted, including `SELECT`, `DESCRIBE`, and `ALL`. It also includes details on **grant options**, indicating whether the principal can delegate permissions to others, and **last updated timestamps** to track permission changes.

## Examples

### List all AWS Lake Formation permissions
Retrieve a list of all **granted permissions**, including the **principal** and **resource type**.

```sql+postgres
select
  principal_identifier,
  database_name,
  table_name,
  permissions
from
  aws_lakeformation_permission;
```

```sql+sqlite
select
  principal_identifier,
  database_name,
  table_name,
  permissions
from
  aws_lakeformation_permission;
```

### Find permissions granted to a specific IAM Role
Identify all Lake Formation permissions granted to an IAM role.

```sql+postgres
select
  principal_identifier,
  resource_catalog_id,
  database_name,
  table_name,
  permissions
from
  aws_lakeformation_permission
where
  principal_identifier = 'arn:aws:iam::123456789012:role/MyLakeFormationRole';
```

```sql+sqlite
select
  principal_identifier,
  resource_catalog_id,
  database_name,
  table_name,
  permissions
from
  aws_lakeformation_permission
where
  principal_identifier = 'arn:aws:iam::123456789012:role/MyLakeFormationRole';
```

### Find permissions granted on a specific database
Retrieve all permissions on a particular database.

```sql+postgres
select
  principal_identifier,
  permissions,
  last_updated
from
  aws_lakeformation_permission
where
  database_name = 'my_database';
```

```sql+sqlite
select
  principal_identifier,
  permissions,
  last_updated
from
  aws_lakeformation_permission
where
  database_name = 'my_database';
```

### List permissions with grant options enabled
Identify permissions where a principal can grant access to others.

```sql+postgres
select
  principal_identifier,
  database_name,
  table_name,
  permissions_with_grant_option
from
  aws_lakeformation_permission
where
  jsonb_array_length(permissions_with_grant_option) > 0;
```

```sql+sqlite
select
  principal_identifier,
  database_name,
  table_name,
  permissions_with_grant_option
from
  aws_lakeformation_permission
where
  jsonb_array_length(permissions_with_grant_option) > 0;
```

### Get permissions associated with LF-tags
Retrieve Lake Formation tag-based access permissions.

```sql+postgres
select
  lf_tag_key,
  lf_tag_values,
  principal_identifier,
  permissions
from
  aws_lakeformation_permission
where
  lf_tag_key is not null;
```

```sql+sqlite
select
  lf_tag_key,
  lf_tag_values,
  principal_identifier,
  permissions
from
  aws_lakeformation_permission
where
  lf_tag_key is not null;
```

### Find IAM principals with access to LF-tags
This query identifies who has permissions on LF-tagged resources, helping enforce tag-based access control (ABAC).

```sql+postgres
select
  p.principal_identifier,
  t.tag_key,
  t.tag_values,
  p.permissions
from
  aws_lakeformation_permission p
  join aws_lakeformation_tag t on p.lf_tag_key = t.tag_key;
```

```sql+sqlite
select
  p.principal_identifier,
  t.tag_key,
  t.tag_values,
  p.permissions
from
  aws_lakeformation_permission p
  join aws_lakeformation_tag t on p.lf_tag_key = t.tag_key;
```

### Find IAM principals with access to registered data locations
This query identifies which IAM principals have access to registered S3 locations in the data lake.

```sql+postgres
select
  p.principal_identifier,
  r.resource_arn as s3_bucket,
  p.permissions
from
  aws_lakeformation_permission p
  join aws_lakeformation_resource r on p.data_location_resource_arn = r.resource_arn;
```

```sql+sqlite
select
  p.principal_identifier,
  r.resource_arn as s3_bucket,
  p.permissions
from
  aws_lakeformation_permission p
  join aws_lakeformation_resource r on p.data_location_resource_arn = r.resource_arn;
```