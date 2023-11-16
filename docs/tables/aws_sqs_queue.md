---
title: "Table: aws_sqs_queue - Query AWS Simple Queue Service (SQS) using SQL"
description: "Allows users to query AWS Simple Queue Service (SQS) to retrieve detailed information about each queue."
---

# Table: aws_sqs_queue - Query AWS Simple Queue Service (SQS) using SQL

The `aws_sqs_queue` table in Steampipe provides information about each queue in AWS Simple Queue Service (SQS). This table allows DevOps engineers to query queue-specific details, including ARN, URL, and associated metadata. Users can utilize this table to gather insights on queues, such as their visibility timeout, message retention period, and delivery delay settings. The schema outlines the various attributes of the SQS queue, including the queue ARN, URL, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sqs_queue` table, you can use the `.inspect aws_sqs_queue` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the queue. This can be used to join this table with other tables that contain AWS resource ARNs.
- `url`: The URL of the queue. This is a unique identifier for the queue and can be used to join with any other table that references AWS SQS queues by URL.
- `name`: The name of the queue. This is a human-readable identifier for the queue and can be useful for joining with tables that reference queues by name.

## Examples

### Basic info

```sql
select
  title,
  delay_seconds,
  max_message_size,
  receive_wait_time_seconds,
  message_retention_seconds,
  visibility_timeout_seconds
from
  aws_sqs_queue;
```

### List unencrypted queues

```sql
select
  title,
  kms_master_key_id,
  sqs_managed_sse_enabled
from
  aws_sqs_queue
where
  kms_master_key_id is null
  and not sqs_managed_sse_enabled;
```

### List queues encrypted with a CMK

```sql
select
  title,
  kms_master_key_id,
  sqs_managed_sse_enabled
from
  aws_sqs_queue
where
  kms_master_key_id is not null;
```

### List queues encrypted with an SQS-owned encryption key

```sql
select
  title,
  kms_master_key_id,
  sqs_managed_sse_enabled
from
  aws_sqs_queue
where
  sqs_managed_sse_enabled;
```

### List queues with a message retention period less than 7 days

```sql
select
  title,
  message_retention_seconds
from
  aws_sqs_queue
where
  message_retention_seconds < '604800';
```

### List queues which are not configured with a dead-letter queue (DLQ)

```sql
select
  title,
  redrive_policy
from
  aws_sqs_queue
where
  redrive_policy is null;
```

### List FIFO queues

```sql
select
  title,
  fifo_queue
from
  aws_sqs_queue
where
  fifo_queue;
```

### List queues with policy statements that grant cross-account access

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_sqs_queue,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  string_to_array(p, ':') as pa,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and (
    pa[5] != account_id
    or p = '*'
  );
```

### List queues with policy statements that grant anoymous access

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_sqs_queue,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  p = '*'
  and s ->> 'Effect' = 'Allow';
```

### List queues with policy statements that grant full access (sqs:*)

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_sqs_queue,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and a in ('*', 'sqs:*');
```
