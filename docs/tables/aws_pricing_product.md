# Table: aws_pricing_product

AWS Pricing Service provides pricing information for all products in AWS.

**Important notes:**

- You ***must*** specify `service_code`, `field` and `value` in a where or join clause in order to use this table.  
- The `field` column represent the product metadata field that you want to filter on. You can filter by just the service code to see all products for a specific service, filter by just the attribute name to see a specific attribute for multiple services, or use both a service code and an attribute name to retrieve only products that match both fields.
- The `value` column represent the service code or attribute value that you want to filter by.

## Examples

### List pricing details of EC2 instance type t3.micro

```sql
select
  service_code,
  field,
  value,
  jsonb_pretty(terms -> 'OnDemand') as on_demand_pricing_detail,
  product -> 'attributes' -> 'location' as location
from
  aws_pricing_product
where
  service_code = 'AmazonEC2' and field = 'instanceType' and value = 't3.micro';
```

### List pricing details of SNS topic in us-west-2 region

```sql
select
  service_code,
  field,
  value,
  jsonb_pretty(terms -> 'OnDemand') as on_demand_pricing_detail
from
  aws_pricing_product
where
  service_code = 'AmazonSNS' and field = 'regionCode' and value = 'us-west-2';
```

### List pricing details of amazon redshift instance type ds1.xlarge

```sql
select
  service_code,
  field,
  value,
  jsonb_pretty(terms) as pricing_detail,
  product -> 'attributes' -> 'location' as location
from
  aws_pricing_product
where
  service_code = 'AmazonRedshift' and field = 'instanceType' and value = 'ds1.xlarge';
```