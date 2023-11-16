---
title: "Table: aws_servicecatalog_portfolio - Query AWS Service Catalog Portfolio using SQL"
description: "Allows users to query AWS Service Catalog Portfolios. The aws_servicecatalog_portfolio table in Steampipe provides information about portfolios within AWS Service Catalog. This table allows DevOps engineers to query portfolio-specific details, including owner, description, created time, and associated metadata. Users can utilize this table to gather insights on portfolios, such as portfolio details, associated products, and more. The schema outlines the various attributes of the portfolio, including the portfolio ARN, creation date, and associated tags."
---

# Table: aws_servicecatalog_portfolio - Query AWS Service Catalog Portfolio using SQL

The `aws_servicecatalog_portfolio` table in Steampipe provides information about portfolios within AWS Service Catalog. This table allows DevOps engineers to query portfolio-specific details, including owner, description, created time, and associated metadata. Users can utilize this table to gather insights on portfolios, such as portfolio details, associated products, and more. The schema outlines the various attributes of the portfolio, including the portfolio ARN, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the aws_servicecatalog_portfolio table, you can use the `.inspect aws_servicecatalog_portfolio` command in Steampipe.

**Key columns**:

- `portfolio_id`: The unique identifier for the portfolio. This column is important for joining with other tables because it is the primary key that represents each individual portfolio.
- `arn`: The Amazon Resource Name (ARN) of the portfolio. This column is useful for joining with other tables because it provides a unique identifier for the portfolio in the AWS ecosystem.
- `name`: The name of the portfolio. This column is important because it provides a human-readable identifier for the portfolio, which can be useful for reporting and analysis.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

### Get budget details of portfolios

```sql
select
  display_name,
  id,
  b ->> 'BudgetName' as budget_name
from
  aws_servicecatalog_portfolio,
  jsonb_array_elements(budgets) as b;
```