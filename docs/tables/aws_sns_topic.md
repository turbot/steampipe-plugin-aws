---
title: "Steampipe Table: aws_sns_topic - Query AWS SNS Topics using SQL"
description: "Allows users to query AWS SNS Topics to gather information about each topic, including its name, owner, ARN, and other related data."
folder: "SNS"
---

# Table: aws_sns_topic - Query AWS SNS Topics using SQL

The AWS Simple Notification Service (SNS) Topics is a web service that coordinates and manages the delivery or sending of messages to subscribing endpoints or clients. It provides a simple, cost-effective method to asynchronously distribute messages to a large number of endpoints, making it a fundamental part of the AWS messaging infrastructure. SNS Topics offer flexibility in terms of message delivery, allowing you to fan out messages to a large number of subscribers, including distributed systems and services, and mobile devices.

## Table Usage Guide

The `aws_sns_topic` table in Steampipe provides you with information about each topic in Amazon Simple Notification Service (SNS). This table allows you as a DevOps engineer to query topic-specific details, including the topic name, owner, ARN, and other associated metadata. You can utilize this table to gather insights on SNS topics, such as topic subscription details, policy attributes, and more. The schema outlines for you the various attributes of the SNS topic, including the topic ARN, owner, subscription count, and associated tags.

## Examples

### List of unencrypted SNS topic
Identify instances where Simple Notification Service (SNS) topics in AWS are not encrypted, which could potentially expose sensitive data. This information is crucial for improving security measures and ensuring data privacy.

```sql+postgres
select
  title,
  kms_master_key_id
from
  aws_sns_topic
where
  kms_master_key_id is null;
```

```sql+sqlite
select
  title,
  kms_master_key_id
from
  aws_sns_topic
where
  kms_master_key_id is null;
```


### List of SNS topics which are not using Customer Managed Keys(CMK)
Identify instances where Simple Notification Service (SNS) topics are not secured with Customer Managed Keys (CMK). This is useful to ensure all your SNS topics have the added security layer of using your own encryption keys.

```sql+postgres
select
  title,
  kms_master_key_id
from
  aws_sns_topic
where
  kms_master_key_id = 'alias/aws/sns';
```

```sql+sqlite
select
  title,
  kms_master_key_id
from
  aws_sns_topic
where
  kms_master_key_id = 'alias/aws/sns';
```


### List of SNS topics without owner tag key
Discover the segments that have SNS topics without an assigned owner. This could be useful in managing and organizing your AWS resources more efficiently.

```sql+postgres
select
  title,
  tags
from
  aws_sns_topic
where
  not tags :: JSONB ? 'owner';
```

```sql+sqlite
select
  title,
  tags
from
  aws_sns_topic
where
  json_extract(tags, '$.owner') is null;
```


### List of SNS topics policy statements that grant anonymous access
Identify instances where anonymous access is granted in SNS topics policy statements. This is useful to uncover potential security risks, as unrestricted access might lead to unauthorized use or data breaches.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(p.value, '$') as principal,
  json_extract(a.value, '$') as action,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Condition') as conditions
from
  aws_sns_topic,
  json_each(json_extract(policy_std, '$.Statement')) as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a
where
  json_extract(p.value, '$') = '*'
  and json_extract(s.value, '$.Effect') = 'Allow';
```


### Topic policy statements that grant full access to the resource
Determine the areas in which topic policy statements are granting full access to the resource. This can help in assessing the security implications and managing access control effectively.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(p.value, '$') as principal,
  json_extract(a.value, '$') as action,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Condition') as conditions
from
  aws_sns_topic,
  json_each(json_extract(policy_std, '$.Statement')) as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a
where
  json_extract(s.value, '$.Effect') = 'Allow'
  and (
    json_extract(a.value, '$') = '*'
    or json_extract(a.value, '$') = 'sns:*'
  );
```

### List of topics that DO NOT enforce encryption in transit
Identify instances where certain topics do not enforce encryption in transit, which could pose potential security risks. This is useful for maintaining data privacy and meeting compliance requirements.

```sql+postgres
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

```sql+sqlite
select
  title
from
  aws_sns_topic
where
  title not in (
    select
      aws_sns_topic.title
    from
      aws_sns_topic,
      json_each(json_extract(policy_std, '$.Statement')) as s,
      json_each(json_extract(s.value, '$.Principal.AWS')) as p,
      json_each(json_extract(s.value, '$.Action')) as a,
      json_each(json_extract(s.value, '$.Condition.Bool."aws:securetransport"')) as ssl
    where
      json_extract(p.value, '$') = '*'
      and json_extract(s.value, '$.Effect') = 'Deny'
      and json_extract(ssl.value, '$') = 'false'
  );
```

### List topics which have delivery status logging for notification messages disabled
Identify instances where certain topics in your AWS SNS service have disabled delivery status logging for notification messages. This can be useful for auditing purposes or to rectify potential communication issues.

```sql+postgres
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

```sql+sqlite
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