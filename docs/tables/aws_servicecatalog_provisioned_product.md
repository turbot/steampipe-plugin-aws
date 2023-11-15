# Table: aws_servicecatalog_provisioned_product

A provisioned product is a resourced instance of a product. For example, provisioning a product based on a CloudFormation template launches a CloudFormation stack and its underlying resources.
This table supports optional quals. Queries with optional quals are optimised to use search filters. Optional quals are supported for the following columns:
  - `arn`
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

```sql
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

```sql
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

### Get product details of the successfully provisioned product

```sql
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

```sql
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

### List the provisioned products of CFN_STACK type

```sql
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

### List the products that went to 'ERROR' state during provisioning

```sql
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
  search_query = 'status:ERROR';
```