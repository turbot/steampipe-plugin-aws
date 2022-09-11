# Table: aws_redshiftserverless_namespace

**Amazon Redshift Serverless** makes it convenient for you to run and scale analytics without having to provision and manage data warehouses. With Amazon Redshift Serverless, data analysts, developers, and data scientists can now use Amazon Redshift to get insights from data in seconds by loading data into and querying records from the data warehouse.

**Namespace** is a collection of database objects and users. The storage-related namespace groups together schemas, tables, users, or AWS Key Management Service keys for encrypting data. Storage properties include the database name and password of the admin user, permissions, and encryption and security. Other resources that are grouped under namespaces include datashares, recovery points, and usage limits.

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

### List default iam role arn associated to each namespace

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