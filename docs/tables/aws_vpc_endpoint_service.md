# Table: aws_vpc_endpoint_service

A VPC endpoint enables you to privately connect your VPC to supported AWS services and VPC endpoint services powered by AWS PrivateLink without requiring an internet gateway, NAT device, VPN connection, or AWS Direct Connect connection.

## Examples

### Basic info

```sql
select
  service_name,
  service_id,
  base_endpoint_dns_names,
  private_dns_name
from
  aws_vpc_endpoint_service;
```

### Get availability zone count for each VPC endpoint service

```sql
select
  service_name,
  jsonb_array_length(availability_zones) as availability_zone_count
from
  aws_vpc_endpoint_service;
```

### Get DNS information for each VPC endpoint service

```sql
select
  service_name,
  service_id,
  base_endpoint_dns_names,
  private_dns_name
from
  aws_vpc_endpoint_service;
```

### List VPC endpoint services with their corresponding service types

```sql
select
  service_name,
  service_id,
  type ->> 'ServiceType' as service_type
from
  aws_vpc_endpoint_service
  cross join jsonb_array_elements(service_type) as type;
```

### List VPC endpoint services which do not support endpoint policies

```sql
select
  service_name,
  service_id,
  vpc_endpoint_policy_supported
from
  aws_vpc_endpoint_service
where
  not vpc_endpoint_policy_supported;
```

### List allowed principals for each VPC endpoint services

```sql
select
  service_name,
  service_id,
  jsonb_pretty(vpc_endpoint_service_permissions) as allowed_principals
from
  aws_vpc_endpoint_service;
```

### Get VPC endpoint connection info for each VPC endpoint service

```sql
select
  service_name,
  service_id,
  c ->> 'VpcEndpointId' as vpc_endpoint_id,
  c ->> 'VpcEndpointOwner' as vpc_endpoint_owner,
  c ->> 'VpcEndpointState' as vpc_endpoint_state,
  jsonb_array_elements_text(c -> 'NetworkLoadBalancerArns') as network_loadBalancer_arns
from
  aws_vpc_endpoint_service,
  jsonb_array_elements(vpc_endpoint_connections) as c
```
