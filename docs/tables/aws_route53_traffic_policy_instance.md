# Table: aws_route53_traffic_policy_instance

Amazon Route 53 Traffic Policy instance identifies the hosted zone where you want to create the policy record and the domain or subdomain name that you want to route the traffic for.

## Examples

### Basic Info

```sql
select
  name,
  id,
  hosted_zone_id,
  ttl,
  region
from 
  aws_route53_traffic_policy_instance;
```

### List associated hosted zone details for each instance

```sql
select 
  i.name,
  i.id,
  h.id as hosted_zone_id,
  h.name as hosted_zone_name,
  h.caller_reference,
  h.private_zone
from 
  aws_route53_traffic_policy_instance i
  join aws_route53_zone h on i.hosted_zone_id = h.id;
```

### List associated traffic policy details for each instance

```sql
select 
  i.name,
  i.id,
  traffic_policy_id,
  p.name as traffic_policy_name,
  traffic_policy_type,
  traffic_policy_version,
  p.document
from 
  aws_route53_traffic_policy_instance i
  join aws_route53_traffic_policy p on i.traffic_policy_id = p.id 
  and i.traffic_policy_version = p.version;
```

### List instances that failed creation

```sql
select
  name,
  id,
  state,
  hosted_zone_id,
  message as failed_reason
from 
  aws_route53_traffic_policy_instance
where
  state = 'Failed';
```