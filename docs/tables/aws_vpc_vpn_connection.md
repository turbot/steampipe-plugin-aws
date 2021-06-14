# Table: aws_vpc_vpn_connection

A VPN connection can be used to configure secure access to your AWS resources with remote and on-premisis networks.

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
