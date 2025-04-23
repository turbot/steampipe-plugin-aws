---
title: "Steampipe Table: aws_lakeformation_tag - Query AWS Lake Formation Tags Using SQL"
description: "Query AWS Lake Formation LF-tags, including their tag keys, possible tag values, and associated AWS accounts."
folder: "Lake Formation"
---

# Table: aws_lakeformation_tag - Query AWS Lake Formation Tags Using SQL

The `aws_lakeformation_tag` table allows you to query **AWS Lake Formation LF-tags**, providing details about the **tag keys, possible tag values, and associated AWS accounts**. This table helps **data governance teams and security administrators** monitor and manage **LF-tag-based access control** effectively.

## Table Usage Guide

The `aws_lakeformation_tag` table provides insights into **LF-tags** applied to AWS Lake Formation resources. LF-tags (Lake Formation tags) enable attribute-based access control (ABAC), allowing administrators to grant permissions dynamically based on tag keys and values instead of manually assigning policies to users or roles. This table helps track **registered LF-tags**, their associated **AWS account, region, and partition**, and provides a list of **possible values** an attribute can take.

## **Examples**

### List all AWS Lake Formation LF-tags
Retrieve a list of all LF-tags registered in AWS Lake Formation, including their key names and possible values.

```sql+postgres
select
  catalog_id
  tag_key,
  tag_values
from
  aws_lakeformation_tag;
```

```sql+sqlite
select
  catalog_id
  tag_key,
  tag_values
from
  aws_lakeformation_tag;
```

### Find LF-tags in a specific AWS Region
Identify LF-tags that are registered in a particular AWS region.

```sql+postgres
select
  tag_key,
  tag_values,
  region
from
  aws_lakeformation_tag
where
  region = 'us-east-1';
```

```sql+sqlite
select
  tag_key,
  tag_values,
  region
from
  aws_lakeformation_tag
where
  region = 'us-east-1';
```

### Get LF-tags associated with for a specific catalog
Find all LF-tags that belong to a given catalog.

```sql+postgres
select
  tag_key,
  tag_values,
  account_id
from
  aws_lakeformation_tag
where
  account_id = '123456789012';
```

```sql+sqlite
select
  tag_key,
  tag_values,
  account_id
from
  aws_lakeformation_tag
where
  catalog_id = '123456789012';
```

### List LF-tags with multiple possible values
Retrieve LF-tags that have multiple values assigned.

```sql+postgres
select
  tag_key,
  tag_values
from
  aws_lakeformation_tag
where
  jsonb_array_length(tag_values) > 1;
```

```sql+sqlite
select
  tag_key,
  tag_values
from
  aws_lakeformation_tag
where
  jsonb_array_length(tag_values) > 1;
```