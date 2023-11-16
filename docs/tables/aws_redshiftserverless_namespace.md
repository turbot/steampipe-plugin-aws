---
title: "Table: aws_redshiftserverless_namespace - Query AWS Redshift Serverless Namespace using SQL"
description: "Allows users to query AWS Redshift Serverless Namespace data. This table provides information about each namespace within an AWS Redshift Serverless cluster. It allows DevOps engineers to query namespace-specific details, including the namespace ARN, creation date, and associated metadata."
---

# Table: aws_redshiftserverless_namespace - Query AWS Redshift Serverless Namespace using SQL

The `aws_redshiftserverless_namespace` table in Steampipe provides information about each namespace within an AWS Redshift Serverless cluster. This table allows DevOps engineers to query namespace-specific details, including the namespace ARN, creation date, and associated metadata. Users can utilize this table to gather insights on namespaces, such as the associated database, the owner of the namespace, and more. The schema outlines the various attributes of the Redshift Serverless Namespace, including the namespace ARN, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshiftserverless_namespace` table, you can use the `.inspect aws_redshiftserverless_namespace` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the namespace. This unique identifier is important for joining this table with other AWS tables.
- `database_name`: The name of the database associated with the namespace. This can be used to join with tables that also contain database information.
- `owner`: The owner of the namespace. This can be useful for identifying who has control over the namespace and can be used to join with tables that contain user information.

## Examples

### Basic info

```sql
select
  namespace_name,
  namespace_arn,
  namespace_id,
  creation_date,
  db_name,
  region,
  status
from
  aws_redshiftserverless_namespace;
```

### List all unavailable namespaces

```sql
select
  namespace_name,
  namespace_arn,
  namespace_id,
  creation_date,
  db_name,
  region,
  status
from
  aws_redshiftserverless_namespace
where
  status <> 'AVAILABLE';
```

### List all unencrypted namespaces

```sql
select
  namespace_name,
  namespace_arn,
  namespace_id,
  creation_date,
  db_name,
  region,
  status
from
  aws_redshiftserverless_namespace
where
  kms_key_id is null;
```

### Get default IAM role ARN associated with each namespace

```sql
select
  namespace_name,
  namespace_arn,
  namespace_id,
  creation_date,
  default_iam_role_arn
from
  aws_redshiftserverless_namespace;
```
