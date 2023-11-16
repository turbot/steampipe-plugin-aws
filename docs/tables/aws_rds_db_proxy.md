---
title: "Table: aws_rds_db_proxy - Query Amazon RDS DB Proxy using SQL"
description: "Allows users to query DB Proxies in Amazon RDS to fetch detailed information about each proxy, including its ARN, name, engine family, role ARN, status, and more."
---

# Table: aws_rds_db_proxy - Query Amazon RDS DB Proxy using SQL

The `aws_rds_db_proxy` table in Steampipe provides information about DB Proxies within Amazon Relational Database Service (RDS). This table allows DevOps engineers, DBAs, and other technical professionals to query proxy-specific details, including the proxy's ARN, name, engine family, role ARN, status, and more. Users can utilize this table to gather insights on DB Proxies, such as their current status, the engine family they are associated with, the IAM role they are utilizing, and more. The schema outlines the various attributes of the DB Proxy, including the proxy ARN, creation date, updated date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_proxy` table, you can use the `.inspect aws_rds_db_proxy` command in Steampipe.

### Key columns:

- `proxy_arn`: The Amazon Resource Name (ARN) for the DB Proxy. This can be used to join this table with other tables that also contain proxy ARNs.
- `name`: The name of the DB Proxy. This can be useful for joining with other tables where the proxy name is known.
- `role_arn`: The Amazon Resource Name (ARN) for the IAM role that is associated with the DB Proxy. This can be used to join with IAM role tables to fetch more details about the role.

## Examples

### List of DB proxies and corresponding engine families

```sql
select
  db_proxy_name,
  status,
  engine_family
from
  aws_rds_db_proxy;
```

### Authentication info of each DB proxy

```sql
select
  db_proxy_name,
  engine_family,
  a ->> 'AuthScheme' as auth_scheme,
  a ->> 'Description' as auth_description,
  a ->> 'IAMAuth' as iam_auth,
  a ->> 'SecretArn' as secret_arn,
  a ->> 'UserName' as user_name
from
  aws_rds_db_proxy,
  jsonb_array_elements(auth) as a;
```
