---
title: "Steampipe Table: aws_iot_thing_group - Query AWS IoT Thing Group using SQL"
description: "Allows users to query AWS IoT Thing Group to gain insights into each group's configuration, including ARN, creation date, version of the group, and parent of the groups."
folder: "IoT Core"
---

# Table: aws_iot_thing_group - Query AWS IoT Thing Groups using SQL

In AWS IoT Core, an IoT Thing Group is a logical collection of IoT devices (Things) that enables collective management and interaction. These groups are integral to AWS IoT Core's device management capabilities, allowing for organized, efficient control and management of device fleets. This is particularly advantageous in large-scale IoT deployments where individual device management is not feasible.

## Table Usage Guide

The `aws_iot_thing_group` table can be utilized to obtain detailed information about IoT Thing Groups. This includes their names, IDs, descriptions, and hierarchical structures. This table is essential for IoT administrators and developers for effective organization and oversight of groups of IoT devices within AWS.

## Examples

### Basic info
Acquire essential details about AWS IoT Thing Groups, including their names, IDs, descriptions, and hierarchical relationships. This fundamental query is crucial for an overview of the groups and their structures.

```sql+postgres
select
  group_name,
  thing_group_id,
  thing_group_description,
  arn,
  creation_date,
  parent_group_name
from
  aws_iot_thing_group;
```

```sql+sqlite
select
  group_name,
  thing_group_id,
  thing_group_description,
  arn,
  creation_date,
  parent_group_name
from
  aws_iot_thing_group;
```

### Filter thing groups by parent group
Identify specific IoT Thing Groups based on their parent group. This query is useful for analyzing the hierarchical organization of your IoT devices.

```sql+postgres
select
  group_name,
  thing_group_id,
  creation_date,
  parent_group_name,
  version
from
  aws_iot_thing_group
where
  parent_group_name = 'foo';
```

```sql+sqlite
select
  group_name,
  thing_group_id,
  creation_date,
  parent_group_name,
  version
from
  aws_iot_thing_group
where
  parent_group_name = 'foo';
```

### List thing groups created in the last 30 days
Discover Thing Groups that have been created in the last 30 days. This query helps in tracking recent additions to your IoT environment.

```sql+postgres
select
  group_name,
  thing_group_id,
  parent_group_name,
  creation_date,
  status
from
  aws_iot_thing_group
where
  creation_date >= now() - interval '30 days';
```

```sql+sqlite
select
  group_name,
  thing_group_id,
  parent_group_name,
  creation_date,
  status
from
  aws_iot_thing_group
where
  datetime(creation_date) >= datetime('now', '-30 days');
```

### List active thing groups
List Thing Groups that are currently being used to manage and organize devices (referred to as "things") within your IoT environment.

```sql+postgres
select
  group_name,
  thing_group_id,
  query_string,
  query_version,
  status
from
  aws_iot_thing_group
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  group_name,
  thing_group_id,
  query_string,
  query_version,
  status
from
  aws_iot_thing_group
where
  status = 'ACTIVE';
```
