# Table: aws_cloudformation_stack_set

AWS CloudFormation Stack Set is a service that allows you to create, update, or delete stacks across multiple accounts and regions with a single CloudFormation template. It provides centralized management and deployment of CloudFormation stacks, making it easier to manage infrastructure resources across a large number of accounts and regions. Stack Sets can be used to enforce compliance policies, deploy standard infrastructure configurations, and simplify the management of infrastructure changes at scale.

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