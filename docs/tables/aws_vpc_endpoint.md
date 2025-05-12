---
title: "Steampipe Table: aws_vpc_endpoint - Query AWS VPC Endpoints using SQL"
description: "Allows users to query AWS VPC Endpoints and retrieve information about each endpoint's configuration, type, status, and related resources such as network interfaces, DNS entries, and security groups."
folder: "VPC"
---

# Table: aws_vpc_endpoint - Query AWS VPC Endpoints using SQL

The AWS VPC Endpoints allow private connectivity to services hosted in AWS, directly from your Virtual Private Cloud (VPC) and without the need for an Internet Gateway, VPN, Network Address Translation (NAT) devices, or firewall proxies. This service is primarily used to securely access AWS services without requiring an internet gateway in your VPC. They enhance the privacy and security of your VPC by not exposing it to the public internet.

## Table Usage Guide

The `aws_vpc_endpoint` table in Steampipe provides you with information about VPC Endpoints within Amazon Virtual Private Cloud (VPC). This table allows you, as a network administrator or DevOps engineer, to query endpoint-specific details, including its service configuration, type (Interface or Gateway), status, and associated resources such as network interfaces, DNS entries, and security groups. You can utilize this table to gather insights on VPC Endpoints, such as their accessibility, security configuration, and integration with other AWS services. The schema outlines the various attributes of the VPC Endpoint for you, including the endpoint ID, VPC ID, service name, and associated tags.

## Examples

### List of VPC endpoint and the corresponding services
Explore which services are associated with each VPC endpoint to better manage network traffic and enhance security in your AWS infrastructure. This could be particularly useful in identifying any misconfigurations or unnecessary connections.

```sql+postgres
select
  vpc_endpoint_id,
  vpc_id,
  service_name
from
  aws_vpc_endpoint;
```

```sql+sqlite
select
  vpc_endpoint_id,
  vpc_id,
  service_name
from
  aws_vpc_endpoint;
```

### Subnet Id count for each VPC endpoints
Explore the number of subnets associated with each VPC endpoint to better manage and organize your network infrastructure. This can aid in optimizing network performance and planning future network development.

```sql+postgres
select
  vpc_endpoint_id,
  jsonb_array_length(subnet_ids) as subnet_id_count
from
  aws_vpc_endpoint;
```

```sql+sqlite
select
  vpc_endpoint_id,
  json_array_length(subnet_ids) as subnet_id_count
from
  aws_vpc_endpoint;
```

### Network details for each VPC endpoint
Determine the areas in which specific network details for each VPC endpoint are configured. This information can be used to assess the network configuration and understand the relationships between different elements within the VPC.

```sql+postgres
select
  vpc_endpoint_id,
  vpc_id,
  jsonb_array_elements(subnet_ids) as subnet_ids,
  jsonb_array_elements(network_interface_ids) as network_interface_ids,
  jsonb_array_elements(route_table_ids) as route_table_ids,
  sg ->> 'GroupName' as sg_name
from
  aws_vpc_endpoint
  cross join jsonb_array_elements(groups) as sg;
```

```sql+sqlite
select
  vpc_endpoint_id,
  vpc_id,
  json_extract(subnet_id.value, '$') as subnet_ids,
  json_extract(network_interface_id.value, '$') as network_interface_ids,
  json_extract(route_table_id.value, '$') as route_table_ids,
  json_extract(sg.value, '$.GroupName') as sg_name
from
  aws_vpc_endpoint,
  json_each(subnet_ids) as subnet_id,
  json_each(network_interface_ids) as network_interface_id,
  json_each(route_table_ids) as route_table_id,
  json_each(groups) as sg;
```

### DNS information for the VPC endpoints
Determine the areas in which DNS is enabled for your VPC endpoints, allowing you to assess the elements within your network's private DNS configuration. This can help you manage and optimize your network infrastructure.

```sql+postgres
select
  vpc_endpoint_id,
  private_dns_enabled,
  dns ->> 'DnsName' as dns_name,
  dns ->> 'HostedZoneId' as hosted_zone_id
from
  aws_vpc_endpoint
  cross join jsonb_array_elements(dns_entries) as dns;
```

```sql+sqlite
select
  vpc_endpoint_id,
  private_dns_enabled,
  json_extract(dns.value, '$.DnsName') as dns_name,
  json_extract(dns.value, '$.HostedZoneId') as hosted_zone_id
from
  aws_vpc_endpoint,
  json_each(dns_entries) as dns;
```

### VPC endpoint count by VPC ID
Explore the number of VPC endpoints associated with each VPC ID to manage network traffic and enhance security within your AWS environment. This can be useful in identifying potential areas of congestion or security vulnerabilities.

```sql+postgres
select
  vpc_id,
  count(vpc_endpoint_id) as vpc_endpoint_count
from
  aws_vpc_endpoint
group by
  vpc_id;
```

```sql+sqlite
select
  vpc_id,
  count(vpc_endpoint_id) as vpc_endpoint_count
from
  aws_vpc_endpoint
group by
  vpc_id;
```

### Count endpoints by endpoint type

```sql
select
  vpc_endpoint_type,
  count(vpc_endpoint_id)
from
  aws_vpc_endpoint
group by
  vpc_endpoint_type;
```

### List 'interface' type VPC Endpoints

```sql
select
  vpc_endpoint_id,
  service_name,
  vpc_id,
  vpc_endpoint_type
from
  aws_vpc_endpoint
where
  vpc_endpoint_type = 'Interface';
```