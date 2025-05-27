---
title: "Steampipe Table: aws_oam_link - Query AWS OAM Links using SQL"
description: "Allows users to query AWS OAM Links to gather information about the link between an AWS resource and an AWS OAM resource."
folder: "OAM"
---

# Table: aws_oam_link - Query AWS OAM Links using SQL

The AWS OAM (Operations Account Management) Link is a service that allows you to manage and organize your AWS accounts. It provides a simple way to create and manage links between your AWS accounts and your AWS Organizations. This makes it easier to manage your accounts, reduce operational overhead, and improve security compliance across your organization.

## Table Usage Guide

The `aws_oam_link` table in Steampipe provides you with information about the links between an AWS resource and an AWS OAM (Operations Account Management) resource. This table allows you, as a DevOps engineer, to query link-specific details, including the link status, link type, and associated metadata. You can utilize this table to gather insights on links, such as their current status, the type of AWS resource linked, the type of OAM resource linked, and more. The schema outlines for you the various attributes of the OAM link, including the link ID, creation date, status, and associated tags.

## Examples

### Basic info
Explore which AWS resources are linked in your environment to understand the connections and dependencies between them. This could be useful for instance, in managing and optimizing your cloud resources or troubleshooting issues.

```sql+postgres
select
  id,
  arn,
  sink_arn,
  label,
  resource_types
from
  aws_oam_link;
```

```sql+sqlite
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
Explore the connections between different components in your network to understand the flow of data. This can help in identifying bottlenecks or potential points of failure, enhancing overall system efficiency and reliability.

```sql+postgres
select
  l.id,
  l.arn,
  s.name as sink_name,
  l.sink_arn
from
  aws_oam_link as l,
  aws_oam_sink as s;
```

```sql+sqlite
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
Determine the areas in which links are sharing data of a specific CloudWatch log group resource type. This can be useful for assessing the elements within your AWS environment that are interacting with these log groups, aiding in both security and resource management.

```sql+postgres
select
  l.id,
  l.arn,
  l.label,
  l.label_template,
  r as resource_type
from
  aws_oam_link as l,
  jsonb_array_elements_text(resource_types) as r
where
  r = 'AWS::Logs::LogGroup';
```

```sql+sqlite
select
  l.id,
  l.arn,
  l.label,
  l.label_template,
  json_extract(r.value, '$') as resource_type
from
  aws_oam_link as l,
  json_each(resource_types) as r
where
  json_extract(r.value, '$') = 'AWS::Logs::LogGroup';
```