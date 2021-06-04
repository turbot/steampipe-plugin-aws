# Table: aws_cost_by_service_usage_type_monthly

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_usage_type_monthly` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Basic info

```sql
select
  service,
  usage_type,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_service_usage_type_monthly
order by
  service,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by service and usage type

```sql
select
  service,
  usage_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  service,
  usage_type;
```

### Top 10 most expensive service usage type (by average monthly unblended_cost_amount)

```sql
select
  service,
  usage_type,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  average desc
limit 10;
```


### Top 10 most expensive service usage type (by total monthly unblended_cost_amount)

```sql
select
  service,
  usage_type,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  sum desc
limit 10;
```
