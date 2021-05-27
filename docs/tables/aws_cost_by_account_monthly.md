# Table: aws_cost_by_account_monthly

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_account_monthly` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

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
  aws_cost_by_account_monthly
order by
  linked_account_id,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by account

```sql
select
  linked_account_id,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_account_monthly
group by
  linked_account_id
order by
  linked_account_id;
```


### Ranked - Most expensive months (unblended_cost_amount) by account

```sql
select
  linked_account_id,
  period_start,
  unblended_cost_amount::numeric::money,
  rank() over(partition by linked_account_id order by unblended_cost_amount desc)
from 
  aws_cost_by_account_monthly;
```



### Month on month growth (unblended_cost_amount) by account

```sql
with cost_data as (
  select
    linked_account_id,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount,1) over(partition by linked_account_id order by period_start desc) as previous_month
  from 
    aws_cost_by_account_monthly
)
select
    linked_account_id,
    period_start,
    this_month::numeric::money,
    previous_month::numeric::money,
    round((100 * ( (this_month - previous_month) / previous_month))::numeric, 2) as percent_change
from
  cost_data
order by
  linked_account_id,
  period_start;
```

