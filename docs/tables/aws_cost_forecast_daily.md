# Table: aws_cost_forecast_daily

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_forecast_daily` table retrieves a forecast for how much Amazon Web Services predicts that you will spend each day over the next 4 months, based on your past costs.

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Examples

### Basic info

```sql

select 
   period_start,
   period_end,
   mean_value::numeric::money   
from 
  aws_cost_forecast_daily
order by
  period_start;
```

