# Table: aws_cost_by_account_daily

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_account_daily` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master), summarized by day, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Basic info

```sql
select
  linked_account_id,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_account_daily
order by
  linked_account_id,
  period_start;
```



### Min, Max, and average daily unblended_cost_amount by account

```sql
select
  linked_account_id,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_account_daily
group by
  linked_account_id
order by
  linked_account_id;
```


### Ranked - Top 10 Most expensive days (unblended_cost_amount) by account

```sql
with ranked_costs as (
  select
    linked_account_id,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by linked_account_id order by unblended_cost_amount desc)
  from 
    aws_cost_by_account_daily
)
select * from ranked_costs where rank <= 10
```
