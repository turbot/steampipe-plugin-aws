---
title: "Table: aws_securityhub_product - Query AWS Security Hub Product using SQL"
description: "Allows users to query AWS Security Hub Product details for comprehensive security and compliance insights."
---

# Table: aws_securityhub_product - Query AWS Security Hub Product using SQL

The `aws_securityhub_product` table in Steampipe provides information about security products within AWS Security Hub. This table allows security analysts and DevOps engineers to query product-specific details, including product ARN, product name, company name, description, and marketplace URL. Users can utilize this table to gather insights on security products, such as their activation status, associated integrations, and more. The schema outlines the various attributes of the security product, including the product ARN, name, company name, description, marketplace URL, and activation status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_product` table, you can use the `.inspect aws_securityhub_product` command in Steampipe.

### Key columns:

- `product_arn`: This is the AWS Resource Name (ARN) of the security product. It is a unique identifier and can be used to join this table with other tables that reference security products.
- `name`: The name of the security product. This can be useful for filtering or sorting the data based on the product name.
- `company_name`: The name of the company that provides the product. This can be useful when you need to join this table with others that contain company-specific data.

## Examples

### Basic info

```sql
select
  name,
  product_arn,
  company_name,
  description
from
  aws_securityhub_product;
```


### List products provided by AWS

```sql
select
  name,
  company_name,
  description
from
  aws_securityhub_product
where
  company_name = 'AWS';
```


### List products that send findings to security hub

```sql
select
  name,
  product_arn,
  company_name
from
  aws_securityhub_product,
  jsonb_array_elements_text(integration_types) as i
where
  i = 'SEND_FINDINGS_TO_SECURITY_HUB';
```
