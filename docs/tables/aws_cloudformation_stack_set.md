---
title: "Steampipe Table: aws_cloudformation_stack_set - Query AWS CloudFormation StackSets using SQL"
description: "Allows users to query AWS CloudFormation StackSets, providing detailed information about each StackSet's configuration, status, and associated AWS resources."
folder: "CloudFormation"
---

# Table: aws_cloudformation_stack_set - Query AWS CloudFormation StackSets using SQL

The AWS CloudFormation StackSets is a feature within the AWS CloudFormation service that allows you to create, update, or delete stacks across multiple accounts and regions with a single AWS CloudFormation template. StackSets takes care of the underlying details of orchestrating stack operations across multiple accounts and regions, ensuring that the stacks are created, updated, or deleted in a specified order. This simplifies the management of AWS resources and enables the easy deployment of regional and global applications.

## Table Usage Guide

The `aws_cloudformation_stack_set` table in Steampipe provides you with information about StackSets within AWS CloudFormation. This table allows you, as a DevOps engineer, to query StackSet-specific details, including its configuration, status, and AWS resources associated with it. You can utilize this table to gather insights on StackSets, such as StackSets with specific configurations, their current status, and more. The schema outlines the various attributes of the StackSet for you, including the StackSet ID, description, status, template body, and associated tags.

## Examples

### Basic info
Explore which AWS CloudFormation stack sets are in use and their current status. This can be useful for auditing purposes, understanding your resource utilization, and identifying any potential issues with your stacks.

```sql+postgres
select
  stack_set_id,
  stack_set_name,
  status,
  arn,
  description
from
  aws_cloudformation_stack_set;
```

```sql+sqlite
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
Determine the areas in which active stack sets are being used within your AWS CloudFormation service. This allows you to monitor and manage your active resources effectively.

```sql+postgres
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

```sql+sqlite
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
This query allows you to delve into the specifics of your stack sets within AWS CloudFormation. It's particularly valuable for understanding the parameters associated with each stack set, which can help in managing and optimizing your cloud resources.

```sql+postgres
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

```sql+sqlite
select
  stack_set_name,
  stack_set_id,
  json_extract(p.value, '$.ParameterKey') as parameter_key,
  json_extract(p.value, '$.ParameterValue') as parameter_value,
  json_extract(p.value, '$.ResolvedValue') as resolved_value,
  json_extract(p.value, '$.UsePreviousValue') as use_previous_value
from
  aws_cloudformation_stack_set,
  json_each(parameters) as p;
```

### Get drift detection details of stack sets
Explore the drift detection status of your stack sets to identify any potential issues or discrepancies. This can help in maintaining the overall health and integrity of your stack sets.

```sql+postgres
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

```sql+sqlite
select
  stack_set_name,
  stack_set_id,
  json_extract(stack_set_drift_detection_details, '$.DriftDetectionStatus') as drift_detection_status,
  json_extract(stack_set_drift_detection_details, '$.DriftStatus') as drift_status,
  json_extract(stack_set_drift_detection_details, '$.DriftedStackInstancesCount') as drifted_stack_instances_count,
  json_extract(stack_set_drift_detection_details, '$.FailedStackInstancesCount') as failed_stack_instances_count,
  json_extract(stack_set_drift_detection_details, '$.InProgressStackInstancesCount') as in_progress_stack_instances_count,
  json_extract(stack_set_drift_detection_details, '$.InSyncStackInstancesCount') as in_sync_stack_instances_count,
  json_extract(stack_set_drift_detection_details, '$.LastDriftCheckTimestamp') as last_drift_check_timestamp,
  json_extract(stack_set_drift_detection_details, '$.TotalStackInstancesCount') as total_stack_instances_count
from
  aws_cloudformation_stack_set;
```