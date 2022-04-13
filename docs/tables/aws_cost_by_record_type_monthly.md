# Table: aws_cost_by_record_type_monthly

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_record_type_monthly` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master) as per record types (fees, usage, costs, tax refunds, and credits), summarized by month, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Basic info

```sql
select
  linked_account_id,
  record_type,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_record_type_monthly
order by
  linked_account_id,
  period_start;
```

### Min, Max, and average monthly unblended_cost_amount by account and record type

```sql
select
  linked_account_id,
  record_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_record_type_monthly
group by
  linked_account_id,
  record_type
order by
  linked_account_id;
```

### Ranked - Most expensive months (unblended_cost_amount) by account and record type

```sql
select
  linked_account_id,
  record_type,
  period_start,
  unblended_cost_amount::numeric::money,
  rank() over(partition by linked_account_id, record_type order by unblended_cost_amount desc)
from 
  aws_cost_by_record_type_monthly;
```