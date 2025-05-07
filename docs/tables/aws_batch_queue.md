---
title: "Steampipe Table: aws_batch_queue - Query AWS Batch Job Queues using SQL"
description: "Allows users to query AWS Batch Job Queues for detailed information about queue configurations, statuses, and related attributes."
folder: "Batch"
---

# Table: aws_batch_queue

AWS Batch Job Queues are resources in the AWS Batch service that store jobs until the AWS Batch Scheduler runs them on a compute environment. Job queues have a priority that the scheduler uses to determine which jobs in which queue should be evaluated for execution first.

## Table Usage Guide

The `aws_batch_queue` table provides insights into job queue configurations within AWS Batch. As a DevOps engineer, explore queue-specific details through this table, including priority settings, compute environment associations, and queue statuses. Utilize it to monitor queue configurations, verify compute environment associations, and ensure appropriate job queue priorities are set.

## Examples

### Basic info
Get a quick overview of all AWS Batch job queues, including their name, state, status, and priority. This is useful for inventory and monitoring purposes.

```sql+postgres
select
  job_queue_name,
  state,
  status,
  priority
from
  aws_batch_queue;
```

```sql+sqlite
select
  job_queue_name,
  state,
  status,
  priority
from
  aws_batch_queue;
```

### List enabled job queues
Identify all job queues that are currently enabled. This helps you focus on active queues that can accept jobs.

```sql+postgres
select
  job_queue_name,
  state,
  priority
from
  aws_batch_queue
where
  state = 'ENABLED';
```

```sql+sqlite
select
  job_queue_name,
  state,
  priority
from
  aws_batch_queue
where
  state = 'ENABLED';
```

### List job queues by state
See how many job queues exist in each state (ENABLED, DISABLED, etc.). This is useful for capacity planning and operational audits.

```sql+postgres
select
  state,
  count(*) as queue_count
from
  aws_batch_queue
group by
  state
order by
  queue_count desc;
```

```sql+sqlite
select
  state,
  count(*) as queue_count
from
  aws_batch_queue
group by
  state
order by
  queue_count desc;
```

### Get compute environments associated with each job queue
List the compute environments attached to each job queue and their order of preference. This helps you understand job placement and compute resource allocation.

```sql+postgres
select
  job_queue_name,
  state,
  c->>'Order' as order_priority,
  c->>'ComputeEnvironment' as compute_environment
from
  aws_batch_queue,
  jsonb_array_elements(compute_environment_order) as c
order by
  job_queue_name,
  (c->>'Order')::int;
```

```sql+sqlite
select
  job_queue_name,
  state,
  json_extract(c.value, '$.Order') as order_priority,
  json_extract(c.value, '$.ComputeEnvironment') as compute_environment
from
  aws_batch_queue,
  json_each(compute_environment_order) as c
order by
  job_queue_name,
  CAST(json_extract(c.value, '$.Order') as integer);
```

### Find job queues with high priority (lower number means higher priority)
Identify job queues with a high priority (lower number). This is useful for understanding which queues are preferred for job scheduling.

```sql+postgres
select
  job_queue_name,
  state,
  priority
from
  aws_batch_queue
where
  priority < 50
order by
  priority;
```

```sql+sqlite
select
  job_queue_name,
  state,
  priority
from
  aws_batch_queue
where
  priority < 50
order by
  priority;
```

### Find job queues with scheduling policies
Find job queues that have a scheduling policy attached. This helps you identify queues with custom scheduling behavior.

```sql+postgres
select
  job_queue_name,
  state,
  scheduling_policy_arn
from
  aws_batch_queue
where
  scheduling_policy_arn is not null;
```

```sql+sqlite
select
  job_queue_name,
  state,
  scheduling_policy_arn
from
  aws_batch_queue
where
  scheduling_policy_arn is not null;
```

### List job queues with their tags
Show all job queues that have tags assigned. This is useful for cost allocation, resource organization, and compliance.

```sql+postgres
select
  job_queue_name,
  state,
  tags
from
  aws_batch_queue
where
  tags is not null;
```

```sql+sqlite
select
  job_queue_name,
  state,
  tags
from
  aws_batch_queue
where
  tags is not null;
```
