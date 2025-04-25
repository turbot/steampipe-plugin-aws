---
title: "Steampipe Table: aws_securityhub_enabled_product_subscription - Query AWS Security Hub Enabled Product Subscriptions using SQL"
description: "Allows users to query details of enabled product subscriptions in AWS Security Hub, providing insights into their activation, integration types, and company information."
folder: "Security Hub"
---

# Table: aws_securityhub_enabled_product_subscription - Query AWS Security Hub Enabled Product Subscriptions using SQL

AWS Security Hub provides a comprehensive view of your security alerts and compliance status across AWS accounts. By enabling specific security products, you can centralize and analyze security findings from various AWS services and partner solutions. The `aws_securityhub_enabled_product_subscription` table in Steampipe allows you to query information about the security products that have been enabled in AWS Security Hub.

## Table Usage Guide

The `aws_securityhub_enabled_product_subscription` table enables security analysts and cloud administrators to gather detailed insights into the products that are enabled in AWS Security Hub. You can query various aspects of these products, such as their activation URLs, integration types, categories, and company details. This table is particularly useful for monitoring the active security products, managing integrations, and ensuring that your security tools are configured correctly.

## Examples

### Basic product information
Retrieve basic information about the enabled security product subscriptions.

```sql+postgres
select
  arn,
  title,
  akas
from
  aws_securityhub_enabled_product_subscription;
```

```sql+sqlite
select
  arn,
  title,
  akas
from
  aws_securityhub_enabled_product_subscription;
```

### List products for enabled subscriptions
Identify all products for the subscriptions that are enabled.

```sql+postgres
select
  s.arn as subscription_arn,
  p.product_arn,
  p.name as product_name,
  p.company_name as product_company_name,
  p.marketplace_url,
  p.integration_types
from
  aws_securityhub_enabled_product_subscription as s,
  aws_securityhub_product as p,
  jsonb_array_elements(p.product_subscription_resource_policy -> 'Statement') as m
where
  (m ->> 'Resource') = s.arn;
```

```sql+sqlite
select
  s.arn as subscription_arn,
  p.product_arn,
  p.name as product_name,
  p.company_name as product_company_name,
  p.marketplace_url,
  p.integration_types
from
  aws_securityhub_enabled_product_subscription s,
  aws_securityhub_product p,
  json_each(p.product_subscription_resource_policy, '$.Statement') as m
where
  json_extract(m.value, '$.Resource') = s.arn;
```
