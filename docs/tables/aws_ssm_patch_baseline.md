# Table: aws_ssm_patch_baseline

A patch baseline defines which patches are approved for installation on your instances.

## Examples

### List all the patch baselines in your account

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
	operating_system = 'UBUNTU';
```


### List all the patch baselines that have rejected patches

```sql
select
	baseline_id,
	name,
	description,
	operating_system,
	created_date,
	rejected_patches,
	region
from
	aws_ssm_patch_baseline
where
	rejected_patches != '[]';
```


### Approval rules details for each patch baseline

```sql
select
	baseline_id,
	p ->> 'ApproveAfterDays' as approve_after_days,
	p ->> 'ApproveUntilDate' as approve_until_date,
	p ->> 'ComplianceLevel' as compliance_level,
	p -> 'PatchFilterGroup' ->> 'PatchFilters' as patch_filters
from
	aws_ssm_patch_baseline,
	jsonb_array_elements(approval_rules -> 'PatchRules') as p;
```