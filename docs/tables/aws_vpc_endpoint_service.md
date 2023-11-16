---
title: "Table: aws_vpc_endpoint_service - Query AWS VPC Endpoint Services using SQL"
description: "Allows users to query AWS VPC Endpoint Services to retrieve detailed information about each service, including service name, service type, and whether or not the service is private."
---

# Table: aws_vpc_endpoint_service - Query AWS VPC Endpoint Services using SQL

The `aws_vpc_endpoint_service` table in Steampipe provides information about AWS VPC Endpoint Services. This table allows DevOps engineers to query service-specific details, including service type, service name, and whether or not the service is private. Users can utilize this table to gather insights on services, such as identifying private services, understanding the types of services available, and more. The schema outlines the various attributes of the VPC Endpoint Service, including the service id, service name, service type, and whether or not the service is private.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_endpoint_service` table, you can use the `.inspect aws_vpc_endpoint_service` command in Steampipe.

### Key columns:

- `service_id`: The ID of the service. This can be used to join this table with other tables that contain information about AWS VPC Endpoint Services.
- `service_name`: The name of the service. This can be used to filter results based on the service name.
- `service_type`: The type of service. This can be used to categorize services based on their type.

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
