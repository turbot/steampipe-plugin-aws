# Table: aws_dlm_lifecycle_policy

You can use Amazon Data Lifecycle Manager to automate the creation, retention, and deletion of EBS snapshots and EBS-backed AMIs.

## Examples

### Basic Info

```sql
select
  policy_id,
  arn,
  date_created
from
  aws_dlm_lifecycle_policy;
```

### List policies where snapshot sharing is scheduled

```sql
select
  policy_id,
  arn,
  date_created,
  policy_type,
  s ->> 'ShareRules' as share_rules
from
  aws_dlm_lifecycle_policy,
  jsonb_array_elements(policy_details -> 'Schedules') s
where 
  s ->> 'ShareRules' is not null;
```

### List policies where cross-region copying is scheduled

```sql
select
  policy_id,
  arn,
  date_created,
  policy_type,
  s ->> 'CrossRegionCopyRules' as cross_region_copy_rules
from
  aws_dlm_lifecycle_policy,
  jsonb_array_elements(policy_details -> 'Schedules') s
where 
  s ->> 'CrossRegionCopyRules' is not null;
  ```

### List maximum snapshots allowed to be retained after each schedule

```sql
select
  policy_id,
  arn,
  date_created,
  policy_type,
  s -> 'RetainRule' ->> 'Count' as retain_count
from
  aws_dlm_lifecycle_policy,
  jsonb_array_elements(policy_details -> 'Schedules') s
where 
  s -> 'RetainRule' is not null;
  ```