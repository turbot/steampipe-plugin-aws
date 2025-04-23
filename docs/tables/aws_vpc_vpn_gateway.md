---
title: "Steampipe Table: aws_vpc_vpn_gateway - Query AWS VPC VPN Gateway using SQL"
description: "Allows users to query AWS VPC VPN Gateway data, providing details about Virtual Private Cloud (VPC) VPN gateways in an AWS account."
folder: "VPC"
---

# Table: aws_vpc_vpn_gateway - Query AWS VPC VPN Gateway using SQL

The AWS VPC VPN Gateway is a component of Amazon Virtual Private Cloud (VPC) that enables the establishment of a secure and private tunnel from your network or device to the AWS global network. It provides connectivity between your virtual network and your on-premises or other cloud network. A VPN gateway is the VPN concentrator on the Amazon side of the VPN connection.

## Table Usage Guide

The `aws_vpc_vpn_gateway` table in Steampipe provides you with information about Virtual Private Cloud (VPC) VPN gateways within AWS. This table allows you as a DevOps engineer, developer, or data analyst to query VPN gateway-specific details, including the state of the VPN gateway, the type of VPN gateway, the availability zone, and the VPC attachments. You can utilize this table to gather insights on VPN gateways, such as the number of VPN gateways in a specific state, the types of VPN gateways used, and the VPCs to which they are attached. The schema outlines the various attributes of the VPN gateway for you, including the VPN gateway ID, the Amazon Resource Name (ARN), and the associated tags.

## Examples

### VPN gateways basic info
Explore the status and type of your VPN gateways within your Amazon Web Services environment. This can help you understand the current configuration and availability of your gateways, which is crucial for maintaining secure and efficient network connections.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have VPN gateways without any VPC attachments. This is useful for identifying unused resources and optimizing cloud infrastructure management.

```sql+postgres
select
  vpn_gateway_id
from
  aws_vpc_vpn_gateway
where
  vpc_attachments is null;
```

```sql+sqlite
select
  vpn_gateway_id
from
  aws_vpc_vpn_gateway
where
  vpc_attachments is null;
```


### List all the VPN gateways attached to default VPC
Explore which VPN gateways are connected to your default VPC. This is beneficial to understand your default network infrastructure and identify any potential security risks or misconfigurations.

```sql+postgres
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

```sql+sqlite
select
  vpn_gateway_id,
  vpc.is_default
from
  aws_vpc_vpn_gateway,
  json_each(vpc_attachments)
  join aws_vpc vpc on json_extract(value, '$.VpcId') = vpc.vpc_id
where
  vpc.is_default = 1;
```