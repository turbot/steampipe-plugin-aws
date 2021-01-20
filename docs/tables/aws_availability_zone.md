# Table: aws_availability_zone

An Availability Zone (AZ) is one or more discrete data centers with redundant power, networking, and connectivity in an AWS Region.

## Examples

### Availability zone info

```sql
select
  name,
  zone_id,
  zone_type,
  group_name,
  region_name
from
  aws_availability_zone;
```


### Count of availability zone per region

```sql
select
  region_name,
  count(name) as zone_count
from
  aws_availability_zone
group by
  region_name;
```


### List of AWS availability zones which are not enabled in the account

```sql
select
  name,
  zone_id,
  region_name,
  opt_in_status
from
  aws_availability_zone
where
  opt_in_status = 'not-opted-in';
```
