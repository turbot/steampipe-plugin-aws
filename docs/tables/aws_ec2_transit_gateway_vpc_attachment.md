---
title: "Table: aws_ec2_transit_gateway_vpc_attachment - Query AWS EC2 Transit Gateway VPC Attachments using SQL"
description: "Allows users to query AWS EC2 Transit Gateway VPC Attachments for details such as the attachment state, creation time, and more."
---

# Table: aws_ec2_transit_gateway_vpc_attachment - Query AWS EC2 Transit Gateway VPC Attachments using SQL

The `aws_ec2_transit_gateway_vpc_attachment` table in Steampipe provides information about the attachments between Virtual Private Clouds (VPCs) and transit gateways in Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query attachment-specific details, including the attachment state, creation time, and associated metadata. Users can utilize this table to gather insights on attachments, such as their status, the VPCs they are associated with, the transit gateways they are connected to, and more. The schema outlines the various attributes of the transit gateway VPC attachment, including the attachment ID, transit gateway ID, VPC ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_transit_gateway_vpc_attachment` table, you can use the `.inspect aws_ec2_transit_gateway_vpc_attachment` command in Steampipe.

### Key columns:

- `transit_gateway_attachment_id`: This is the unique identifier for the transit gateway attachment. It can be used to join this table with other tables to get more detailed information about the specific attachment.
- `transit_gateway_id`: This column holds the ID of the transit gateway. It can be useful for joining with other tables that contain information about transit gateways.
- `vpc_id`: This column contains the ID of the VPC that is attached to the transit gateway. It can be used to join with other tables that contain information about VPCs.

## Examples

### Basic transit gateway vpc attachment info

```sql
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

```sql
select
  resource_type,
  count(transit_gateway_attachment_id) as count
from
  aws_ec2_transit_gateway_vpc_attachment
group by
  resource_type;
```