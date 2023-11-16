---
title: "Table: aws_pricing_service_attribute - Query AWS Pricing Service Attributes using SQL"
description: "Allows users to query AWS Pricing Service Attributes to gain insights into product attributes and their respective prices."
---

# Table: aws_pricing_service_attribute - Query AWS Pricing Service Attributes using SQL

The `aws_pricing_service_attribute` table in Steampipe provides information about product attributes within AWS Pricing. This table allows DevOps engineers, financial analysts, and cloud architects to query product-specific details, including service codes, product families, and associated metadata. Users can utilize this table to gather insights on AWS products, such as their pricing details, product families, usage types, and more. The schema outlines the various attributes of the AWS product, including the service code, product family, instance type, location, and usage type.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_pricing_service_attribute` table, you can use the `.inspect aws_pricing_service_attribute` command in Steampipe.

**Key columns**:

- `service_code`: This column provides the service code for the AWS product. It is useful for identifying the specific AWS service that the product belongs to.
- `product_family`: This column provides the product family for the AWS product. It is useful for categorizing AWS products and can be used to join this table with others that also contain product family information.
- `location`: This column provides the location where the AWS product is available. It is useful for understanding the geographical distribution of AWS products and can be used to join this table with others that also contain location information.

## Examples

### Basic info

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute;
```

### List attribute details of AWS Backup service

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup';
```

### List supported attribute values of AWS Backup service for termType attribute

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup' and attribute_name = 'termType';
```
