# Table: aws_vpc_nat_gateway_metric_bytes_out_to_destination

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_vpc_nat_gateway_metric_bytes_out_to_destination` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info

```sql
select
  nat_gateway_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_vpc_nat_gateway_metric_bytes_out_to_destination
order by
  nat_gateway_id,
  timestamp;
```

### Show unused nat gateways

```sql
select
  g.nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway as g
  left join aws_vpc_nat_gateway_metric_bytes_out_to_destination as d
  on g.nat_gateway_id = d.nat_gateway_id
group by
  g.nat_gateway_id,
  vpc_id,
  subnet_id
having
  sum(average) = 0;
```
