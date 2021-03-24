# Table: aws_efs_access_point

Amazon EFS access points are application-specific entry points into an EFS file system that make it easier to manage application access to shared datasets.

## Examples

### Basic info

```sql
select
  name,
  access_point_id,
  access_point_arn,
  file_system_id,
  life_cycle_state,
  owner_id,
  root_directory
from
  aws_efs_access_point;
```


### List access points for a specific file system

```sql
select
  name,
  access_point_id,
  file_system_id,
  owner_id,
  root_directory
from
  aws_efs_access_point
where
  file_system_id = 'fs-82c7d9fa';
```


### List access points that are in available life cycle state

```sql
select
  name,
  access_point_id,
  life_cycle_state,
  file_system_id,
  owner_id,
  root_directory
from
  aws_efs_access_point
where
  life_cycle_state = 'available';
```