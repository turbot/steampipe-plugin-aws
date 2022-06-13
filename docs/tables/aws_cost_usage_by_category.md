# Table: aws_cost_usage_by_category

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_usage_by_category` table provides a simplified yet flexible view of cost for your account by cost category against the organization master account.  You must specify a granularity (`MONTHLY`, `DAILY`), and you may pass category key or category value for filter out the cost for usage by category otherwise it will show all the cost by usage for all available category.

This tables requires an '=' qualifier for columns: granularity and "=" or "<>" qualifier for category_key and category_value columns.

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Monthly net unblended cost by account and service
```sql
select 
  period_start,
  category_key,
  category_value,
  net_unblended_cost_amount::numeric::money
from 
  aws_nagraj_master_acc.aws_cost_usage_by_category
where 
  granularity = 'MONTHLY'
and
  category_key = 'Test_Cost_Category'
order by
  category_key,
  period_start;
```

### Top 5 most expensive (net unblended cost) in each cost by category
```sql
with ranked_costs as (
  select 
    category_key,
    category_value,
    sum(net_unblended_cost_amount)::numeric::money as net_unblended_cost,
    rank() over(partition by category_key order by sum(net_unblended_cost_amount) desc)
  from 
    aws_nagraj_master_acc.aws_cost_usage_by_category
  where 
    granularity = 'MONTHLY'
  group by
    category_key,
    category_value
  order by
    category_key,
    net_unblended_cost desc
)
select * from ranked_costs where rank <= 5;
```

### Monthly net unblended cost by category

```sql
select 
  period_start,
  category_key,
  category_value,
  net_unblended_cost_amount::numeric::money
from 
  aws_nagraj_master_acc.aws_cost_usage_by_category
where 
  granularity = 'MONTHLY'
and
  category_key = 'nagaraj_category'
order by
  category_key,
  period_start;
```

### List monthly cost and usage by cost category value

```sql
select 
  period_start,
  category_key,
  category_value,
  net_unblended_cost_amount::numeric::money
from 
  aws_nagraj_master_acc.aws_cost_usage_by_category
where 
  granularity = 'MONTHLY'
and
  category_value = '49'
order by
  category_key,
  period_start;
```
