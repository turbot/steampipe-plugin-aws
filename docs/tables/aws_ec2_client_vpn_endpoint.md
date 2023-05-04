# Table: aws_ec2_client_vpn_endpoint

An AWS Client VPN Endpoint is a fully-managed VPN service offered by Amazon Web Services (AWS) that enables secure and private access to resources in a virtual private cloud (VPC) or on-premises network.

## Examples

### Basic Info

```sql
select
  title,
  description,
  status
  client_vpn_endpoint_id,
  transfer_protocol,
  creation_time,
  tags
from
  aws_ec2_client_vpn_endpoint;
```
