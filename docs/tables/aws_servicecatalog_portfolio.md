---
title: "Steampipe Table: aws_servicecatalog_portfolio - Query AWS Service Catalog Portfolio using SQL"
description: "Allows users to query AWS Service Catalog Portfolios. The aws_servicecatalog_portfolio table in Steampipe provides information about portfolios within AWS Service Catalog. This table allows DevOps engineers to query portfolio-specific details, including owner, description, created time, and associated metadata. Users can utilize this table to gather insights on portfolios, such as portfolio details, associated products, and more. The schema outlines the various attributes of the portfolio, including the portfolio ARN, creation date, and associated tags."
folder: "Service Catalog"
---

# Table: aws_servicecatalog_portfolio - Query AWS Service Catalog Portfolio using SQL

The AWS Service Catalog Portfolio is a part of the AWS Service Catalog that allows you to manage, organize, and govern your cloud resources. It enables you to create and manage catalogs of IT services that are approved for use on AWS. These IT Service Catalogs can help you achieve consistent governance and meet your compliance requirements, while enabling users to quickly deploy only the approved IT services they need.

## Table Usage Guide

The `aws_servicecatalog_portfolio` table in Steampipe provides you with information about portfolios within AWS Service Catalog. This table allows you, as a DevOps engineer, to query portfolio-specific details, including owner, description, created time, and associated metadata. You can utilize this table to gather insights on portfolios, such as portfolio details, associated products, and more. The schema outlines the various attributes of the portfolio for you, including the portfolio ARN, creation date, and associated tags.

## Examples

### Basic info
Explore the portfolios available in your AWS Service Catalog to understand which resources are accessible for deployment across your organization. This can help manage and govern cloud resources more effectively.

```sql+postgres
select
  display_name,
  id,
  arn,
  region,
  akas
from
  aws_servicecatalog_portfolio;
```

```sql+sqlite
select
  display_name,
  id,
  arn,
  region,
  akas
from
  aws_servicecatalog_portfolio;
```

### List portfolios of a provider
Determine the portfolios associated with a specific provider to understand the range of services they offer. This can help in comparing and selecting suitable service providers based on your needs.

```sql+postgres
select
  display_name,
  id,
  description,
  provider_name
from
  aws_servicecatalog_portfolio
where
  provider_name = 'my-portfolio';
```

```sql+sqlite
select
  display_name,
  id,
  description,
  provider_name
from
  aws_servicecatalog_portfolio
where
  provider_name = 'my-portfolio';
```

### List portfolios created in the last 30 days
Gain insights into the most recently created portfolios within the past month. This is beneficial for tracking the latest additions and updates in your AWS Service Catalog.

```sql+postgres
select
  display_name,
  id,
  description,
  created_time
from
  aws_servicecatalog_portfolio
where
  created_time >= now() - interval '30' day;
```

```sql+sqlite
select
  display_name,
  id,
  description,
  created_time
from
  aws_servicecatalog_portfolio
where
  created_time >= datetime('now', '-30 day');
```

### Get budget details of portfolios
Discover the segments that are costing you the most by understanding the budget allocation across various portfolios. This can aid in better financial management by identifying areas that require cost optimization.

```sql+postgres
select
  sp.display_name,
  sp.id,
  b ->> 'BudgetName' as budget_name
from
  aws_servicecatalog_portfolio as sp,
  jsonb_array_elements(budgets) as b;
```

```sql+sqlite
select
  sp.display_name,
  sp.id,
  json_extract(b.value, '$.BudgetName') as budget_name
from
  aws_servicecatalog_portfolio as sp,
  json_each(budgets) as b;
```