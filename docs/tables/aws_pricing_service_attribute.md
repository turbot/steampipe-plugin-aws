---
title: "Steampipe Table: aws_pricing_service_attribute - Query AWS Pricing Service Attributes using SQL"
description: "Allows users to query AWS Pricing Service Attributes to gain insights into product attributes and their respective prices."
folder: "Pricing"
---

# Table: aws_pricing_service_attribute - Query AWS Pricing Service Attributes using SQL

The AWS Pricing Service Attributes are part of the AWS Price List Service API, which provides programmatic access to information about product prices. It allows you to query for the prices of AWS services, using attributes such as service code, location, or usage type. This service can help you manage your AWS costs by providing detailed and up-to-date pricing information.

## Table Usage Guide

The `aws_pricing_service_attribute` table in Steampipe provides you with information about product attributes within AWS Pricing. This table empowers you, as a DevOps engineer, financial analyst, or cloud architect, to query product-specific details, including service codes, product families, and associated metadata. You can utilize this table to gather insights on AWS products, such as their pricing details, product families, usage types, and more. The schema outlines the various attributes of the AWS product for you, including the service code, product family, instance type, location, and usage type.

## Examples

### Basic info
Discover the segments that determine pricing in your AWS services by analyzing the attributes and their respective values. This can help you understand cost drivers and potentially identify areas for cost optimization.

```sql+postgres
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute;
```

```sql+sqlite
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute;
```

### List attribute details of AWS Backup service
Analyze the settings to understand the specific attributes and associated values of the AWS Backup service. This can be beneficial in assessing the cost implications of different attributes, helping to optimize budget allocation for AWS services.

```sql+postgres
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup';
```

```sql+sqlite
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
Discover the segments that are supported by the AWS Backup service for the 'termType' attribute. This is useful in understanding the different term types available, aiding in effective cost and resource management.

```sql+postgres
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup' and attribute_name = 'termType';
```

```sql+sqlite
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup' and attribute_name = 'termType';
```