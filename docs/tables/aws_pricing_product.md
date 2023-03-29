# Table: aws_pricing_product

Returns costs, and offer definition, for the corresponding service.

## Examples

### List pricing offers for on-demand shared EC2 c5.2xlarge without pre-installed software, with Linux OS

```sql
select
  term,
  purchase_option,
  lease_contract_length,
  unit,
  price_per_unit::numeric::money,
  currency,
  begin_range,
  end_range,
  effective_date,
  description,
  attributes ->> 'instanceType',
  attributes ->> 'vcpu',
  attributes ->> 'memory',
  attributes ->> 'operatingSystem',
  attributes ->> 'preInstalledSw'
from
  aws_pricing_product
where
  service_code = 'AmazonEC2'
  and filters = '{
  "regionCode": "eu-west-3",
  "locationType": "AWS Region",
  "instanceType": "c5.2xlarge",
  "operatingSystem": "Linux",
  "tenancy": "Shared",
  "preInstalledSw": "NA",
  "capacityStatus": "Used" }'::jsonb;
```


### List pricing offers for Mysql RDS db.m5.xlarge instance in eu-west-3 in single-az deployment

```sql
select
  term,
  purchase_option,
  lease_contract_length,
  unit,
  price_per_unit::numeric::money,
  currency,
  attributes ->> 'instanceType',
  attributes ->> 'vcpu',
  attributes ->> 'memory',
  attributes ->> 'databaseEngine',
  attributes ->> 'deploymentOption'
from
  aws_pricing_product
where
  service_code = 'AmazonRDS'
  and filters = '{
  "regionCode": "eu-west-3",
  "locationType": "AWS Region",
  "instanceType": "db.m5.xlarge",
  "databaseEngine": "MySQL",
  "deploymentOption": "Single-AZ" }'::jsonb;
```

### List pricing offers for Redis ElasticCache cache.m5.xlarge instances in eu-west-3

```sql
select
  term,
  purchase_option,
  lease_contract_length,
  unit,
  price_per_unit::numeric::money,
  currency,
  attributes ->> 'instanceType',
  attributes ->> 'vcpu',
  attributes ->> 'memory',
  attributes ->> 'cacheEngine'
from
  aws_pricing_product
where
  service_code = 'AmazonElastiCache'
  and filters = '{
  "regionCode": "eu-west-3",
  "locationType": "AWS Region",
  "instanceType": "cache.m5.xlarge",
  "cacheEngine": "Redis" }'::jsonb;
```
