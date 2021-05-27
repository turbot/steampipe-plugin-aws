# Table: aws_cost_by_service_monthly

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_monthly` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Basic info

```sql
select
  service,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_service_monthly
order by
  service,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by service

```sql
select
  service,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  service;
```

### Top 10 most expensive service (by average monthly unblended_cost_amount)

```sql
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  average desc
limit 10;
```


### Top 10 most expensive service (by total monthly unblended_cost_amount)

```sql
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  sum desc
limit 10;
```


### Ranked - Most expensive month (unblended_cost_amount) by service

```sql
with ranked_costs as (
  select
    service,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by service order by unblended_cost_amount desc)
  from 
    aws_cost_by_service_monthly
)
select * from ranked_costs where rank = 1
```



###  Month on month growth (unblended_cost_amount) by service

```sql
with cost_data as (
  select
    service,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount,1) over(partition by service order by period_start desc) as previous_month
  from 
    aws_cost_by_service_monthly
)
select
    service,
    period_start,
    this_month::numeric::money,
    previous_month::numeric::money,
    case 
      when previous_month = 0 and this_month = 0  then 0
      when previous_month = 0 then 999
      else round((100 * ( (this_month - previous_month) / previous_month))::numeric, 2) 
    end as percent_change
from
  cost_data
order by
  service,
  period_start;
```