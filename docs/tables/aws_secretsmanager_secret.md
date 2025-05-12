---
title: "Steampipe Table: aws_secretsmanager_secret - Query AWS Secrets Manager Secret using SQL"
description: "Allows users to query AWS Secrets Manager Secret data, including metadata, versions, rotation configuration, and more."
folder: "ECR"
---

# Table: aws_secretsmanager_secret - Query AWS Secrets Manager Secret using SQL

The AWS Secrets Manager Secret is a secure and scalable service that enables you to easily manage secrets throughout their lifecycle. Using Secrets Manager, you can secure, audit, and manage secrets used to access resources in the AWS Cloud, on third-party services, and on-premises. This service helps protect access to your applications, services, and IT resources without the upfront investment and on-going maintenance costs of operating your own infrastructure.

## Table Usage Guide

The `aws_secretsmanager_secret` table in Steampipe provides you with information about secrets within AWS Secrets Manager. This table allows you, as a DevOps engineer, to query secret-specific details, including metadata, versions, rotation configuration, and more. You can utilize this table to gather insights on secrets, such as secret rotation status, associated resource policies, and more. The schema outlines the various attributes of the secret for you, including the secret ARN, name, description, rotation rules, and associated tags.

## Examples

### Basic info
Gain insights into the creation and last accessed dates of your AWS Secrets Manager secrets. This can help in managing secret lifecycle, ensuring secrets are regularly updated or identifying unused secrets.

```sql+postgres
select
  name,
  created_date,
  description,
  last_accessed_date
from
  aws_secretsmanager_secret;
```

```sql+sqlite
select
  name,
  created_date,
  description,
  last_accessed_date
from
  aws_secretsmanager_secret;
```


### List secrets that do not automatically rotate
Discover the segments that contain secrets which do not have an automatic rotation feature enabled. This is useful for identifying potential security risks and ensuring best practices for data safety.

```sql+postgres
select
  name,
  created_date,
  description,
  rotation_enabled
from
  aws_secretsmanager_secret
where
  not rotation_enabled;
```

```sql+sqlite
select
  name,
  created_date,
  description,
  rotation_enabled
from
  aws_secretsmanager_secret
where
  rotation_enabled = 0;
```


### List secrets that automatically rotate every 7 days
Identify the secrets in your AWS Secrets Manager that are set to automatically rotate more frequently than every 7 days. This can be useful for maintaining a high level of security by ensuring that secrets are updated regularly.

```sql+postgres
select
  name,
  created_date,
  description,
  rotation_enabled,
  rotation_rules
from
  aws_secretsmanager_secret
where
  rotation_rules -> 'AutomaticallyAfterDays' > '7';
```

```sql+sqlite
select
  name,
  created_date,
  description,
  rotation_enabled,
  rotation_rules
from
  aws_secretsmanager_secret
where
  json_extract(rotation_rules, '$.AutomaticallyAfterDays') > 7;
```


### List secrets that are not replicated in other regions
Determine the areas in which certain secrets are not replicated across different regions. This can be useful for ensuring data redundancy and mitigating risks associated with data loss in specific geographical locations.

```sql+postgres
select
  name,
  created_date,
  description,
  replication_status
from
  aws_secretsmanager_secret
where
  replication_status is null;
```

```sql+sqlite
select
  name,
  created_date,
  description,
  replication_status
from
  aws_secretsmanager_secret
where
  replication_status is null;
```

### List policy details for the secrets
Determine the specifics of policies pertaining to your secrets. This query is useful for gaining insights into your secret management policies, helping you understand and manage your security better.

```sql+postgres
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_secretsmanager_secret;
```

```sql+sqlite
select
  name,
  policy,
  policy_std
from
  aws_secretsmanager_secret;
```