# Table: aws_vpc_verified_access_endpoint

A Verified Access endpoint represents an application. Each endpoint is associated with a Verified Access group and inherits the access policy for the group. You can optionally attach an application-specific endpoint policy to each endpoint.

## Examples

### Basic info

```sql
select
  verified_access_endpoint_id,
  verified_access_instance_id,
  verified_access_group_id,
  creation_time,
  verified_access_instance_id,
  domain_certificate_arn,
  device_validation_domain,
  status_code
from
  aws_vpc_verified_access_endpoint;
```

### List endpoints older than 30 days

```sql
select
  verified_access_endpoint_id,
  creation_time,
  description,
  status_code
from
  aws_vpc_verified_access_endpoint
where
  creation_time <= now() - interval '30' day;
```

### List endpoints that are not in active state

```sql
select
  verified_access_endpoint_id,
  status_code,
  creation_time,
  deletion_time,
  description,
  device_validation_domain
from
  aws_vpc_verified_access_endpoint
where
  status_code <> 'active';
```

### Get group details for each endpoint

```sql
select
  e.verified_access_endpoint_id,
  e.creation_time,
  i.creation_time as group_create_time,
  i.verified_access_group_id
from
  aws_vpc_verified_access_endpoint as e,
  aws_vpc_verified_access_group as i
where
  g.verified_access_group_id = i.verified_access_group_id;
```

### Get trusted provider details for each endpoint

```sql
select
  e.verified_access_group_id,
  e.creation_time,
  i.creation_time as instance_create_time,
  i.verified_access_instance_id,
  jsonb_pretty(i.verified_access_trust_providers) as verified_access_trust_providers
from
  aws_vpc_verified_access_endpoint as e,
  aws_vpc_verified_access_instance as i
where
  e.verified_access_instance_id = i.verified_access_instance_id;
```

### Count of endpoints per instance

```sql
select
  verified_access_instance_id,
  count(verified_access_endpoint_id) instance_count
from
  aws_vpc_verified_access_endpoint
group by
  verified_access_instance_id;
```

### Get network interface details for each endpoint

```sql
select
  e.verified_access_endpoint_id,
  i.network_interface_id,
  i.interface_type,
  i.private_ip_address,
  i.association_public_ip,
  jsonb_pretty(i.groups) as security_groups
from
  aws_vpc_verified_access_endpoint as e,
  aws_ec2_network_interface as i
where
  e.network_interface_options ->> 'NetworkInterfaceId' = i.network_interface_id;
```