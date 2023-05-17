# Table: aws_ec2_client_vpn_endpoint

An AWS Client VPN Endpoint is a fully-managed VPN service offered by Amazon Web Services (AWS) that enables secure and private access to resources in a virtual private cloud (VPC) or on-premises network.

## Examples

### Basic Info

```sql
select
  title,
  description,
  status,
  client_vpn_endpoint_id,
  transport_protocol,
  creation_time,
  tags
from
  aws_ec2_client_vpn_endpoint;
```

### List client VPN endpoints that are not in available state

```sql
select
  title,
  status,
  client_vpn_endpoint_id,
  transport_protocol,
  tags
from
  aws_ec2_client_vpn_endpoint
where
  status ->> 'Code' <> 'available';
```

### List client VPN endpoints created in the last 30 days

```sql
select
  title,
  status ->> 'Code' as status,
  client_vpn_endpoint_id,
  transport_protocol,
  tags
from
  aws_ec2_client_vpn_endpoint
where
  creation_time >= now() - interval '30' day;
```

### Get the security group and the VPC details of client VPN endpoints 

```sql
select
  title,
  status ->> 'Code' as status,
  client_vpn_endpoint_id,
  security_group_ids,
  vpc_id,
  vpn_port,
  vpn_protocol,
  transport_protocol,
  tags
from
  aws_ec2_client_vpn_endpoint
where
  creation_time >= now() - interval '30' day;
```

### Get the security group and the VPC details of client VPN endpoints 

```sql
select
  title,
  status ->> 'Code' as status,
  client_vpn_endpoint_id,
  security_group_ids,
  vpc_id,
  vpn_port,
  vpn_protocol,
  transport_protocol,
  tags
from
  aws_ec2_client_vpn_endpoint;
```

### Get the logging configuration of client VPN endpoints 

```sql
select
  title,
  status ->> 'Code' as status,
  client_vpn_endpoint_id,
  client_log_options ->> 'Enabled' as client_log_options_enabled,
  client_log_options ->> 'CloudwatchLogGroup' as client_log_options_cloudwatch_log_group,
  client_log_options ->> 'CloudwatchLogStream' as client_log_options_cloudwatch_log_stream,
  tags
from
  aws_ec2_client_vpn_endpoint;
```

### Get the authentication information of client VPN endpoints 

```sql
select
  title,
  status ->> 'Code' as status,
  client_vpn_endpoint_id,
  autentication ->> 'Type' as authentication_options_type,
  autentication -> 'MutualAuthentication' ->> 'ClientRootCertificateChain' as authentication_client_root_certificate_chain,
  authentication_options,
  tags
from
  aws_ec2_client_vpn_endpoint,
  jsonb_array_elements(authentication_options) as autentication;
```