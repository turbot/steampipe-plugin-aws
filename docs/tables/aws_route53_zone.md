# Table: aws_route53_zone

A hosted zone is analogous to a traditional DNS zone file; it represents a collection of records that can be managed together, belonging to a single parent domain name. All resource record sets within a hosted zone must have the hosted zone’s domain name as a suffix.


## Examples

### Basic Zone Info
```sql
select
  name,
  id,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone;
```

### List private zones  
```sql
select
  name,
  id,
  comment,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  private_zone;
```


### List public zones  
```sql
select
  name,
  id,
  comment,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  not private_zone;
```


### Find zones by subdomain name
```sql
select
  name,
  id,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  name like '%.turbot.com.'
```
