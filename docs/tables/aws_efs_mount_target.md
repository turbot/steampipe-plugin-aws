# Table: aws_efs_mount_target

A mount target is an NFS endpoint that lives in a VCN subnet of your choice and provides network access for file systems.

## Examples

### Basic info

```sql
select
  mount_target_id,
  file_system_id,
  life_cycle_state,
  availability_zone_id,
  availability_zone_name
from
  aws_efs_mount_target;
```

### Get network details for each mount target

```sql
select
  mount_target_id,
  network_interface_id,
  subnet_id,
  vpc_id
from
  aws_efs_mount_target;
```
