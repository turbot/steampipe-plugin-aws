---
title: "Steampipe Table: aws_securityhub_product - Query AWS Security Hub Products using SQL"
description: "Allows users to query AWS Security Hub Product details for comprehensive security and compliance insights."
folder: "Security Hub"
---

# Table: aws_securityhub_product - Query AWS Security Hub Products using SQL

The AWS Security Hub Product is a service that provides a comprehensive view of your high-priority security alerts and compliance status across AWS accounts. It aggregates, organizes, and prioritizes your security alerts, or findings, from multiple AWS services, such as Amazon GuardDuty, Amazon Inspector, and Amazon Macie, as well as from AWS Partner solutions. The findings are then visually summarized on integrated dashboards with actionable graphs and tables.

## Table Usage Guide

The `aws_securityhub_product` table in Steampipe provides you with information about security products within AWS Security Hub. This table allows you as a security analyst or DevOps engineer to query product-specific details, including product ARN, product name, company name, description, and marketplace URL. You can utilize this table to gather insights on security products, such as their activation status, associated integrations, and more. The schema outlines the various attributes of the security product for you, including the product ARN, name, company name, description, marketplace URL, and activation status.

## Examples

### Basic info
Explore which security products are in use within your AWS environment and gain insights into their associated companies and descriptions. This information can be valuable for auditing purposes, compliance checks, or for understanding the overall security posture of your AWS infrastructure.

```sql+postgres
select
  name,
  product_arn,
  company_name,
  description
from
  aws_securityhub_product;
```

```sql+sqlite
select
  name,
  product_arn,
  company_name,
  description
from
  aws_securityhub_product;
```


### List products provided by AWS
Discover the range of products provided directly by AWS, enabling you to understand the scope of services and solutions offered by the company. This can assist in identifying potential resources for your specific needs.

```sql+postgres
select
  name,
  company_name,
  description
from
  aws_securityhub_product
where
  company_name = 'AWS';
```

```sql+sqlite
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
Determine the areas in which specific products are configured to send findings to the security hub. This can be particularly useful for organizations looking to enhance their security posture by ensuring all relevant findings are centralized for further analysis.

```sql+postgres
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

```sql+sqlite
select
  name,
  product_arn,
  company_name
from
  aws_securityhub_product,
  json_each(integration_types)
where
  value = 'SEND_FINDINGS_TO_SECURITY_HUB';
```