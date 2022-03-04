# Table: aws_sqs_queue

Amazon Simple Queue Service (SQS) is a fully managed message queuing service that enables to decouple and scale micro services, distributed systems, and serverless applications

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
