# Table: aws_cost_by_service_daily

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_daily` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by day, for the last year.  


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
  aws_cost_by_service_daily
order by
  service,
  period_start;
```



### Min, Max, and average daily unblended_cost_amount by service

```sql
select
  service,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_daily
group by
  service
order by
  service;
```

### Top 10 most expensive service (by average daily unblended_cost_amount)

```sql
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_daily
group by
  service
order by
  average desc
limit 10;
```


### Top 10 most expensive service (by total daily unblended_cost_amount)

```sql
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_daily
group by
  service
order by
  sum desc
limit 10;
```


### Ranked - Top 10 Most expensive days (unblended_cost_amount) by service

```sql
with ranked_costs as (
  select
    service,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by service order by unblended_cost_amount desc)
  from 
    aws_cost_by_service_daily
)
select * from ranked_costs where rank <= 10
```
