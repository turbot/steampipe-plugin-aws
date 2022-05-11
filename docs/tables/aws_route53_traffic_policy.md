# Table: aws_route53_traffic_policy

Amazon Route 53 Traffic Policy is a domain name system service that allows an Amazon Web Services customer to define how end-user traffic is routed to application endpoints through a visual interface.

## Examples

### Basic Info

```sql
select
  name,
  id,
  version,
  document,
  region
from 
  aws_route53_traffic_policy;
```

### List policies' latest version

```sql
select 
  name,
  policy.id,
  policy.version, 
  comment 
from 
  aws_route53_traffic_policy policy,
  (select
    id,
    max(version) as version
  from 
    aws_route53_traffic_policy 
  group by 
    id) as latest
where 
  latest.id = policy.id 
  and latest.version = policy.version;
```

### List total policies in each dns type

```sql
select
  document ->> 'RecordType' as dns_type,
  count(id) as "policies"
from
  aws_route53_traffic_policy
group by 
  dns_type;
```