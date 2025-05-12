---
title: "Steampipe Table: aws_servicecatalog_provisioned_product - Query AWS Service Catalog Provisioned Product using SQL"
description: "Allows users to query AWS Service Catalog Provisioned Product data including product details, owner, type, and associated metadata."
folder: "Service Catalog"
---

# Table: aws_servicecatalog_provisioned_product - Query AWS Service Catalog Provisioned Product using SQL

A provisioned product is a resourced instance of a product. For example, provisioning a product based on a CloudFormation template launches a CloudFormation stack and its underlying resources.

## Table Usage Guide

The `aws_servicecatalog_provisioned_product` table in Steampipe provides you with information about provisioned products within AWS Service Catalog. This table allows you, as a DevOps engineer, to query product-specific details, including product owner, type, and associated metadata. You can utilize this table to gather insights on products such as their distribution status, launch paths, provisioning artifacts, and more. The schema outlines the various attributes of the AWS Service Catalog Provisioned Product for you, including the product ARN, creation time, product type, owner.

**Important notes:**
This table supports optional quals. Queries with optional quals are optimised to use search filters. Optional quals are supported for the following columns:
	- `created_time`
	- `id`
	- `last_record_id`
	- `idempotency_token`
	- `name`
	- `product_id`
	- `type`
	- `status`
	- `last_provisioning_record_id`
	- `last_successful_provisioning_record_id`

## Examples

### Basic info
This query can be very useful for getting a comprehensive overview of provisioned products in your AWS environment, particularly for inventory management, auditing, and tracking the status and details of various service catalog products. It helps in understanding what services are currently deployed, their status, and key identifiers that might be required for further management or automation tasks.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product;
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product;
```

### List the provisioned products created in the last 7 days
The query you've provided is useful for retrieving information about AWS Service Catalog provisioned products that were created within the last 7 days. It's particularly valuable for monitoring and managing recent resource provisioning activities in AWS.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product
where
  created_time >= (current_date - interval '7' day)
order by
  created_time;
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product
where
  created_time >= date('now', '-7 day')
order by
  created_time;
```

### Get product details of the successfully provisioned product
Filters provisioned products to include only those products where a provisioning process has been successfully completed at least once. This can be particularly useful for maintaining an inventory of active and successfully set up resources, aiding in tracking and managing AWS resources effectively.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product
where
  last_successful_provisioning_record_id is not null;
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product
where
  last_successful_provisioning_record_id is not null;
```

### List details of the successfully provisioned product
This allows you to see detailed information about each provisioned product, such as its name, type, support URL, and support email, alongside its provisioning ID. It's particularly useful for having a consolidated view of product details along with their provisioning status, especially for those products with a successful provisioning record.

```sql+postgres
select
  pr.id as provisioning_id,
  p.name as product_name,
  p.id as product_view_id,
  p.product_id,
  p.type as product_type,
  p.support_url as product_support_url,
  p.support_email as product_support_email
from
  aws_servicecatalog_provisioned_product as pr,
  aws_servicecatalog_product as p
where
  pr.product_id = p.product_id
  and last_successful_provisioning_record_id is not null;
```

```sql+sqlite
select
  pr.id as provisioning_id,
  p.name as product_name,
  p.id as product_view_id,
  p.product_id,
  p.type as product_type,
  p.support_url as product_support_url,
  p.support_email as product_support_email
from
  aws_servicecatalog_provisioned_product as pr
join
  aws_servicecatalog_product as p on pr.product_id = p.product_id
where
  pr.last_successful_provisioning_record_id is not null;
```

### List the provisioned products of CFN_STACK type
Ensures that the provisioned products that have been successfully deployed at least once. This is beneficial for monitoring and auditing successful deployments, understanding resource utilization, and managing CloudFormation-based resources effectively.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product
where
  type = 'CFN_STACK'
  and last_successful_provisioning_record_id is not null;
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  product_id,
  status,
  created_time,
  last_provisioning_record_id
from
  aws_servicecatalog_provisioned_product
where
  type = 'CFN_STACK'
  and last_successful_provisioning_record_id is not null;
```