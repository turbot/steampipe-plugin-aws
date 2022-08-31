# Table: aws_sns_topic

Amazon Simple Notification Service (Amazon SNS) is a fully managed messaging service for both application-to-application (A2A) and application-to-person (A2P) communication. SNS topic is a logical access point that acts as a communication channel.

## Examples

### List of unencrypted SNS topic

```sql
select
  title,
  kms_master_key_id
from
  aws_sns_topic
where
  kms_master_key_id is null;
```


### List of SNS topics which are not using Customer Managed Keys(CMK)

```sql
select
  title,
  kms_master_key_id
from
  aws_sns_topic
where
  kms_master_key_id = 'alias/aws/sns';
```


### List of SNS topics without owner tag key

```sql
select
  title,
  tags
from
  aws_sns_topic
where
  not tags :: JSONB ? 'owner';
```


### List of SNS topics policy statements that grant anonymous access

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_sns_topic,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  p = '*'
  and s ->> 'Effect' = 'Allow';
```


### Topic policy statements that grant full access to the resource

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_sns_topic,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and a in ('*', 'sns:*');
```


### List of topics that DO NOT enforce encryption in transit

```sql
select
  title
from
  aws_sns_topic
where
  title not in (
    select
      title
    from
      aws_sns_topic,
      jsonb_array_elements(policy_std -> 'Statement') as s,
      jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
      jsonb_array_elements_text(s -> 'Action') as a,
      jsonb_array_elements_text(
        s -> 'Condition' -> 'Bool' -> 'aws:securetransport'
      ) as ssl
    where
      p = '*'
      and s ->> 'Effect' = 'Deny'
      and ssl :: bool = false
  );
```

### List topics which have delivery status logging for notification messages disabled

```sql
select
  title,
  topic_arn,
  region
from 
  aws_sns_topic
where
  application_failure_feedback_role_arn is null and
  firehose_failure_feedback_role_arn is null and
  http_failure_feedback_role_arn is null and
  lambda_failure_feedback_role_arn is null and
  sqs_failure_feedback_role_arn is null;
```
