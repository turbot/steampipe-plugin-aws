# Table: aws_ec2_transit_gateway_route_table

Transit gateway route tables are used to configure routing for transit gateway attachments.

## Examples

### Basic transit gateway route table info

```sql
select
  transit_gateway_route_table_id,
  transit_gateway_id,
  default_association_route_table,
  default_propagation_route_table
from
  aws_ec2_transit_gateway_route_table;
```


### Count of transit gateway route table by transit gateway

```sql
select
  transit_gateway_id,
  count(transit_gateway_route_table_id) as transit_gateway_route_table_count
from
  aws_ec2_transit_gateway_route_table
group by
  transit_gateway_id;
```