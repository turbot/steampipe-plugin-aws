# Table: aws_ec2_transit_gateway_vpc_attachment

Attaching transit gateway to VPC.

## Examples

### Basic transit gateway vpc attachment info

```sql
select
  transit_gateway_attachment_id,
  transit_gateway_id,
  state,
  transit_gateway_owner_id,
  creation_time,
  association_state
from
  aws_ec2_transit_gateway_vpc_attachment;
```


### Count of transit gateway vpc attachment by transit gateway id

```sql
select
  resource_type,
  count(transit_gateway_attachment_id) as count
from
  aws_ec2_transit_gateway_vpc_attachment
group by
  resource_type;
```