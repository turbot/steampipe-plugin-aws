---
title: "Table: aws_vpc_vpn_connection - Query AWS VPC VPN Connection using SQL"
description: "Allows users to query VPN connections in an AWS VPC."
---

# Table: aws_vpc_vpn_connection - Query AWS VPC VPN Connection using SQL

The `aws_vpc_vpn_connection` table in Steampipe provides information about VPN connections within AWS Virtual Private Cloud (VPC). This table allows DevOps engineers to query VPN connection-specific details, including the VPN connection ID, state, VPN gateway configurations, customer gateway configurations, and associated metadata. Users can utilize this table to gather insights on VPN connections, such as connection states, associated VPN and customer gateways, static routes, and more. The schema outlines the various attributes of the VPN connection, including the VPN connection ID, creation time, VPN gateway ID, customer gateway ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_vpn_connection` table, you can use the `.inspect aws_vpc_vpn_connection` command in Steampipe.

**Key columns**:

- `vpn_connection_id`: The ID of the VPN connection. This can be used to join this table with other tables that contain VPN connection information.
- `vpn_gateway_id`: The ID of the VPN gateway at the AWS side of the VPN connection. This can be used to join this table with other tables that contain VPN gateway information.
- `customer_gateway_id`: The ID of the customer gateway at your end of the VPN connection. This can be used to join this table with other tables that contain customer gateway information.

## Examples

### Basic info

```sql
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

```sql
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


### List VPN connections with tunnel status UP

```sql
select
  vpn_connection_id,
  arn,
  t ->> 'Status' as status
from
  aws_vpc_vpn_connection,
  jsonb_array_elements(vgw_telemetry) as t
where  t ->> 'Status' = 'UP';
```
