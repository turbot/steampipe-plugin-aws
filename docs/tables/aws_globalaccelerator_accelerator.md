# Table: aws_globalaccelerator_accelerator

An accelerator directs traffic to endpoints over the AWS global network to improve the performance of your internet
applications. Each accelerator includes one or more listeners.

## Examples

### Basic info

```sql
select
  name,
  created_time,
  dns_name,
  enabled,
  ip_address_type,
  last_modified_time,
  status
from
  aws_globalaccelerator_accelerator;
```

### List IPs used by global accelerators

```sql
 select
   name,
   created_time,
   dns_name,
   enabled,
   ip_address_type,
   last_modified_time,
   status,
   anycast_ip
from
  aws_globalaccelerator_accelerator,
  jsonb_array_elements(ip_sets -> 0 -> 'IpAddresses') as anycast_ip;
```

### List global accelerators without owner tag key

```sql
select
  name,
  tags
from
  aws_globalaccelerator_accelerator
where
  not tags::JSONB ? 'owner';
```
