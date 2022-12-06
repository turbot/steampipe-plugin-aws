# Table: aws_pipes_pipe

Amazon EventBridge Pipes connects sources to targets. It reduces the need for specialized knowledge and integration code when developing event driven architectures, fostering consistency across your company’s applications. To set up a pipe, you choose the source, add optional filtering, define optional enrichment, and choose the target for the event data.

## Examples

### Basic info

```sql
select
  name,
  arn,
  current_state,
  creation_time,
  role_arn
from
  aws_pipes_pipe;
```

### List pipes that are not in desired state

```sql
select
  name,
  arn,
  description,
  creation_time,
  current_state,
  desired_state
from
  aws_pipes_pipe
where
  desired_state <> current_state;
```

### Get the target parameters information for each pipe

```sql
select
  name,
  target_parameters ->> 'BatchJobParameters' as batch_job_parameters,
  target_parameters ->> 'CloudWatchLogsParameters' as cloud_watch_logs_parameters,
  target_parameters ->> 'EcsTaskParameters' as ecs_task_parameters,
  target_parameters ->> 'EventBridgeEventBusParameters' as event_bridge_event_bus_parameters,
  target_parameters ->> 'HttpParameters' as http_parameters,
  target_parameters ->> 'InputTemplate' as input_template,
  target_parameters ->> 'KinesisStreamParameters' as kinesis_stream_parameters,
  target_parameters ->> 'LambdaFunctionParameters' as lambda_function_parameters,
  target_parameters ->> 'RedshiftDataParameters' as redshift_data_parameters,
  target_parameters ->> 'SageMakerPipelineParameters' as sage_maker_pipeline_parameters,
  target_parameters ->> 'SqsQueueParameters' as sqs_queue_parameters,
  target_parameters ->> 'StepFunctionStateMachineParameters' as step_function_state_machine_parameters
from
  aws_pipes_pipe;
```

## Get enrichment parameters details for each pipe

```sql
select
  name,
  enrichment_parameters ->> 'HttpParameters' as http_parameters,
  enrichment_parameters ->> 'InputTemplate' as input_template
from
  aws_pipes_pipe;
```

## List pipes created in the last 30 days

```sql
select
  name,
  creation_time,
  current_state,
  desired_state,
  enrichment,
  target
from
  aws_pipes_pipe
where
  creation_time <= now() - interval '30' day;
```

## Get role details for pipes

```sql
select
  p.name,
  r.arn as role_arn
  r.role_id,
  r.permissions_boundary_arn,
  r.role_last_used_region,
  r.inline_policies,
  r.assume_role_policy
from
  aws_pipes_pipe as p,
  aws_iam_role as r
where
  p.role_arn = r.arn;
```