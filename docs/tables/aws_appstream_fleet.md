---
title: "Table: aws_appstream_fleet - Query AWS AppStream Fleet using SQL"
description: "Allows users to query AWS AppStream Fleets for detailed information about each fleet, including its state, instance type, and associated stack details."
---

# Table: aws_appstream_fleet - Query AWS AppStream Fleet using SQL

The `aws_appstream_fleet` table in Steampipe provides information about fleets within AWS AppStream. This table allows DevOps engineers to query fleet-specific details, including the fleet state, instance type, associated stack details, and more. Users can utilize this table to gather insights on fleets, such as the fleet's current capacity, the fleet's idle disconnect timeout settings, and the fleet's stream view. The schema outlines the various attributes of the AppStream Fleet, including the fleet ARN, creation time, fleet type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_appstream_fleet` table, you can use the `.inspect aws_appstream_fleet` command in Steampipe.

* Key columns:

    - `name`: The name of the fleet. This column can be used to join with other tables that contain fleet-specific information.
    - `arn`: The Amazon Resource Name (ARN) of the fleet. This column can be used to join with other tables that contain ARN-specific information.
    - `state`: The current state of the fleet. This column can be used to join with other tables that contain state-specific information.

## Examples

### Basic info

```sql
select
  name,
  arn,
  instance_type,
  description,
  created_time,
  display_name,
  state,
  directory_name,
  enable_default_internet_access
from
  aws_appstream_fleet;
```

### List fleets that have default internet access anabled

```sql
select
  name,
  arn,
  instance_type,
  description,
  created_time,
  display_name,
  state,
  enable_default_internet_access
from
  aws_appstream_fleet
where enable_default_internet_access;
```

### List on-demand fleets

```sql
select
  name,
  created_time,
  fleet_type,
  instance_type,
  display_name,
  image_arn,
  image_name
from
  aws_appstream_fleet
where
  fleet_type = 'ON_DEMAND';
```

### List fleets that are created in last 30 days

```sql
select
  name,
  created_time,
  display_name,
  enable_default_internet_access,
  max_concurrent_sessions,
  max_user_duration_in_seconds
from
  aws_appstream_fleet
where
  created_time >= now() - interval '30' day;
```

### List fleets that are using private images

```sql
select
  f.name,
  f.created_time,
  f.display_name,
  f.image_arn,
  i.base_image_arn,
  i.image_builder_name,
  i.visibility
from
  aws_appstream_fleet as f,
  aws_appstream_image as i
where
  i.arn = f.image_arn
and
  i.visibility = 'PRIVATE';
```

### Get compute capacity status of each fleet

```sql
select
  name,
  arn,
  compute_capacity_status ->> 'Available' as available,
  compute_capacity_status ->> 'Desired' as desired,
  compute_capacity_status ->> 'InUse' as in_use,
  compute_capacity_status ->> 'Running' as running
from
  aws_appstream_fleet;
```

### Get error details of failed images

```sql
select
  name,
  arn,
  e ->> 'ErrorCode' as error_code,
  e ->> 'ErrorMessage' as error_message
from
  aws_appstream_fleet,
  jsonb_array_elements(fleet_errors) as e;
```

### Get VPC config details of each fleet

```sql
select
  name,
  arn,
  vpc_config -> 'SecurityGroupIds' as security_group_ids,
  vpc_config -> 'SubnetIds' as subnet_ids
from
  aws_appstream_fleet;
```

### Count fleets by instance type

```sql
select
  name,
  instance_type,
  Count(instance_type) as number_of_fleets
from
  aws_appstream_fleet
group by
  instance_type,
  name;
```

### List fleets that are in running state

```sql
select
  name,
  arn,
  state,
  created_time,
  description
from
  aws_appstream_fleet
where
  state = 'RUNNING';
```