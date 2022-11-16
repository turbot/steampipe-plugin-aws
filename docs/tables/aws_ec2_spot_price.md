# Table: aws_ec2_spot_price

Returns the list of prices for Spot EC2 instances.

## Examples

### List EC2 spot prices for Linux m5.4xlarge instance in eu-west-3a and eu-west-3b availability zones in the last month

```sql
select
  availability_zone,
  instance_type,
  product_description,
  spot_price::numeric as spot_price,
  create_timestamp as start_time,
  lead(create_timestamp, 1, now()) over (partition by instance_type, availability_zone, product_description order by create_timestamp) as stop_time
from
  aws_ec2_spot_price
where
  instance_type = 'm5.4xlarge'
  and product_description = 'Linux/UNIX'
  and availability_zone in
  (
    'eu-west-3a',
    'eu-west-3b'
  )
  and start_time = now() - interval '1' month
  and end_time = now() - interval '1' minute
```
