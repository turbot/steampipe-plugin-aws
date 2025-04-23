---
title: "Steampipe Table: aws_sqs_queue - Query AWS Simple Queue Service (SQS) using SQL"
description: "Allows users to query AWS Simple Queue Service (SQS) to retrieve detailed information about each queue."
folder: "SQS"
---

# Table: aws_sqs_queue - Query AWS Simple Queue Service (SQS) using SQL

The AWS Simple Queue Service (SQS) is a fully managed message queuing service that enables you to decouple and scale microservices, distributed systems, and serverless applications. SQS eliminates the complexity and overhead associated with managing and operating message oriented middleware, and empowers developers to focus on differentiating work. Using SQS, you can send, store, and receive messages between software components at any volume, without losing messages or requiring other services to be available.

## Table Usage Guide

The `aws_sqs_queue` table in Steampipe provides you with information about each queue in AWS Simple Queue Service (SQS). This table allows you, as a DevOps engineer, to query queue-specific details, including ARN, URL, and associated metadata. You can utilize this table to gather insights on queues, such as their visibility timeout, message retention period, and delivery delay settings. The schema outlines the various attributes of the SQS queue for you, including the queue ARN, URL, and associated tags.

## Examples

### Basic info
Explore the configuration of your AWS Simple Queue Service (SQS) to understand factors like delay, message size, wait time, and visibility timeout. This can help optimize your queue management by adjusting these parameters for better performance and efficiency.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in your system where message queues are not protected by encryption, which could potentially expose sensitive data to unauthorized individuals.

```sql+postgres
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

```sql+sqlite
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
Identify instances where AWS Simple Queue Service (SQS) queues are encrypted with a Customer Master Key (CMK) for enhanced security measures. This can be useful to verify if the queues are following your organization's security protocols.

```sql+postgres
select
  title,
  kms_master_key_id,
  sqs_managed_sse_enabled
from
  aws_sqs_queue
where
  kms_master_key_id is not null;
```

```sql+sqlite
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
Explore which AWS Simple Queue Service (SQS) queues are encrypted using an SQS-owned key. This can be useful for understanding your encryption practices and ensuring that sensitive data is properly secured.

```sql+postgres
select
  title,
  kms_master_key_id,
  sqs_managed_sse_enabled
from
  aws_sqs_queue
where
  sqs_managed_sse_enabled;
```

```sql+sqlite
select
  title,
  kms_master_key_id,
  sqs_managed_sse_enabled
from
  aws_sqs_queue
where
  sqs_managed_sse_enabled = 1;
```

### List queues with a message retention period less than 7 days
Discover the segments that have a message retention period of less than a week. This query is useful to identify potential areas of data loss in your AWS Simple Queue Service (SQS) due to short retention periods.

```sql+postgres
select
  title,
  message_retention_seconds
from
  aws_sqs_queue
where
  message_retention_seconds < '604800';
```

```sql+sqlite
select
  title,
  message_retention_seconds
from
  aws_sqs_queue
where
  message_retention_seconds < 604800;
```

### List queues which are not configured with a dead-letter queue (DLQ)
Determine the areas in your system where queues are lacking a dead-letter queue (DLQ) configuration. This can be useful for identifying potential points of failure where messages could be lost.

```sql+postgres
select
  title,
  redrive_policy
from
  aws_sqs_queue
where
  redrive_policy is null;
```

```sql+sqlite
select
  title,
  redrive_policy
from
  aws_sqs_queue
where
  redrive_policy is null;
```

### List FIFO queues
Discover the segments that utilize first-in, first-out (FIFO) queues in AWS Simple Queue Service (SQS), allowing you to better manage and prioritize tasks in your applications.

```sql+postgres
select
  title,
  fifo_queue
from
  aws_sqs_queue
where
  fifo_queue;
```

```sql+sqlite
select
  title,
  fifo_queue
from
  aws_sqs_queue
where
  fifo_queue = 1;
```

### List queues with policy statements that grant cross-account access
Discover the segments that have policy statements granting cross-account access within your queue system. This can be useful in identifying potential security risks and ensuring proper access management.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List queues with policy statements that grant anoymous access
Determine the areas in your AWS SQS queues where policy statements permit anonymous access. This is useful for identifying potential security vulnerabilities and ensuring that your queues are only accessible to authorized users.

```sql+postgres
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

```sql+sqlite
select
  title,
  p as principal,
  a as action,
  json_extract(s, '$.Effect') as effect,
  json_extract(s, '$.Condition') as conditions
from
  aws_sqs_queue,
  json_each(json_extract(policy_std, '$.Statement')) as s,
  json_each(json_extract(json_extract(s.value, '$.Principal'), '$.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a
where
  p.value = '*'
  and json_extract(s.value, '$.Effect') = 'Allow';
```

### List queues with policy statements that grant full access (sqs:*)
Determine the areas in your AWS SQS queues where policy statements grant full access. This is useful for identifying potential security risks and ensuring that access permissions are appropriately restricted.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```


