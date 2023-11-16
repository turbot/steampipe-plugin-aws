---
title: "Table: aws_servicecatalog_product - Query AWS Service Catalog Product using SQL"
description: "Allows users to query AWS Service Catalog Product data including product details, owner, type, and associated metadata."
---

# Table: aws_servicecatalog_product - Query AWS Service Catalog Product using SQL

The `aws_servicecatalog_product` table in Steampipe provides information about products within AWS Service Catalog. This table allows DevOps engineers to query product-specific details, including product owner, type, and associated metadata. Users can utilize this table to gather insights on products such as their distribution status, launch paths, provisioning artifacts, and more. The schema outlines the various attributes of the AWS Service Catalog Product, including the product ARN, creation time, product type, owner, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_servicecatalog_product` table, you can use the `.inspect aws_servicecatalog_product` command in Steampipe.

**Key columns**:

- `name`: The name of the product. This can be used to join with other tables that contain product information.
- `arn`: The Amazon Resource Name (ARN) of the product. This unique identifier is useful for joining with other tables that reference AWS resources.
- `product_id`: The identifier of the product. This can be used to join with other tables that contain product-specific information.

## Examples

### Basic info

```sql
select
  name,
  id,
  product_id,
  type,
  akas,
  support_url,
  support_email
from
  aws_servicecatalog_product;
```

### List products that have a default path

```sql
select
  name,
  id,
  product_id,
  type,
  distributor,
  owner,
  has_default_path
from
  aws_servicecatalog_product
where
  has_default_path;
```

### List products that are owned by AWS

```sql
select
  name,
  id,
  product_id,
  type,
  support_url,
  support_description
from
  aws_servicecatalog_product
where
  type = 'MARKETPLACE';
```

### Get budget details of each product

```sql
select
  name,
  id,
  owner,
  product_id,
  short_description,
  b ->> 'BudgetName' as budget_name
from
  aws_servicecatalog_product,
  jsonb_array_elements(budgets) as b;
```

### Get launch path details of each product

```sql
select
  name,
  id,
  owner,
  short_description,
  l ->> 'Id' as launch_path_id,
  l ->> 'Name' as launch_path_name
from
  aws_servicecatalog_product,
  jsonb_array_elements(launch_paths) as l;
```

### Get provisioning artifact details for each product

```sql
select
  name,
  id,
  p ->> 'Id' as provisioning_artifact_id,
  p ->> 'Name' as provisioning_artifact_name,
  p ->> 'CreatedTime' as provisioning_artifact_created_time,
  p ->> 'Description' as provisioning_artifact_description,
  p ->> 'Guidance' as provisioning_artifact_guidance
from
  aws_servicecatalog_product,
  jsonb_array_elements(provisioning_artifacts) as p;
```