# Table: aws_dax_subnet_group

Amazon DynamoDB Accelerator (DAX) is a fully managed, highly available, in-memory cache for Amazon DynamoDB that delivers up to a 10 times performance improvement—from milliseconds to microseconds—even at millions of requests per second.

## Examples

### Basic info

```sql
select
  subnet_group_name,
  description,
  vpc_id,
  subnets,
  region
from
  aws_dax_subnet_group;
```

### List VPC details for each subnet group

```sql
select
  subnet_group_name,
  v.vpc_id,
  v.arn as vpc_arn,
  v.cidr_block as vpc_cidr_block,
  v.state as vpc_state,
  v.is_default as is_default_vpc,
  v.region
from
  aws_dax_subnet_group g
join aws_vpc v
  on v.vpc_id = g.vpc_id;
```

### List subnet details for each subnet group

```sql
select
  subnet_group_name,
  g.vpc_id,
  vs.subnet_arn,
  vs.cidr_block as subnet_cidr_block,
  vs.state as subnet_state,
  vs.availability_zone as subnet_availability_zone,
  vs.region
from
  aws_dax_subnet_group g,
  jsonb_array_elements(subnets) s
join aws_vpc_subnet vs
  on vs.subnet_id = s ->> 'SubnetIdentifier';
```