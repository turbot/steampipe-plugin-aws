# Table: aws_appstream_fleet

AWS AppStream Fleet is a group of streaming instances that deliver desktop applications to users. It is a service provided by Amazon Web Services (AWS) that allows you to stream desktop applications from the cloud to any device with an internet connection. A fleet consists of a collection of streaming instances that are provisioned and managed as a single entity. These instances run the applications and stream the user interface and input/output back and forth between the user's device and the cloud. AppStream Fleet helps simplify the management and delivery of applications, allowing users to access their applications securely from anywhere.

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