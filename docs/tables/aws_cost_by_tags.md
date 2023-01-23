# Table: aws_cost_by_tags

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_tags` table provides a simplified view of cost by tags in your account.

User must need to provide one of the tag key. Otherwise the table will return an empty row.

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Basic info

```sql
select
  tag_key_1,
  tag_value_1,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from
  aws_cost_by_tags
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name';
```

### Min, Max, and average daily unblended_cost_amount by tag

```sql
select
  tag_key_1,
  tag_value_1,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from
  aws_cost_by_tags
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name'
group by
  tag_key_1, tag_value_1;
```

### Ranked - Top 10 Most expensive days (unblended_cost_amount) by tag

```sql
with ranked_costs as (
  select
    tag_key_1,
    tag_value_1,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by tag_key_1 order by unblended_cost_amount desc)
  from
  aws_cost_by_tags
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name'
)
select * from ranked_costs where rank <= 10;
```
