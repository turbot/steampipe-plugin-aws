---
title: "Steampipe Table: aws_vpc_vpn_connection - Query AWS VPC VPN Connection using SQL"
description: "Allows users to query VPN connections in an AWS VPC."
folder: "VPC"
---

# Table: aws_vpc_vpn_connection - Query AWS VPC VPN Connection using SQL

The AWS VPC VPN Connection is a component within the Amazon Virtual Private Cloud (VPC) service that allows you to securely link your AWS environment with your on-premises networks. It utilizes industry-standard encryption protocols to provide high security and low latency for your network traffic. This VPN connection is a vital tool for hybrid cloud architectures, enabling seamless and secure communication between AWS and your data center.

## Table Usage Guide

The `aws_vpc_vpn_connection` table in Steampipe provides you with information about VPN connections within AWS Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query VPN connection-specific details, including the VPN connection ID, state, VPN gateway configurations, customer gateway configurations, and associated metadata. You can utilize this table to gather insights on VPN connections, such as connection states, associated VPN and customer gateways, static routes, and more. The schema outlines the various attributes of the VPN connection for you, including the VPN connection ID, creation time, VPN gateway ID, customer gateway ID, and associated tags.

## Examples

### Basic info
Explore the status and types of Virtual Private Network connections within your Amazon Web Services Virtual Private Cloud. This information is useful to understand the connectivity between your network and the AWS network, helping in maintaining secure and reliable connections.

```sql+postgres
select
  vpn_connection_id,
  state,
  type,
  vpn_gateway_id,
  customer_gateway_id,
  region
from
  aws_vpc_vpn_connection;
```

```sql+sqlite
select
  vpn_connection_id,
  state,
  type,
  vpn_gateway_id,
  customer_gateway_id,
  region
from
  aws_vpc_vpn_connection;
```


### Get option configurations for each VPN connection
Explore the configuration settings of each VPN connection to understand its specific features such as acceleration enablement, local and remote network details, and tunnel options. This can assist in optimizing network performance and security measures.

```sql+postgres
select
  vpn_connection_id,
  options -> 'EnableAcceleration' as enable_acceleration,
  options ->> 'LocalIpv4NetworkCidr' as local_ipv4_network_cidr,
  options ->> 'LocalIpv6NetworkCidr' as local_ipv6_network_cidr,
  options ->> 'RemoteIpv4NetworkCidr' as remote_ipv4_network_cidr,
  options ->> 'RemoteIpv6NetworkCidr' as remote_ipv6_network_cidr,
  options -> 'StaticRoutesOnly' as static_routes_only,
  options ->> 'TunnelInsideIpVersion' as tunnel_inside_ip_version,
  options ->> 'TunnelOptions' as tunnel_options
from
  aws_vpc_vpn_connection;
```

```sql+sqlite
select
  vpn_connection_id,
  json_extract(options, '$.EnableAcceleration') as enable_acceleration,
  json_extract(options, '$.LocalIpv4NetworkCidr') as local_ipv4_network_cidr,
  json_extract(options, '$.LocalIpv6NetworkCidr') as local_ipv6_network_cidr,
  json_extract(options, '$.RemoteIpv4NetworkCidr') as remote_ipv4_network_cidr,
  json_extract(options, '$.RemoteIpv6NetworkCidr') as remote_ipv6_network_cidr,
  json_extract(options, '$.StaticRoutesOnly') as static_routes_only,
  json_extract(options, '$.TunnelInsideIpVersion') as tunnel_inside_ip_version,
  json_extract(options, '$.TunnelOptions') as tunnel_options
from
  aws_vpc_vpn_connection;
```


### List VPN connections with tunnel status UP
This query is used to identify all active VPN connections within your AWS VPC. It's useful for maintaining a real-time overview of your network's connectivity status, helping to ensure secure and uninterrupted data transmission.

```sql+postgres
select
  vpn_connection_id,
  arn,
  t ->> 'Status' as status
from
  aws_vpc_vpn_connection,
  jsonb_array_elements(vgw_telemetry) as t
where  t ->> 'Status' = 'UP';
```

```sql+sqlite
select
  vpn_connection_id,
  arn,
  json_extract(t.value, '$.Status') as status
from
  aws_vpc_vpn_connection,
  json_each(vgw_telemetry) as t
where  json_extract(t.value, '$.Status') = 'UP';
```