# Table: aws_vpc_vpn_gateway

On the AWS side of the Site-to-Site VPN connection, a virtual private gateway or transit gateway provides two VPN endpoints (tunnels) for automatic failover. AWS Client VPN is a managed client-based VPN service that enables you to securely access your AWS resources or your on-premises network.

## Examples

### VPN gateways basic info

```sql
select
  vpn_gateway_id,
  state,
  type,
  amazon_side_asn,
  availability_zone,
  vpc_attachments
from
  aws_vpc_vpn_gateway;
```


### List Unattached VPN gateways

```sql
select
  vpn_gateway_id
from
  aws_vpc_vpn_gateway
where
  vpc_attachments is null;
```


### List all the VPN gateways attached to default VPC

```sql
select
  vpn_gateway_id,
  vpc.is_default
from
  aws_vpc_vpn_gateway
  cross join jsonb_array_elements(vpc_attachments) as i
  join aws_vpc vpc on i ->> 'VpcId' = vpc.vpc_id
where
  vpc.is_default = true;
```