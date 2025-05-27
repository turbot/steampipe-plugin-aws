---
title: "Steampipe Table: aws_pricing_product - Query AWS Pricing Product using SQL"
description: "Allows users to query AWS Pricing Product details such as the product's description, pricing details, and associated attributes."
folder: "Pricing"
---

# Table: aws_pricing_product - Query AWS Pricing Product using SQL

The AWS Pricing Product is a service that provides pricing information for AWS services. It allows you to retrieve current or historical prices and understand your costs for AWS services better. This service is essential for cost management and optimization in an AWS environment.

## Table Usage Guide

The `aws_pricing_product` table in Steampipe provides you with information about pricing products within AWS Pricing. This table allows you, whether you're a financial analyst, cloud cost manager, or DevOps engineer, to query product-specific details, including product descriptions, pricing details, and associated attributes. You can utilize this table to gather insights on products, such as the cost of each AWS service, the pricing model, and the location. The schema outlines the various attributes of the pricing product for you, including the product description, pricing details, and associated attributes.

## Examples

### List pricing offers for on-demand shared EC2 c5.2xlarge without pre-installed software, with Linux OS
Determine the pricing options for on-demand shared EC2 c5.2xlarge instances without pre-installed software, running on Linux OS. This is useful for cost planning and budgeting for your AWS resources.

```sql+postgres
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

```sql+sqlite
select
  term,
  purchase_option,
  lease_contract_length,
  unit,
  CAST(price_per_unit AS REAL) as price_per_unit,
  currency,
  begin_range,
  end_range,
  effective_date,
  description,
  json_extract(attributes, '$.instanceType'),
  json_extract(attributes, '$.vcpu'),
  json_extract(attributes, '$.memory'),
  json_extract(attributes, '$.operatingSystem'),
  json_extract(attributes, '$.preInstalledSw')
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
  "capacityStatus": "Used" }';
```


### List pricing offers for Mysql RDS db.m5.xlarge instance in eu-west-3 in single-az deployment
Explore pricing offers for a specific type of MySQL RDS instance in a certain region and deployment. This is useful for cost analysis and budget planning for AWS services.

```sql+postgres
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

```sql+sqlite
select
  term,
  purchase_option,
  lease_contract_length,
  unit,
  price_per_unit,
  currency,
  json_extract(attributes, '$.instanceType'),
  json_extract(attributes, '$.vcpu'),
  json_extract(attributes, '$.memory'),
  json_extract(attributes, '$.databaseEngine'),
  json_extract(attributes, '$.deploymentOption')
from
  aws_pricing_product
where
  service_code = 'AmazonRDS'
  and filters = '{
  "regionCode": "eu-west-3",
  "locationType": "AWS Region",
  "instanceType": "db.m5.xlarge",
  "databaseEngine": "MySQL",
  "deploymentOption": "Single-AZ" }';
```

### List pricing offers for Redis ElasticCache cache.m5.xlarge instances in eu-west-3
Explore the different pricing options available for Redis ElasticCache cache.m5.xlarge instances in the eu-west-3 region. This can help to make informed decisions about cost management and budgeting.

```sql+postgres
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

```sql+sqlite
select
  term,
  purchase_option,
  lease_contract_length,
  unit,
  price_per_unit,
  currency,
  json_extract(attributes, '$.instanceType'),
  json_extract(attributes, '$.vcpu'),
  json_extract(attributes, '$.memory'),
  json_extract(attributes, '$.cacheEngine')
from
  aws_pricing_product
where
  service_code = 'AmazonElastiCache'
  and filters = '{
  "regionCode": "eu-west-3",
  "locationType": "AWS Region",
  "instanceType": "cache.m5.xlarge",
  "cacheEngine": "Redis" }';
```