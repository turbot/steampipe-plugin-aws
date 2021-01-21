# Table: aws_ec2_instance_availability

Available instance types in each region

## Examples

### List of instance types available in us-east-1 region

```sql
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'us-east-1';
```


### Check if r5.12xlarge instance type available in af-south-1

```sql
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'af-south'
  and instance_type = 'r5.12xlarge';
```
