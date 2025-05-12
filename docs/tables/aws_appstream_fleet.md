---
title: "Steampipe Table: aws_appstream_fleet - Query AWS AppStream Fleet using SQL"
description: "Allows users to query AWS AppStream Fleets for detailed information about each fleet, including its state, instance type, and associated stack details."
folder: "AppStream"
---

# Table: aws_appstream_fleet - Query AWS AppStream Fleet using SQL

The AWS AppStream Fleet is a part of Amazon AppStream 2.0, a fully managed, secure application streaming service that allows you to stream desktop applications from AWS to any device running a web browser. It provides users instant-on access to the applications they need, and a responsive, fluid user experience on the device of their choice. An AppStream Fleet consists of streaming instances that run the image builder to stream applications to users.

## Table Usage Guide

The `aws_appstream_fleet` table in Steampipe provides you with information about fleets within AWS AppStream. This table allows you, as a DevOps engineer, to query fleet-specific details, including the fleet state, instance type, associated stack details, and more. You can utilize this table to gather insights on fleets, such as the fleet's current capacity, the fleet's idle disconnect timeout settings, and the fleet's stream view. The schema outlines the various attributes of the AppStream Fleet for you, including the fleet ARN, creation time, fleet type, and associated tags.

## Examples

### Basic info
Explore the characteristics of your AWS AppStream fleet, such as its creation time, state, and whether default internet access is enabled. This can help you understand the configuration and status of your fleet for better resource management.

```sql+postgres
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

```sql+sqlite
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
Determine the fleets that have their default internet access enabled. This is beneficial for assessing which fleets are potentially exposed to internet-based threats, thereby assisting in risk management and security planning.

```sql+postgres
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

```sql+sqlite
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
where enable_default_internet_access = 1;
```

### List on-demand fleets
Identify instances where on-demand fleets in AWS AppStream are being used, allowing users to understand the scope and details of their on-demand resource utilization. This information can be valuable for cost management and resource allocation strategies.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have been established within the last month to understand their internet access status, maximum concurrent sessions, and user duration limits. This can be beneficial for assessing recent changes or additions to your fleet configurations.

```sql+postgres
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

```sql+sqlite
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
  created_time >= datetime('now','-30 day');
```

### List fleets that are using private images
Explore which fleets are utilizing private images, allowing you to assess the level of privacy and security in your AWS AppStream fleets. This can be particularly useful in managing resource allocation and ensuring compliance with internal policies regarding data privacy.

```sql+postgres
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

```sql+sqlite
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
Assess the elements within each fleet in terms of compute capacity to ensure efficient resource management and optimal performance. This can help in identifying any discrepancies between desired and actual usage, thereby aiding in capacity planning and optimization.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  json_extract(compute_capacity_status, '$.Available') as available,
  json_extract(compute_capacity_status, '$.Desired') as desired,
  json_extract(compute_capacity_status, '$.InUse') as in_use,
  json_extract(compute_capacity_status, '$.Running') as running
from
  aws_appstream_fleet;
```

### Get error details of failed images
Identify instances where images have failed within the AWS AppStream fleet by analyzing the associated error codes and messages. This can assist in troubleshooting and rectifying issues promptly.

```sql+postgres
select
  name,
  arn,
  e ->> 'ErrorCode' as error_code,
  e ->> 'ErrorMessage' as error_message
from
  aws_appstream_fleet,
  jsonb_array_elements(fleet_errors) as e;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(e.value, '$.ErrorCode') as error_code,
  json_extract(e.value, '$.ErrorMessage') as error_message
from
  aws_appstream_fleet,
  json_each(fleet_errors) as e;
```

### Get VPC config details of each fleet
Analyze the settings to understand the configuration details of each fleet in your AWS Appstream service. This can help in managing network access and security for your fleets by identifying their associated security groups and subnets.

```sql+postgres
select
  name,
  arn,
  vpc_config -> 'SecurityGroupIds' as security_group_ids,
  vpc_config -> 'SubnetIds' as subnet_ids
from
  aws_appstream_fleet;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(vpc_config, '$.SecurityGroupIds') as security_group_ids,
  json_extract(vpc_config, '$.SubnetIds') as subnet_ids
from
  aws_appstream_fleet;
```

### Count fleets by instance type
Identify the variety of fleets based on their instance type within your AWS AppStream service. This can help optimize resource allocation by showing where the most and least populated instance types are.

```sql+postgres
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

```sql+sqlite
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
Explore which fleets are currently active and operational. This is useful for monitoring the status of your resources and ensuring they are functioning as expected.

```sql+postgres
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

```sql+sqlite
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