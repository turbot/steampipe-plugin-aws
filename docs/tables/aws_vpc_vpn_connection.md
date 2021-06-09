# Table: aws_vpc_vpn_connection

AWS Client VPN is a managed client-based VPN service that enables securely access AWS resources or on-premises network.

## Examples

### VPN connection basic info

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


### List VPN connection options info

```sql
select
  vpn_connection_id,
  options -> 'EnableAcceleration' as enable_acceleration,
  options ->> 'LocalIpv4NetworkCidr' as local_ipv4_network_cidr,
  options ->> 'LocalIpv6NetworkCidr' as local_ipv6_network_cidr,
  options ->> 'RemoteIpv4NetworkCidr' as remote_ipv4_network_cidr,
  options ->> 'RemoteIpv6NetworkCidr' as remote_ipv6_network_cidr,
  options -> 'StaticRoutesOnly' as static_routes_only,
  options ->> 'TunnelInsideIpVersion' as tunnel_inside_ip_iersion,
  options ->> '"TunnelOptions"' as "tunnel_options"
from
  aws_vpc_vpn_connection;
```


### List VPN tunnel info of each VPN connection

```sql
select
  vpn_connection_id,
  t -> 'AcceptedRouteCount' as accepted_route_count,
  t ->> 'CertificateArn' as certificate_arn,
  t ->> 'LastStatusChange' as last_status_change,
  t ->> 'OutsideIpAddress' as outside_ip_address,
  t ->> 'Status' as status,
  t ->> 'StatusMessage' as status_message
from
  aws_vpc_vpn_connection,
  jsonb_array_elements(vgw_telemetry) as t;
```