---
title: "Steampipe Table: aws_rds_db_proxy - Query Amazon RDS DB Proxy using SQL"
description: "Allows users to query DB Proxies in Amazon RDS to fetch detailed information about each proxy, including its ARN, name, engine family, role ARN, status, and more."
folder: "RDS"
---

# Table: aws_rds_db_proxy - Query Amazon RDS DB Proxy using SQL

The Amazon RDS DB Proxy is a fully managed, highly available database proxy for Amazon Relational Database Service (RDS) that makes applications more scalable, more resilient to database failures, and more secure. The DB Proxy maintains a pool of established connections to your RDS instances, improving efficiency by reducing the time and resources associated with establishing new connections. It also provides failover support for the databases and includes features for enhanced security.

## Table Usage Guide

The `aws_rds_db_proxy` table in Steampipe provides you with information about DB Proxies within Amazon Relational Database Service (RDS). This table enables you, as a DevOps engineer, DBA, or other technical professional, to query proxy-specific details, including the proxy's ARN, name, engine family, role ARN, status, and more. You can utilize this table to gather insights on DB Proxies, such as their current status, the engine family they are associated with, the IAM role they are utilizing, and more. The schema outlines the various attributes of the DB Proxy for you, including the proxy ARN, creation date, updated date, and associated tags.

## Examples

### List of DB proxies and corresponding engine families
Explore which database proxies are currently active and their corresponding engine families. This can help in managing and optimizing your database resources effectively.

```sql+postgres
select
  db_proxy_name,
  status,
  engine_family
from
  aws_rds_db_proxy;
```

```sql+sqlite
select
  db_proxy_name,
  status,
  engine_family
from
  aws_rds_db_proxy;
```

### Authentication info of each DB proxy
Explore the authentication information for each database proxy. This can be useful to understand the authentication methods in use, which can aid in security audits and compliance checks.

```sql+postgres
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

```sql+sqlite
select
  db_proxy_name,
  engine_family,
  json_extract(a.value, '$.AuthScheme') as auth_scheme,
  json_extract(a.value, '$.Description') as auth_description,
  json_extract(a.value, '$.IAMAuth') as iam_auth,
  json_extract(a.value, '$.SecretArn') as secret_arn,
  json_extract(a.value, '$.UserName') as user_name
from
  aws_rds_db_proxy,
  json_each(auth) as a;
```