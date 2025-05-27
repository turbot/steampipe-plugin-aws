---
title: "Steampipe Table: aws_redshiftserverless_namespace - Query AWS Redshift Serverless Namespace using SQL"
description: "Allows users to query AWS Redshift Serverless Namespace data. This table provides information about each namespace within an AWS Redshift Serverless cluster. It allows DevOps engineers to query namespace-specific details, including the namespace ARN, creation date, and associated metadata."
folder: "Redshift"
---

# Table: aws_redshiftserverless_namespace - Query AWS Redshift Serverless Namespace using SQL

The AWS Redshift Serverless Namespace is a component of Amazon Redshift, a fully managed, petabyte-scale data warehouse service in the cloud. It allows for the use of SQL, a standard querying language, to interact with and manage the data within the Redshift environment. This serverless feature provides on-demand, scalable capacity, making it easier to set up, operate, and scale a relational database.

## Table Usage Guide

The `aws_redshiftserverless_namespace` table in Steampipe provides you with information about each namespace within an AWS Redshift Serverless cluster. This table allows you, as a DevOps engineer, to query namespace-specific details, including the namespace ARN, creation date, and associated metadata. You can utilize this table to gather insights on namespaces, such as the associated database, the owner of the namespace, and more. The schema outlines the various attributes of the Redshift Serverless Namespace for you, including the namespace ARN, creation date, and associated tags.

## Examples

### Basic info
Explore the basic information of your AWS Redshift Serverless namespaces such as their creation dates, regions, and statuses. This can help you understand the distribution and status of your resources for better resource management.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which AWS Redshift Serverless namespaces are not currently available. This can aid in troubleshooting or planning resource allocation.

```sql+postgres
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

```sql+sqlite
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
  status != 'AVAILABLE';
```

### List all unencrypted namespaces
Identify instances where namespaces are not encrypted in order to enhance security measures and protect sensitive data. This can help in preventing unauthorized access and potential data breaches.

```sql+postgres
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

```sql+sqlite
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
Explore which default IAM roles are associated with each namespace to better understand your AWS Redshift serverless configurations. This can be useful in identifying potential security risks or misconfigurations in your AWS environment.

```sql+postgres
select
  namespace_name,
  namespace_arn,
  namespace_id,
  creation_date,
  default_iam_role_arn
from
  aws_redshiftserverless_namespace;
```

```sql+sqlite
select
  namespace_name,
  namespace_arn,
  namespace_id,
  creation_date,
  default_iam_role_arn
from
  aws_redshiftserverless_namespace;
```