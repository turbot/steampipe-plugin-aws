# Table: aws_ssm_managed_instance_patch_state

AWS Systems Manager (SSM) provides the capability to manage and patch instances in your AWS environment. The managed instance patch states in AWS SSM refer to the different states that an instance can be in during the patching process.

## Examples

### Basic info

```sql
select
  instance_id,
  baseline_id,
  operation,
  patch_group,
  failed_count,
  installed_count,
  installed_other_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count number of patches installed form patch base line

```sql
select
  instance_id,
  baseline_id,
  installed_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count number of patches installed not form patch base line

```sql
select
  instance_id,
  baseline_id,
  installed_other_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count security non compliant per node

```sql
select
  instance_id,
  baseline_id,
  security_non_compliant_count
from
  aws_ssm_managed_instance_patch_state;
```

### List patch operations in the last 10 days

```sql
select
  instance_id,
  baseline_id,
  operation,
  operation_end_time,
  operation_start_time
from
  aws_ssm_managed_instance_patch_state
where
  operation_end_time >= now() - interval '10' day;
```

### List scan patches

```sql
select
  instance_id,
  baseline_id,
  operation
from
  aws_ssm_managed_instance_patch_state
where
  operation = 'Scan';
```