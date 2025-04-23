---
title: "Steampipe Table: aws_ssm_managed_instance_patch_state - Query AWS Systems Manager Managed Instance Patch State using SQL"
description: "Allows users to query AWS Systems Manager Managed Instance Patch State to gather information about the patch state of managed instances. This includes the instance ID, patch group, owner information, installed patches, and more."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_managed_instance_patch_state - Query AWS Systems Manager Managed Instance Patch State using SQL

The AWS Systems Manager Managed Instance Patch State is a feature of AWS Systems Manager that provides information about the patch state of your managed instances. It allows you to determine the patch compliance of instances in your managed environment, helping you to maintain the security and compliance of your instances. This feature can be queried using SQL, providing detailed information about the patch status of each instance.

## Table Usage Guide

The `aws_ssm_managed_instance_patch_state` table in Steampipe provides you with information about the patch state of managed instances within AWS Systems Manager (SSM). This table allows you, as a DevOps engineer, to query specific details related to the patch state, including the instance ID, patch group, owner information, installed patches, and more. You can utilize this table to gather insights on patch compliance and to monitor the patching status of your managed instances. The schema outlines the various attributes of the managed instance patch state for you, including the instance ID, patch group, owner information, installed patches, and associated tags.

## Examples

### Basic info
Analyze the status of patch installation in AWS managed instances to understand the effectiveness of patching operations. This helps in identifying instances where patch installation has failed, thereby enabling timely troubleshooting and ensuring system security.

```sql+postgres
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

```sql+sqlite
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

### Count the number of patches installed from patch base line
Determine the total number of installed patches from a baseline to assess the level of system updates in your AWS managed instances. This can help in identifying systems that may be lagging in updates, aiding in maintaining security and performance standards.

```sql+postgres
select
  instance_id,
  baseline_id,
  installed_count
from
  aws_ssm_managed_instance_patch_state;
```

```sql+sqlite
select
  instance_id,
  baseline_id,
  installed_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count the number of patches installed not from patch base line
Determine the areas in which patches have been installed outside of the baseline, allowing for a better understanding of potential security vulnerabilities or inconsistencies in system management.

```sql+postgres
select
  instance_id,
  baseline_id,
  installed_other_count
from
  aws_ssm_managed_instance_patch_state;
```

```sql+sqlite
select
  instance_id,
  baseline_id,
  installed_other_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count of non-compliant security patches for each instance
Determine the areas in which non-compliant security patches exist for each instance. This helps in identifying potential security vulnerabilities and aids in maintaining system integrity.

```sql+postgres
select
  instance_id,
  baseline_id,
  security_non_compliant_count
from
  aws_ssm_managed_instance_patch_state;
```

```sql+sqlite
select
  instance_id,
  baseline_id,
  security_non_compliant_count
from
  aws_ssm_managed_instance_patch_state;
```

### List patch operations in the last 10 days
Explore the recent activities of patch operations within the last 10 days. This can be beneficial for monitoring and maintaining the health and security of your managed instances.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  baseline_id,
  operation,
  operation_end_time,
  operation_start_time
from
  aws_ssm_managed_instance_patch_state
where
  operation_end_time >= datetime('now', '-10 day');
```

### List scan patches
Discover the segments that are currently in the 'Scan' operation state within your managed instances. This can be particularly useful in understanding and managing your system's security patching process.

```sql+postgres
select
  instance_id,
  baseline_id,
  operation
from
  aws_ssm_managed_instance_patch_state
where
  operation = 'Scan';
```

```sql+sqlite
select
  instance_id,
  baseline_id,
  operation
from
  aws_ssm_managed_instance_patch_state
where
  operation = 'Scan';
```