# Table: aws_ssm_patch_baseline

A patch baseline defines which patches are approved for installation on your instances. It is possible to specify approved or rejected patches one by one.

## Examples

### List all the patch baseline in your account

```sql
select
  baseline_id,
  name,
  description,
  operating_system,
  created_date,
  region
from
  aws_ssm_patch_baseline;
```


### List all the patch baselines for a specific type of operating system

```sql
select
  baseline_id,
  name,
  description,
  created_date,
  region
from
  aws_ssm_patch_baseline
where
  operating_system='UBUNTU';
```

### List all the patch baselines which have no rejected patches

```sql
select
  baseline_id,
  name,
  description,
  operating_system,
  created_date,
  region
from
  aws_ssm_patch_baseline
where
  rejected_patches='[]';
```

### List the ApproveAfterDays, ApproveUntilDate & PatchFilters fields from approval rules for all the patch baselines

```sql
select
  baseline_id,
  p ->> 'ApproveAfterDays' as approve_after_days,
  p ->> 'ApproveUntilDate' as approve_until_date,
  p -> 'PatchFilterGroup' ->> 'PatchFilters' as patch_filters
from
  aws_ssm_patch_baseline
  cross join jsonb_array_elements(approval_rules -> 'PatchRules') as p;
```