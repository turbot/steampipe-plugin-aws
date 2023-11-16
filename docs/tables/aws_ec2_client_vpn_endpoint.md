---
title: "Table: aws_ec2_client_vpn_endpoint - Query AWS EC2 Client VPN Endpoints using SQL"
description: "Allows users to query AWS EC2 Client VPN Endpoints to retrieve detailed information about the configuration, status, and associated network details of each endpoint."
---

# Table: aws_ec2_client_vpn_endpoint - Query AWS EC2 Client VPN Endpoints using SQL

The `aws_ec2_client_vpn_endpoint` table in Steampipe provides information about the Client VPN endpoints within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers, security analysts, and other IT professionals to query VPN endpoint-specific details, including the endpoint configuration, associated network details, connection logs, and associated metadata. Users can utilize this table to gather insights on VPN endpoints, such as the associated VPC, Subnets, Security Groups, and more. The schema outlines the various attributes of the VPN endpoint, including the endpoint ID, creation time, DNS server, VPN protocol, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_client_vpn_endpoint` table, you can use the `.inspect aws_ec2_client_vpn_endpoint` command in Steampipe.

### Key columns:

- `vpn_endpoint_id`: The ID of the VPN endpoint. This is the primary key of the table and can be used to join with other tables to get more detailed information.
- `vpc_id`: The ID of the VPC associated with the VPN endpoint. This can be used to join with the `aws_vpc` table to get more information about the associated VPC.
- `security_group_ids`: The IDs of the security groups associated with the VPN endpoint. These can be used to join with the `aws_security_group` table to get more information about the security groups.

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