# Table: aws_region

An AWS Region is a collection of AWS resources in a geographic area. Each AWS Region is isolated and independent of the other Regions. Regions provide fault tolerance, stability, and resilience, and can also reduce latency.

## Examples

### AWS region info

```sql
select
  name,
  opt_in_status
from
  aws_region;
```


### List of AWS regions which are enable

```sql
select
  name,
  opt_in_status
from
  aws_region
where
  opt_in_status = 'not-opted-in';
```
