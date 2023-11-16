---
title: "Table: aws_oam_link - Query AWS OAM Links using SQL"
description: "Allows users to query AWS OAM Links to gather information about the link between an AWS resource and an AWS OAM resource."
---

# Table: aws_oam_link - Query AWS OAM Links using SQL

The `aws_oam_link` table in Steampipe provides information about the links between an AWS resource and an AWS OAM (Operations Account Management) resource. This table allows DevOps engineers to query link-specific details, including the link status, link type, and associated metadata. Users can utilize this table to gather insights on links, such as their current status, the type of AWS resource linked, the type of OAM resource linked, and more. The schema outlines the various attributes of the OAM link, including the link ID, creation date, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_oam_link` table, you can use the `.inspect aws_oam_link` command in Steampipe.

### Key columns:

- `link_id`: This is the unique identifier for the link. It can be used to join this table with other tables that contain link-specific information.
- `resource_arn`: This is the Amazon Resource Name (ARN) of the AWS resource that is linked. It can be used to join with other tables that contain resource-specific information.
- `oam_resource_arn`: This is the Amazon Resource Name (ARN) of the OAM resource that is linked. It can be used to join with other tables that contain OAM resource-specific information.

## Examples

### Basic info

```sql
select
  id,
  arn,
  sink_arn,
  label,
  resource_types
from
  aws_oam_link;
```

### Get sink details of each link

```sql
select
  l.id,
  l.arn,
  s.name as sink_name,
  l.sink_arn
from
  aws_oam_link as l,
  aws_oam_sink as s;
```

### List links that share data of CloudWatch log group resource type

```sql
select
  id,
  arn,
  label,
  label_template,
  r as resource_type
from
  aws_oam_link,
  jsonb_array_elements_text(resource_types) as r
where
  r = 'AWS::Logs::LogGroup';
```