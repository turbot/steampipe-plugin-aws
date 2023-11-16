---
title: "Table: aws_cloudformation_stack_set - Query AWS CloudFormation StackSets using SQL"
description: "Allows users to query AWS CloudFormation StackSets, providing detailed information about each StackSet's configuration, status, and associated AWS resources."
---

# Table: aws_cloudformation_stack_set - Query AWS CloudFormation StackSets using SQL

The `aws_cloudformation_stack_set` table in Steampipe provides information about StackSets within AWS CloudFormation. This table allows DevOps engineers to query StackSet-specific details, including its configuration, status, and AWS resources associated with it. Users can utilize this table to gather insights on StackSets, such as StackSets with specific configurations, their current status, and more. The schema outlines the various attributes of the StackSet, including the StackSet ID, description, status, template body, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudformation_stack_set` table, you can use the `.inspect aws_cloudformation_stack_set` command in Steampipe.

### Key columns:

- `stack_set_name`: The name of the AWS CloudFormation StackSet. This can be used to join this table with other tables that also contain StackSet names.
- `stack_set_id`: The unique identifier for the StackSet. This is useful for joining with other tables that may reference the StackSet by its ID.
- `status`: The current status of the StackSet. This can be useful for filtering StackSets based on their status or for joining with other tables that track resource statuses.

## Examples

### Basic info

```sql
select
  stack_set_id,
  stack_set_name,
  status,
  arn,
  description
from
  aws_cloudformation_stack_set;
```

### List active stack sets

```sql
select
  stack_set_id,
  stack_set_name,
  status,
  permission_model,
  auto_deployment
from
  aws_cloudformation_stack_set
where
  status = 'ACTIVE';
```

### Get parameter details of stack sets

```sql
select
  stack_set_name,
  stack_set_id,
  p ->> 'ParameterKey' as parameter_key,
  p ->> 'ParameterValue' as parameter_value,
  p ->> 'ResolvedValue' as resolved_value,
  p ->> 'UsePreviousValue' as use_previous_value
from
  aws_cloudformation_stack_set,
  jsonb_array_elements(parameters) as p;
```

### Get drift detection details of stack sets

```sql
select
  stack_set_name,
  stack_set_id,
  stack_set_drift_detection_details ->> 'DriftDetectionStatus' as drift_detection_status,
  stack_set_drift_detection_details ->> 'DriftStatus' as drift_status,
  stack_set_drift_detection_details ->> 'DriftedStackInstancesCount' as drifted_stack_instances_count,
  stack_set_drift_detection_details ->> 'FailedStackInstancesCount' as failed_stack_instances_count,
  stack_set_drift_detection_details ->> 'InProgressStackInstancesCount' as in_progress_stack_instances_count,
  stack_set_drift_detection_details ->> 'InSyncStackInstancesCount' as in_sync_stack_instances_count,
  stack_set_drift_detection_details ->> 'LastDriftCheckTimestamp' as last_drift_check_timestamp,
  stack_set_drift_detection_details ->> 'TotalStackInstancesCount' as total_stack_instances_count
from
  aws_cloudformation_stack_set;
```