---
title: "Steampipe Table: aws_servicecatalog_product - Query AWS Service Catalog Product using SQL"
description: "Allows users to query AWS Service Catalog Product data including product details, owner, type, and associated metadata."
folder: "Service Catalog"
---

# Table: aws_servicecatalog_product - Query AWS Service Catalog Product using SQL

The AWS Service Catalog Product allows you to create, manage, and distribute catalogs of approved IT services. These services can be anything from virtual machine images, servers, software, and databases to complete multi-tier application architectures. It helps organizations to manage their IT services and resources with a standardized approach and eliminates the need for manual processes.

## Table Usage Guide

The `aws_servicecatalog_product` table in Steampipe provides you with information about products within AWS Service Catalog. This table allows you, as a DevOps engineer, to query product-specific details, including product owner, type, and associated metadata. You can utilize this table to gather insights on products such as their distribution status, launch paths, provisioning artifacts, and more. The schema outlines the various attributes of the AWS Service Catalog Product for you, including the product ARN, creation time, product type, owner, and associated tags.

## Examples

### Basic info
Explore which AWS Service Catalog products are available, along with their associated identifiers and support details. This can be useful in managing and supporting your organization's AWS resources effectively.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have a preset path in the AWS Service Catalog. This can be beneficial to assess the elements within your AWS environment that come with predefined settings.

```sql+postgres
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

```sql+sqlite
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
  has_default_path = 1;
```

### List products that are owned by AWS
Discover the segments that are owned by AWS in the marketplace. This could be useful in gaining insights into the variety of products and their support details provided by AWS in the marketplace.

```sql+postgres
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

```sql+sqlite
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
Determine the budget allocation for each product in your AWS Service Catalog. This helps in financial planning and resource management by identifying where your money is being spent.

```sql+postgres
select
  sp.name,
  sp.id,
  sp.owner,
  sp.product_id,
  sp.short_description,
  b ->> 'BudgetName' as budget_name
from
  aws_servicecatalog_product as sp,
  jsonb_array_elements(budgets) as b;
```

```sql+sqlite
select
  sp.name,
  sp.id,
  sp.owner,
  sp.product_id,
  sp.short_description,
  json_extract(b.value, '$.BudgetName') as budget_name
from
  aws_servicecatalog_product as sp,
  json_each(budgets) as b;
```

### Get launch path details of each product
Identify the launch path details for each product in your AWS Service Catalog. This can help you understand the unique identifiers and names associated with each product's launch path, aiding in efficient product management and organization.

```sql+postgres
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

```sql+sqlite
select
  sp.name,
  sp.id,
  sp.owner,
  short_description,
  json_extract(l.value, '$.Id') as launch_path_id,
  json_extract(l.value, '$.Name') as launch_path_name
from
  aws_servicecatalog_product as sp,
  json_each(aws_servicecatalog_product.launch_paths) as l;
```

### Get provisioning artifact details for each product
This query enables users to gain insights into the details of each provisioning artifact associated with their products on AWS Service Catalog. It is useful for tracking and managing product versions and configurations, which can aid in maintaining standardized environments across an organization.

```sql+postgres
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

```sql+sqlite
select
  sp.name,
  sp.id,
  json_extract(p.value, '$.Id') as provisioning_artifact_id,
  json_extract(p.value, '$.Name') as provisioning_artifact_name,
  json_extract(p.value, '$.CreatedTime') as provisioning_artifact_created_time,
  json_extract(p.value, '$.Description') as provisioning_artifact_description,
  json_extract(p.value, '$.Guidance') as provisioning_artifact_guidance
from
  aws_servicecatalog_product as sp,
  json_each(provisioning_artifacts) as p;
```