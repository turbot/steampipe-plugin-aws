---
title: "Table: aws_ssm_managed_instance_patch_state - Query AWS Systems Manager Managed Instance Patch State using SQL"
description: "Allows users to query AWS Systems Manager Managed Instance Patch State to gather information about the patch state of managed instances. This includes the instance ID, patch group, owner information, installed patches, and more."
---

# Table: aws_ssm_managed_instance_patch_state - Query AWS Systems Manager Managed Instance Patch State using SQL

The `aws_ssm_managed_instance_patch_state` table in Steampipe provides information about the patch state of managed instances within AWS Systems Manager (SSM). This table allows DevOps engineers to query specific details related to the patch state, including the instance ID, patch group, owner information, installed patches, and more. Users can utilize this table to gather insights on patch compliance and to monitor the patching status of their managed instances. The schema outlines the various attributes of the managed instance patch state, including the instance ID, patch group, owner information, installed patches, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_managed_instance_patch_state` table, you can use the `.inspect aws_ssm_managed_instance_patch_state` command in Steampipe.

**Key columns**:

- `instance_id`: This is the ID of the managed instance. It is a key column as it uniquely identifies each instance and can be used to join with other tables that contain instance-specific information.
- `patch_group`: This column represents the patch group that the instance is part of. It is useful for querying instances based on their patch group.
- `owner_information`: This column contains information about the owner of the instance. It is important for tracking ownership and responsibility for each instance.

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

### Count the number of patches installed from patch base line

```sql
select
  instance_id,
  baseline_id,
  installed_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count the number of patches installed not from patch base line

```sql
select
  instance_id,
  baseline_id,
  installed_other_count
from
  aws_ssm_managed_instance_patch_state;
```

### Count of non-compliant security patches for each instance

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
