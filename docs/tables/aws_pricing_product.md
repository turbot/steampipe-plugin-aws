---
title: "Table: aws_pricing_product - Query AWS Pricing Product using SQL"
description: "Allows users to query AWS Pricing Product details such as the product's description, pricing details, and associated attributes."
---

# Table: aws_pricing_product - Query AWS Pricing Product using SQL

The `aws_pricing_product` table in Steampipe provides information about pricing products within AWS Pricing. This table allows financial analysts, cloud cost managers, and DevOps engineers to query product-specific details, including product descriptions, pricing details, and associated attributes. Users can utilize this table to gather insights on products, such as the cost of each AWS service, the pricing model, and the location. The schema outlines the various attributes of the pricing product, including the product description, pricing details, and associated attributes.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_pricing_product` table, you can use the `.inspect aws_pricing_product` command in Steampipe.

### Key columns:

- `product_family` - This column is important as it allows users to filter or join with other tables based on the product family. For example, you can find all pricing details for a specific product family.
- `servicecode` - This column is useful for joining with other tables that contain service code information. This allows for more detailed queries related to specific AWS services.
- `location` - This column is important as it allows users to filter or join with other tables based on the location. This can be useful for queries related to the cost of services in specific geographical locations.

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
