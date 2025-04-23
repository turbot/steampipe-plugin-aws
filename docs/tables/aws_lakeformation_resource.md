---
title: "Steampipe Table: aws_lakeformation_resource - Query AWS Lake Formation Resources Using SQL"
description: "Query AWS Lake Formation registered resources, including their ARNs, associated IAM roles, and data access configurations."
folder: "Lake Formation"
---

# Table: aws_lakeformation_resource - Query AWS Lake Formation Resources Using SQL

The `aws_lakeformation_resource` table allows you to query **AWS Lake Formation registered resources**, including details about the **Amazon S3 locations** registered with Lake Formation, the IAM role used for registration, and whether hybrid access is enabled. This table helps **data governance teams and security administrators** monitor and manage **data lake access control** effectively.

## Table Usage Guide

The `aws_lakeformation_resource` table provides insights into registered **Lake Formation resources**, enabling users to identify **registered S3 locations** managed by Lake Formation and determine which **IAM role** was used for resource registration. It also allows users to check if **hybrid access** is enabled, which permits both **Lake Formation permissions and S3 bucket policies** to manage access. Additionally, this table helps track **when a resource was last modified** and filter resources based on attributes such as **AWS account, region, and partition**, making it a valuable tool for data governance and access control.

## **Examples**

### List all registered AWS Lake Formation resources
Retrieve a list of all resources registered in AWS Lake Formation, along with their associated IAM roles and modification timestamps.

```sql+postgres
select
  resource_arn,
  role_arn,
  last_modified
from
  aws_lakeformation_resource;
```

```sql+sqlite
select
  resource_arn,
  role_arn,
  last_modified
from
  aws_lakeformation_resource;
```

### Find resources with hybrid access enabled
Identify resources where **both Lake Formation and S3 bucket policies** manage access.

```sql+postgres
select
  resource_arn,
  role_arn,
  hybrid_access_enabled
from
  aws_lakeformation_resource
where
  hybrid_access_enabled = true;
```

```sql+sqlite
select
  resource_arn,
  role_arn,
  hybrid_access_enabled
from
  aws_lakeformation_resource
where
  hybrid_access_enabled = true;
```

### Get resources registered with a specific IAM role
Find all resources registered by a specific **IAM role** in AWS Lake Formation.

```sql+postgres
select
  resource_arn,
  role_arn
from
  aws_lakeformation_resource
where
  role_arn = 'arn:aws:iam::123456789012:role/MyLakeFormationRole';
```

```sql+sqlite
select
  resource_arn,
  role_arn
from
  aws_lakeformation_resource
where
  role_arn = 'arn:aws:iam::123456789012:role/MyLakeFormationRole';
```

### Check for federated Lake Formation resources
List all resources that are **federated**, meaning they are accessible across AWS accounts.

```sql+postgres
select
  resource_arn,
  role_arn,
  with_federation
from
  aws_lakeformation_resource
where
  with_federation = true;
```

```sql+sqlite
select
  resource_arn,
  role_arn,
  with_federation
from
  aws_lakeformation_resource
where
  with_federation = true;
```