---
title: "Steampipe Table: aws_ec2_client_vpn_endpoint - Query AWS EC2 Client VPN Endpoints using SQL"
description: "Allows users to query AWS EC2 Client VPN Endpoints to retrieve detailed information about the configuration, status, and associated network details of each endpoint."
folder: "EC2"
---

# Table: aws_ec2_client_vpn_endpoint - Query AWS EC2 Client VPN Endpoints using SQL

The AWS EC2 Client VPN Endpoint is a scalable, end-to-end managed VPN service that enables users to securely access their AWS resources and home network. It provides secure and scalable compute capacity in the AWS Cloud, allowing users to launch virtual servers. With EC2 Client VPN, you can access your resources from any location using an OpenVPN-based VPN client.

## Table Usage Guide

The `aws_ec2_client_vpn_endpoint` table in Steampipe provides you with information about the Client VPN endpoints within AWS Elastic Compute Cloud (EC2). This table enables you, as a DevOps engineer, security analyst, or other IT professional, to query VPN endpoint-specific details, including the endpoint configuration, associated network details, connection logs, and associated metadata. You can utilize this table to gather insights on VPN endpoints, such as the associated VPC, Subnets, Security Groups, and more. The schema outlines the various attributes of the VPN endpoint for you, including the endpoint ID, creation time, DNS server, VPN protocol, and associated tags.

## Examples

### Basic Info
Explore the status and configuration details of your AWS EC2 Client VPN endpoints to understand their operational state and settings. This can be beneficial for assessing your network's security posture and troubleshooting connectivity issues.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which your client VPN endpoints are not available. This can be useful for troubleshooting connectivity issues or managing network resources.

```sql+postgres
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

```sql+sqlite
select
  title,
  status,
  client_vpn_endpoint_id,
  transport_protocol,
  tags
from
  aws_ec2_client_vpn_endpoint
where
  json_extract(status, '$.Code') <> 'available';
```

### List client VPN endpoints created in the last 30 days
Determine the areas in which new client VPN endpoints have been established in the past month. This can help manage and monitor recent network expansions or changes.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(status, '$.Code') as status,
  client_vpn_endpoint_id,
  transport_protocol,
  tags
from
  aws_ec2_client_vpn_endpoint
where
  creation_time >= datetime('now', '-30 day');
```

### Get the security group and the VPC details of client VPN endpoints
Determine the security setup of recently created VPN endpoints, including their associated security groups and VPC details. This is useful for reviewing and auditing the security configurations of new VPN connections in your network.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(status, '$.Code') as status,
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
  creation_time >= datetime('now', '-30 day');
```

### Get the security group and the VPC details of client VPN endpoints
Explore the security settings and network details of your client VPN endpoints. This can help in assessing the security measures in place and understanding the network configuration, which is crucial for maintaining a secure and efficient VPN service.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(status, '$.Code') as status,
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
Determine the status of client VPN endpoints and assess whether their logging configurations are enabled. This can be useful for monitoring and troubleshooting VPN connectivity issues.

```sql+postgres
select
  title,
  status ->> 'Code' as status,
  client_vpn_endpoint_id,
  connection_log_options ->> 'Enabled' as connection_log_options_enabled,
  connection_log_options ->> 'CloudwatchLogGroup' as connection_log_options_cloudwatch_log_group,
  connection_log_options ->> 'CloudwatchLogStream' as connection_log_options_cloudwatch_log_stream,
  tags
from
  aws_ec2_client_vpn_endpoint;
```

```sql+sqlite
select
  title,
  json_extract(status, '$.Code') as status,
  client_vpn_endpoint_id,
  json_extract(connection_log_options, '$.Enabled') as connection_log_options_enabled,
  json_extract(connection_log_options, '$.CloudwatchLogGroup') as connection_log_options_cloudwatch_log_group,
  json_extract(connection_log_options, '$.CloudwatchLogStream') as connection_log_options_cloudwatch_log_stream,
  tags
from
  aws_ec2_client_vpn_endpoint;
```

### Get the authentication information of client VPN endpoints
This query is used to gain insights into the authentication information of client VPN endpoints within the AWS EC2 service. It's particularly useful for understanding the type of authentication being used and the details of the mutual authentication, which can help in assessing security measures and compliance requirements.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(status, '$.Code') as status,
  client_vpn_endpoint_id,
  json_extract(autentication.value, '$.Type') as authentication_options_type,
  json_extract(json_extract(autentication.value, '$.MutualAuthentication'), '$.ClientRootCertificateChain') as authentication_client_root_certificate_chain,
  authentication_options,
  tags
from
  aws_ec2_client_vpn_endpoint,
  json_each(authentication_options) as autentication;
```