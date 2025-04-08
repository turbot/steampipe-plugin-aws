---
title: "Steampipe Table: aws_vpc_endpoint_service - Query AWS VPC Endpoint Services using SQL"
description: "Allows users to query AWS VPC Endpoint Services to retrieve detailed information about each service, including service name, service type, and whether or not the service is private."
folder: "VPC"
---

# Table: aws_vpc_endpoint_service - Query AWS VPC Endpoint Services using SQL

The AWS VPC Endpoint Service is a feature that allows private connectivity between your VPC and the service of another AWS account, without requiring access over the internet or through a VPN connection. This service enables you to expose your application or service behind load balancers inside your VPC to other AWS accounts. It also supports private connectivity over AWS Direct Connect, providing a more reliable and consistent network experience than internet-based connections.

## Table Usage Guide

The `aws_vpc_endpoint_service` table in Steampipe provides you with information about AWS VPC Endpoint Services. This table allows you, as a DevOps engineer, to query service-specific details, including service type, service name, and whether or not the service is private. You can utilize this table to gather insights on services, such as identifying private services, understanding the types of services available, and more. The schema outlines the various attributes of the VPC Endpoint Service for you, including the service id, service name, service type, and whether or not the service is private.

## Examples

### Basic info
Explore the various services within your AWS VPC by identifying their names, IDs, and associated DNS details. This can be useful for understanding the structure and connectivity of your VPC, particularly when troubleshooting or optimizing network configurations.

```sql+postgres
select
  service_name,
  service_id,
  base_endpoint_dns_names,
  private_dns_name
from
  aws_vpc_endpoint_service;
```

```sql+sqlite
select
  service_name,
  service_id,
  base_endpoint_dns_names,
  private_dns_name
from
  aws_vpc_endpoint_service;
```

### Get availability zone count for each VPC endpoint service
Discover the number of availability zones associated with each VPC endpoint service. This is useful in understanding the distribution and availability of services across different zones.

```sql+postgres
select
  service_name,
  jsonb_array_length(availability_zones) as availability_zone_count
from
  aws_vpc_endpoint_service;
```

```sql+sqlite
select
  service_name,
  json_array_length(json_extract(availability_zones, '
```

### Get DNS information for each VPC endpoint service
Discover the segments that consist of DNS details for each VPC endpoint service. This could be useful in managing network traffic and ensuring secure and efficient communication within your AWS environment.

```sql+postgres
select
  service_name,
  service_id,
  base_endpoint_dns_names,
  private_dns_name
from
  aws_vpc_endpoint_service;
```

```sql+sqlite
select
  service_name,
  service_id,
  base_endpoint_dns_names,
  private_dns_name
from
  aws_vpc_endpoint_service;
```

### List VPC endpoint services with their corresponding service types
Explore which VPC endpoint services are linked with their corresponding types. This can be useful in managing and optimizing your AWS VPC environment by understanding the association between services and their types.

```sql+postgres
select
  service_name,
  service_id,
  type ->> 'ServiceType' as service_type
from
  aws_vpc_endpoint_service
  cross join jsonb_array_elements(service_type) as type;
```

```sql+sqlite
select
  service_name,
  service_id,
  json_extract(type, '$.ServiceType') as service_type
from
  aws_vpc_endpoint_service,
  json_each(service_type) as type;
```

### List VPC endpoint services which do not support endpoint policies
Discover the segments that are not supported by VPC endpoint policies. This is useful to identify potential vulnerabilities and ensure compliance with security standards.

```sql+postgres
select
  service_name,
  service_id,
  vpc_endpoint_policy_supported
from
  aws_vpc_endpoint_service
where
  not vpc_endpoint_policy_supported;
```

```sql+sqlite
select
  service_name,
  service_id,
  vpc_endpoint_policy_supported
from
  aws_vpc_endpoint_service
where
  vpc_endpoint_policy_supported = 0;
```

### List allowed principals for each VPC endpoint services
Determine the areas in which specific permissions are allowed for each VPC endpoint service within your AWS environment. This can be especially useful for understanding and managing access control and security configurations.

```sql+postgres
select
  service_name,
  service_id,
  jsonb_pretty(vpc_endpoint_service_permissions) as allowed_principals
from
  aws_vpc_endpoint_service;
```

```sql+sqlite
select
  service_name,
  service_id,
  vpc_endpoint_service_permissions as allowed_principals
from
  aws_vpc_endpoint_service;
```

### Get VPC endpoint connection info for each VPC endpoint service
Explore the connection details for each VPC endpoint service to gain insights into their network load balancer associations and status. This can help in assessing the health and configuration of your VPC endpoint services.

```sql+postgres
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

```sql+sqlite
select
  service_name,
  service_id,
  json_extract(c.value, '$.VpcEndpointId') as vpc_endpoint_id,
  json_extract(c.value, '$.VpcEndpointOwner') as vpc_endpoint_owner,
  json_extract(c.value, '$.VpcEndpointState') as vpc_endpoint_state,
  json_each.value as network_loadBalancer_arns
from
  aws_vpc_endpoint_service,
  json_each(vpc_endpoint_connections) as c,
  json_each(json_extract(c.value, '$.NetworkLoadBalancerArns'))
```
