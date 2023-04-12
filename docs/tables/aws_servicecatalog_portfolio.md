# Table: aws_servicecatalog_portfolio

AWS Service Catalog allows IT administrators to create, manage, and distribute catalogs of approved products to end users, who can then access the products they need in a personalized portal.

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