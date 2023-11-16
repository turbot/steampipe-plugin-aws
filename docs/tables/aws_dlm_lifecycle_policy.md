---
title: "Table: aws_dlm_lifecycle_policy - Query AWS DLM Lifecycle Policies using SQL"
description: "Allows users to query AWS DLM Lifecycle Policies to retrieve detailed information about each policy, including its configuration, status, and tags."
---

# Table: aws_dlm_lifecycle_policy - Query AWS DLM Lifecycle Policies using SQL

The `aws_dlm_lifecycle_policy` table in Steampipe provides information about DLM (Data Lifecycle Manager) lifecycle policies within AWS. This table allows DevOps engineers to query policy-specific details, including policy ID, policy description, state, status message, and execution details. Users can utilize this table to gather insights on policies, such as the policy execution frequency, target tags, and retention rules. The schema outlines the various attributes of the DLM lifecycle policy, including policy ARN, creation date, policy details, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dlm_lifecycle_policy` table, you can use the `.inspect aws_dlm_lifecycle_policy` command in Steampipe.

### Key columns:

- `policy_id`: The unique identifier for the lifecycle policy. This key column can be used to join this table with other tables to retrieve specific policy details.
- `policy_arn`: The Amazon Resource Name (ARN) of the lifecycle policy. This key column can be used to join this table with other tables to get detailed information about the policy's resource.
- `state`: The state of the lifecycle policy (ENABLED | DISABLED). This key column can be used to filter policies based on their current state.

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