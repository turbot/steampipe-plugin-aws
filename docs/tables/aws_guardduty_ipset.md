# Table: aws_guardduty_ipset

An IPSet is a list of trusted IP addresses from which secure communication is permitted with AWS infrastructure and applications.

## Examples

### Basic info

```sql
select
  detector_id,
  ipset_id,
  name,
  format,
  location
from
  aws_guardduty_ipset;
```


### List ipsets which are active

```sql
select
  ipset_id,
  name,
  status
from
  aws_guardduty_ipset
where
  status = 'ACTIVE';
```
