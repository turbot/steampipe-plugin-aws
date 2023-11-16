---
title: "Table: aws_ssm_patch_baseline - Query AWS SSM Patch Baseline using SQL"
description: "Allows users to query AWS SSM Patch Baseline data to retrieve information about each patch baseline in your AWS account."
---

# Table: aws_ssm_patch_baseline - Query AWS SSM Patch Baseline using SQL

The `aws_ssm_patch_baseline` table in Steampipe allows you to query information about each patch baseline in your AWS account. This table provides DevOps engineers and system administrators with patch-specific details, including the patch baseline ID, name, operating system, approval rules, and more. Users can utilize this table to gather insights on patch baselines, such as the approved and rejected patches, patch compliance levels, and the patch groups the baseline is associated with. The schema outlines the various attributes of the AWS SSM Patch Baseline, including the baseline ID, creation date, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_patch_baseline` table, you can use the `.inspect aws_ssm_patch_baseline` command in Steampipe.

### Key columns:

- `baseline_id`: This is the ID of the patch baseline. It can be used to join this table with other tables that contain patch baseline information.
- `name`: This is the name of the patch baseline. It provides a human-readable identifier that can be useful in queries and reports.
- `operating_system`: This indicates the operating system of the patch baseline. It can be used to filter or join data based on the operating system.

## Examples

### Basic info

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

### List patch baselines for a specific operating system

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

### List patch baselines that have rejected patches

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

### Get approval rules details for each patch baseline

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
