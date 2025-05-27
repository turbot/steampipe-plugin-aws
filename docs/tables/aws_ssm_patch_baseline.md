---
title: "Steampipe Table: aws_ssm_patch_baseline - Query AWS SSM Patch Baseline using SQL"
description: "Allows users to query AWS SSM Patch Baseline data to retrieve information about each patch baseline in your AWS account."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_patch_baseline - Query AWS SSM Patch Baseline using SQL

The AWS Systems Manager Patch Manager (SSM Patch Manager) is a service that allows you to automate the process of patching managed instances. A patch baseline defines which patches are approved for installation on your instances. You can specify approved or rejected patches one by one or by using patch filters based on product, classification, or severity.

## Table Usage Guide

The `aws_ssm_patch_baseline` table in Steampipe allows you to query information about each patch baseline in your AWS account. This table provides you, as a DevOps engineer or system administrator, with patch-specific details, including the patch baseline ID, name, operating system, approval rules, and more. You can utilize this table to gather insights on patch baselines, such as the approved and rejected patches, patch compliance levels, and the patch groups the baseline is associated with. The schema outlines the various attributes of the AWS SSM Patch Baseline, including the baseline ID, creation date, description, and associated tags for you.

## Examples

### Basic info
Analyze the settings to understand the basic information about the patch baselines in your AWS Simple Systems Manager. This helps in gaining insights into the operating system, creation date, and geographical region of these patch baselines, aiding in better management and maintenance of your system's security.

```sql+postgres
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

```sql+sqlite
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
Gain insights into the patch baselines that are specific to the Ubuntu operating system. This is useful for maintaining system security and ensuring you're aware of all available updates.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have patch baselines with rejected patches in AWS SSM. This is useful in identifying potential issues in your patch management process and addressing them promptly.

```sql+postgres
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

```sql+sqlite
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
  json_array_length(rejected_patches) != 0;
```

### Get approval rules details for each patch baseline
Determine the specific approval rules for each patch baseline in your AWS Simple Systems Manager. This helps in understanding the timeline and compliance level for each patch, aiding in efficient system management and ensuring security compliance.

```sql+postgres
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

```sql+sqlite
select
  baseline_id,
  json_extract(p.value, '$.ApproveAfterDays') as approve_after_days,
  json_extract(p.value, '$.ApproveUntilDate') as approve_until_date,
  json_extract(p.value, '$.ComplianceLevel') as compliance_level,
  json_extract(json_extract(p.value, '$.PatchFilterGroup'), '$.PatchFilters') as patch_filters
from
  aws_ssm_patch_baseline,
  json_each(approval_rules, '$.PatchRules') as p;
```