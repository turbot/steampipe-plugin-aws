# Table: aws_pricing_product_json

AWS Pricing Service provides pricing information for all products in AWS.

**Important notes:**

- You ***must*** specify `service_code` and `attributes` in a where or join clause in order to use this table.  
- The `attributes` column represent the product metadata field and value that you want to filter on. You can filter by just the service code to see all products for a specific service, filter by just the attribute name to see a specific attribute for multiple services, or use both a service code and an attribute name to retrieve only products that match both fields.
- You can get attribute details for each service in `aws_service_attribute` table.

## Examples

### List on-demand pricing details of EC2 instance type t3.micro in us-east-1

```sql
select
  service_code,
  attributes,
  jsonb_pretty(terms -> 'OnDemand') as on_demand_pricing_detail,
  product -> 'attributes' ->> 'location' as location
from
  aws_pricing_product_json
where
  service_code = 'AmazonEC2' 
  and attributes = jsonb_object('{instanceType, regionCode}','{t3.micro, us-east-1}');
```

### List reserved pricing details of EC2 instance type c5.xlarge with operating system linux and no license required

```sql
select
  service_code,
  attributes,
  jsonb_pretty(terms -> 'Reserved') as reserved_pricing_detail,
  product -> 'attributes' ->> 'location' as location
from
  aws_pricing_product_json
where
  service_code = 'AmazonEC2' 
  and attributes = jsonb_object('{instanceType, operatingSystem, licenseModel}','{c5.xlarge, Linux, No License required}');
```

### List pricing details of amazon redshift instance type ds1.xlarge in us-west-2

```sql
select
  service_code,
  attributes,
  jsonb_pretty(terms) as pricing_detail,
  product -> 'attributes' -> 'location' as location
from
  aws_pricing_product_json
where
  service_code = 'AmazonRedshift' 
  and attributes = jsonb_object('{instanceType, tenancy, regionCode}','{ds1.xlarge, Shared, us-east-1}');
```