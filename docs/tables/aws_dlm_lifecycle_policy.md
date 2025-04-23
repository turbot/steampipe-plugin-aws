---
title: "Steampipe Table: aws_dlm_lifecycle_policy - Query AWS DLM Lifecycle Policies using SQL"
description: "Allows users to query AWS DLM Lifecycle Policies to retrieve detailed information about each policy, including its configuration, status, and tags."
folder: "DLM"
---

# Table: aws_dlm_lifecycle_policy - Query AWS DLM Lifecycle Policies using SQL

The AWS DLM (Data Lifecycle Manager) Lifecycle Policy is a service that automates the creation, retention, and deletion of Amazon EBS volume snapshots. This service eliminates the need for custom scripts and manual operations to manage the lifecycle of EBS volume snapshots. It allows you to manage the lifecycle of your snapshots with policy-based management, reducing the cost and effort of data backup, disaster recovery, and migration tasks.

## Table Usage Guide

The `aws_dlm_lifecycle_policy` table in Steampipe provides you with information about DLM (Data Lifecycle Manager) lifecycle policies within AWS. This table enables you, as a DevOps engineer, to query policy-specific details, including policy ID, policy description, state, status message, and execution details. You can utilize this table to gather insights on policies, such as the policy execution frequency, target tags, and retention rules. The schema outlines the various attributes of the DLM lifecycle policy for you, including policy ARN, creation date, policy details, and associated tags.

## Examples

### Basic Info
Explore which AWS Data Lifecycle Manager policies have been created and when, to manage and monitor the lifecycle of your AWS resources effectively.

```sql+postgres
select
  policy_id,
  arn,
  date_created
from
  aws_dlm_lifecycle_policy;
```

```sql+sqlite
select
  policy_id,
  arn,
  date_created
from
  aws_dlm_lifecycle_policy;
```

### List policies where snapshot sharing is scheduled
Determine the areas in which snapshot sharing is scheduled within your policy settings. This helps to identify potential security risks and ensure data integrity.

```sql+postgres
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

```sql+sqlite
select
  policy_id,
  arn,
  date_created,
  policy_type,
  json_extract(s.value, '$.ShareRules') as share_rules
from
  aws_dlm_lifecycle_policy,
  json_each(json_extract(policy_details, '$.Schedules')) as s
where 
  json_extract(s.value, '$.ShareRules') is not null;
```

### List policies where cross-region copying is scheduled
Explore policies that have cross-region copying scheduled. This is useful to identify and manage data replication across different geographical areas for redundancy and disaster recovery purposes.

```sql+postgres
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

```sql+sqlite
select
  policy_id,
  arn,
  date_created,
  policy_type,
  json_extract(s.value, '$.CrossRegionCopyRules') as cross_region_copy_rules
from
  aws_dlm_lifecycle_policy,
  json_each(json_extract(policy_details, '$.Schedules')) as s
where 
  json_extract(s.value, '$.CrossRegionCopyRules') is not null;
```

### List maximum snapshots allowed to be retained after each schedule
Discover the segments that have rules for cross-region copying in your AWS DLM lifecycle policies. This can be useful to manage and optimize your data lifecycle, especially if you have policies that need to retain a certain number of snapshots across different regions for backup or disaster recovery purposes.

```sql+postgres
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

```sql+sqlite
select
  policy_id,
  arn,
  date_created,
  policy_type,
  json_extract(json_extract(s.value, '$.RetainRule'), '$.Count') as retain_count
from
  aws_dlm_lifecycle_policy,
  json_each(json_extract(policy_details, '$.Schedules')) as s
where 
  json_extract(s.value, '$.RetainRule') is not null;
```