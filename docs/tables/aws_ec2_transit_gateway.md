---
title: "Steampipe Table: aws_ec2_transit_gateway - Query AWS EC2 Transit Gateway using SQL"
description: "Allows users to query AWS EC2 Transit Gateway resources for detailed information on configuration, status, and associations."
folder: "EC2"
---

# Table: aws_ec2_transit_gateway - Query AWS EC2 Transit Gateway using SQL

The AWS EC2 Transit Gateway is a service that simplifies the process of networking connectivity across multiple Amazon Virtual Private Clouds (VPCs) and on-premises networks. It acts as a hub that controls how traffic is routed among all connected networks which simplifies your network architecture. With Transit Gateway, you can manage connectivity for thousands of VPCs, easily scale connectivity across multiple AWS accounts, and segregate your network traffic to improve security.

## Table Usage Guide

The `aws_ec2_transit_gateway` table in Steampipe provides you with information about Transit Gateways within Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query Transit Gateway-specific details, including its configuration, state, and associations. You can utilize this table to gather insights on Transit Gateways, such as its attached VPCs, VPN connections, Direct Connect gateways, and more. The schema outlines the various attributes of the Transit Gateway for you, including the transit gateway ID, creation time, state, and associated tags.

## Examples

### Basic Transit Gateway info
Gain insights into the status and ownership details of your AWS Transit Gateway configurations, along with their creation times, to better manage your network transit connectivity. This can be particularly useful for auditing, tracking changes, and troubleshooting network issues.

```sql+postgres
select
  transit_gateway_id,
  state,
  owner_id,
  creation_time
from
  aws_ec2_transit_gateway;
```

```sql+sqlite
select
  transit_gateway_id,
  state,
  owner_id,
  creation_time
from
  aws_ec2_transit_gateway;
```


### List transit gateways which automatically accepts shared account attachment
Determine the areas in which transit gateways are set to automatically accept shared account attachments. This is useful to identify potential security risks and ensure proper management of your AWS resources.

```sql+postgres
select
  transit_gateway_id,
  auto_accept_shared_attachments
from
  aws_ec2_transit_gateway
where
  auto_accept_shared_attachments = 'enable';
```

```sql+sqlite
select
  transit_gateway_id,
  auto_accept_shared_attachments
from
  aws_ec2_transit_gateway
where
  auto_accept_shared_attachments = 'enable';
```


### Find the number of transit gateways by default route table id
Determine the areas in which transit gateways are most commonly associated by default route table ID, which can aid in understanding network traffic distribution and optimizing resource allocation within your AWS EC2 environment.

```sql+postgres
select
  association_default_route_table_id,
  count(transit_gateway_id) as transit_gateway
from
  aws_ec2_transit_gateway
group by
  association_default_route_table_id;
```

```sql+sqlite
select
  association_default_route_table_id,
  count(transit_gateway_id) as transit_gateway
from
  aws_ec2_transit_gateway
group by
  association_default_route_table_id;
```


### Map all transit gateways to the application to which they belong with an application tag
Discover the segments that have transit gateways without an application tag, enabling you to identify and categorize untagged resources for better resource management and organization.

```sql+postgres
select
  transit_gateway_id,
  tags
from
  aws_ec2_transit_gateway
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  transit_gateway_id,
  tags
from
  aws_ec2_transit_gateway
where
  json_extract(tags, '$.application') IS NULL;
```