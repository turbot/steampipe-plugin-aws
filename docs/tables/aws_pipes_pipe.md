---
title: "Table: aws_pipes_pipe - Query AWS Pipes using SQL"
description: "Allows users to query AWS Pipes to obtain detailed information about individual pipes."
---

# Table: aws_pipes_pipe - Query AWS Pipes using SQL

The `aws_pipes_pipe` table in Steampipe provides information about individual pipes within AWS Pipes. This table allows DevOps engineers to query pipe-specific details, including the pipe name, pipe ARN, and the creation time. Users can utilize this table to gather insights on pipes, such as pipe statuses, pipe types, and more. The schema outlines the various attributes of the pipe, including the pipe name, pipe ARN, creation time, and pipe status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_pipes_pipe` table, you can use the `.inspect aws_pipes_pipe` command in Steampipe.

### Key columns:

- `name`: The name of the pipe. This can be used to join this table with other tables that also contain pipe names.
- `arn`: The Amazon Resource Number (ARN) of the pipe. This unique identifier can be used to join this table with other tables that also contain pipe ARNs.
- `status`: The status of the pipe. This can be used to filter pipes based on their status.

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
  target_parameters ->> 'CloudWatchLogsParameters' as cloudwatch_logs_parameters,
  target_parameters ->> 'EcsTaskParameters' as ecs_task_parameters,
  target_parameters ->> 'EventBridgeEventBusParameters' as eventbridge_event_bus_parameters,
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

## Get enrichment parameter details for each pipe

```sql
select
  name,
  enrichment_parameters ->> 'HttpParameters' as http_parameters,
  enrichment_parameters ->> 'InputTemplate' as input_template
from
  aws_pipes_pipe;
```

### List pipes created within the last 30 days

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
  creation_time >= now() - interval '30' day;
```

### Get IAM role details for pipes

```sql
select
  p.name,
  r.arn as role_arn,
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