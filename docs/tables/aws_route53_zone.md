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

### List VPCs associated with zones

```sql
select 
  name,
  id,
  v ->> 'VPCId' as vpc_id,
  v ->> 'VPCRegion' as vpc_region
from
  aws_route53_zone,
  jsonb_array_elements(vpcs) as v;
```

### Get VPC details associated with zones

```sql
select 
  name,
  id,
  v.vpc_id as vpc_id,
  v.cidr_block as cidr_block,
  v.is_default as is_default,
  v.dhcp_options_id as dhcp_options_id
from
  aws_route53_zone,
  jsonb_array_elements(vpcs) as p,
  aws_vpc as v
where
  p ->> 'VPCId' = v.vpc_id;
```