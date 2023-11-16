---
title: "Table: aws_sns_topic - Query AWS SNS Topics using SQL"
description: "Allows users to query AWS SNS Topics to gather information about each topic, including its name, owner, ARN, and other related data."
---

# Table: aws_sns_topic - Query AWS SNS Topics using SQL

The `aws_sns_topic` table in Steampipe provides information about each topic in Amazon Simple Notification Service (SNS). This table allows DevOps engineers to query topic-specific details, including the topic name, owner, ARN, and other associated metadata. Users can utilize this table to gather insights on SNS topics, such as topic subscription details, policy attributes, and more. The schema outlines the various attributes of the SNS topic, including the topic ARN, owner, subscription count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sns_topic` table, you can use the `.inspect aws_sns_topic` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the SNS Topic. It can be used to join this table with other AWS resource tables that reference the SNS Topic ARN.
- `owner`: The AWS account ID of the topic's owner. This can be helpful in scenarios where you need to join this table with other tables that contain owner information.
- `policy`: The policy that defines who can publish and subscribe to the topic. This can be used to join with other tables that need policy information for security analysis.

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
