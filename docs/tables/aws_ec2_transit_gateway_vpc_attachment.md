---
title: "Steampipe Table: aws_ec2_transit_gateway_vpc_attachment - Query AWS EC2 Transit Gateway VPC Attachments using SQL"
description: "Allows users to query AWS EC2 Transit Gateway VPC Attachments for details such as the attachment state, creation time, and more."
folder: "EC2"
---

# Table: aws_ec2_transit_gateway_vpc_attachment - Query AWS EC2 Transit Gateway VPC Attachments using SQL

The AWS EC2 Transit Gateway VPC Attachment is a resource that allows you to attach an Amazon VPC to a transit gateway. This attachment enables connectivity between the VPC and other networks connected to the transit gateway. It simplifies network architecture, reduces operational overhead, and provides a central gateway for connectivity.

## Table Usage Guide

The `aws_ec2_transit_gateway_vpc_attachment` table in Steampipe provides you with information about the attachments between Virtual Private Clouds (VPCs) and transit gateways in Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query attachment-specific details, including the attachment state, creation time, and associated metadata. You can utilize this table to gather insights on attachments, such as their status, the VPCs they are associated with, the transit gateways they are connected to, and more. The schema outlines the various attributes of the transit gateway VPC attachment for you, including the attachment ID, transit gateway ID, VPC ID, and associated tags.

## Examples

### Basic transit gateway vpc attachment info
Determine the areas in which your AWS EC2 Transit Gateway is attached to a VPC. This helps you understand the status and ownership of these connections, as well as when they were created.

```sql+postgres
select
  transit_gateway_attachment_id,
  transit_gateway_id,
  state,
  transit_gateway_owner_id,
  creation_time,
  association_state
from
  aws_ec2_transit_gateway_vpc_attachment;
```

```sql+sqlite
select
  transit_gateway_attachment_id,
  transit_gateway_id,
  state,
  transit_gateway_owner_id,
  creation_time,
  association_state
from
  aws_ec2_transit_gateway_vpc_attachment;
```


### Count of transit gateway vpc attachment by transit gateway id
Analyze your AWS EC2 Transit Gateway setup to understand the distribution of VPC attachments across different types of resources. This could be useful in optimizing resource allocation and identifying potential areas for cost savings.

```sql+postgres
select
  resource_type,
  count(transit_gateway_attachment_id) as count
from
  aws_ec2_transit_gateway_vpc_attachment
group by
  resource_type;
```

```sql+sqlite
select
  resource_type,
  count(transit_gateway_attachment_id) as count
from
  aws_ec2_transit_gateway_vpc_attachment
group by
  resource_type;
```