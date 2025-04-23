---
title: "Steampipe Table: aws_pipes_pipe - Query AWS Pipes using SQL"
description: "Allows users to query AWS Pipes to obtain detailed information about individual pipes."
folder: "Pipes"
---

# Table: aws_pipes_pipe - Query AWS Pipes using SQL

The AWS Pipes service is a tool that allows you to interact with your AWS resources using SQL syntax. It provides a unified interface to query and manipulate your AWS data, facilitating easy integration with existing SQL-based tools and workflows. AWS Pipes essentially turns your AWS data into a relational database, offering powerful querying capabilities for data analysis and management.

## Table Usage Guide

The `aws_pipes_pipe` table in Steampipe provides you with information about individual pipes within AWS Pipes. This table allows you, as a DevOps engineer, to query pipe-specific details, including the pipe name, pipe ARN, and the creation time. You can utilize this table to gather insights on pipes, such as pipe statuses, pipe types, and more. The schema outlines the various attributes of the pipe for you, including the pipe name, pipe ARN, creation time, and pipe status.

## Examples

### Basic info
Determine the areas in which specific AWS pipes are currently active and when they were created. This information is useful for auditing purposes, helping to track the lifecycle and usage of each pipe within your AWS infrastructure.

```sql+postgres
select
  name,
  arn,
  current_state,
  creation_time,
  role_arn
from
  aws_pipes_pipe;
```

```sql+sqlite
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
Identify instances where the current state of your AWS pipes does not match the desired state. This can be useful for troubleshooting or identifying potential issues within your AWS environment.

```sql+postgres
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

```sql+sqlite
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
  desired_state != current_state;
```

### Get the target parameters information for each pipe
Explore the different types of parameters that each pipe targets, which can help in understanding the specific configurations and settings associated with each pipe. This can be useful in identifying potential bottlenecks or areas for improvement in the system's data processing capabilities. 

Determine the enrichment parameters for each pipe, providing insights into how each pipe is enhancing the data flow. This could be beneficial in optimizing the data processing and enhancing the overall system performance.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(target_parameters, '$.BatchJobParameters') as batch_job_parameters,
  json_extract(target_parameters, '$.CloudWatchLogsParameters') as cloudwatch_logs_parameters,
  json_extract(target_parameters, '$.EcsTaskParameters') as ecs_task_parameters,
  json_extract(target_parameters, '$.EventBridgeEventBusParameters') as eventbridge_event_bus_parameters,
  json_extract(target_parameters, '$.HttpParameters') as http_parameters,
  json_extract(target_parameters, '$.InputTemplate') as input_template,
  json_extract(target_parameters, '$.KinesisStreamParameters') as kinesis_stream_parameters,
  json_extract(target_parameters, '$.LambdaFunctionParameters') as lambda_function_parameters,
  json_extract(target_parameters, '$.RedshiftDataParameters') as redshift_data_parameters,
  json_extract(target_parameters, '$.SageMakerPipelineParameters') as sage_maker_pipeline_parameters,
  json_extract(target_parameters, '$.SqsQueueParameters') as sqs_queue_parameters
  json_extract(target_parameters, '$.StepFunctionStateMachineParameters') as step_function_state_machine_parameters
from
  aws_pipes_pipe;
```

## Get enrichment parameter details for each pipe

```sql+postgres
select
  name,
  enrichment_parameters ->> 'HttpParameters' as http_parameters,
  enrichment_parameters ->> 'InputTemplate' as input_template
from
  aws_pipes_pipe;
```

```sql+sqlite
select
  name,
  json_extract(enrichment_parameters, '$.HttpParameters') as http_parameters,
  json_extract(enrichment_parameters, '$.InputTemplate') as input_template
from
  aws_pipes_pipe;
```

### List pipes created within the last 30 days
Identify instances where new pipelines have been established within the past month. This can be useful for tracking recent changes or additions to your AWS resources.

```sql+postgres
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

```sql+sqlite
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
  creation_time >= datetime('now', '-30 day');
```

### Get IAM role details for pipes
This query is useful for identifying the specific AWS IAM roles associated with your Steampipe pipelines. It provides insights into the permissions, boundaries, and usage of these roles, helping you manage security and access controls more effectively.

```sql+postgres
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

```sql+sqlite
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