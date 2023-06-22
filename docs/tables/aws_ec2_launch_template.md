# Table: aws_ec2_launch_template

A launch template is similar to a launch configuration, in that it specifies instance configuration information.

## Examples

### Basic info

```sql
select
  launch_template_name,
  launch_template_id,
  created_time,
  created_by,
  default_version_number,
  latest_version_number
from
  aws_ec2_launch_template;
```

### List launch templates created by a user

```sql
select
  launch_template_name,
  launch_template_id,
  create_time,
  created_by
from
  aws_ec2_launch_template
where
  created_by like '%turbot';
```

### List launch templates created in the last 30 days

```sql
select
  launch_template_name,
  launch_template_id,
  create_time
from
  aws_ec2_launch_template
where
  create_time >= now() - interval '30' day;
```