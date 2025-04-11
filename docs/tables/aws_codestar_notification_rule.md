---
title: "Steampipe Table: aws_codestar_notification_rule - Query AWS CodeStar notification rules using SQL"
description: "Allows users to query CodeStar notification rules in the AWS Developer Tools to retrieve information about notification rules."
folder: "CodeStar"
---

# Table: aws_codestar_notification_rule - Query AWS CodeStar notification rules using SQL

The AWS CodeStar notification rules allow you to set up notifications for the AWS Developer Tools, including AWS CodePipeline and AWS CodeBuild, to various destinations including AWS SNS and AWS Chatbot.

## Table Usage Guide

The `aws_codestar_notification_rule` table in Steampipe provides you with information about notification rules. This table allows you, as a DevOps engineer, to query notification rule details, including the notification rule ARN, status, level of detail, enabled event types, as well as the ARN of the resource producing notifications and the notification targets. You can use this table to gather insights on notification rules, and combine it with other tables such as `aws_codepipeline_pipeline` to check notification rules are set up consistently.

## Examples

### Basic info
Review the configured rules and their status.

```sql+postgres
select
  name,
  resource,
  detail_type,
  status
from
  aws_codestar_notification_rule;
```

```sql+sqlite
select
  name,
  resource,
  detail_type,
  status
from
  aws_codestar_notification_rule;
```

### Identify which CI/CD pipelines have notification rules
Determine which AWS CodePipeline pipelines do or do not have associated notification rules.

```sql+postgres
select
  pipeline.name as pipeline,
  notification_rule.name notification_rule,
  notification_rule.status
from
  aws_codepipeline_pipeline as pipeline
  left join aws_codestar_notification_rule as notification_rule on pipeline.arn = notification_rule.resource;
```

```sql+sqlite
select
  pipeline.name as pipeline,
  notification_rule.name as notification_rule,
  notification_rule.status
from
  aws_codepipeline_pipeline as pipeline
  left join aws_codestar_notification_rule as notification_rule on pipeline.arn = notification_rule.resource;
```

### Check for notification rules with no targets
Determine which notification rules lack targets. This query uses PostgreSQL's JSON querying capabilities to count the number of targets configured.

```sql+postgres
select
  name
from
  aws_codestar_notification_rule
where
  jsonb_array_length(targets) = 0;
```

```sql+sqlite
select
  name
from
  aws_codestar_notification_rule
where
  json_array_length(targets) = 0;
```

### Name the SNS topics associated with notification rules
Determine which AWS SNS topics the notification rules are targeting. This query uses PostgreSQL's JSON querying capabilities to join on the notification rule targets. Note that due to the cross join, this query will not list notification rules that don't have any targets.

```sql+postgres
select
  notification_rule.name as notification_rule,
  target ->> 'TargetType' as target_type,
  topic.title as target_topic
from
  aws_codestar_notification_rule as notification_rule cross
  join jsonb_array_elements(notification_rule.targets) as target
  left join aws_sns_topic as topic on target ->> 'TargetAddress' = topic.topic_arn;
```

```sql+sqlite
select
  notification_rule.name as notification_rule,
  json_extract(target.value, '$.TargetType') as target_type,
  topic.title as target_topic
from
  aws_codestar_notification_rule as notification_rule
  cross join json_each(notification_rule.targets) as target
  left join aws_sns_topic as topic on json_extract(target.value, '$.TargetAddress') = topic.topic_arn;
```

### Using CTE to retain notification rules without targets
By using a Common Table Expression (`with` query), it is possible to join on targets without discarding notification rules that don't have any targets.

```sql+postgres
with rule_target as (
  select
    arn,
    target ->> 'TargetAddress' as target_address,
    target ->> 'TargetStatus' as target_status,
    target ->> 'TargetType' as target_type
  from
    aws_codestar_notification_rule cross
    join jsonb_array_elements(targets) as target
)
select
  notification_rule.name as notification_rule,
  rule_target.target_type,
  topic.title as target_topic
from
  aws_codestar_notification_rule as notification_rule
  left join rule_target on rule_target.arn = notification_rule.arn
  left join aws_sns_topic as topic on rule_target.target_address = topic.topic_arn;
```

```sql+sqlite
with rule_target as (
  select
    notification_rule.arn,
    json_extract(target.value, '$.TargetAddress') as target_address,
    json_extract(target.value, '$.TargetStatus') as target_status,
    json_extract(target.value, '$.TargetType') as target_type
  from
    aws_codestar_notification_rule as notification_rule
    cross join json_each(notification_rule.targets) as target
)
select
  notification_rule.name as notification_rule,
  rule_target.target_type,
  topic.title as target_topic
from
  aws_codestar_notification_rule as notification_rule
  left join rule_target on rule_target.arn = notification_rule.arn
  left join aws_sns_topic as topic on rule_target.target_address = topic.topic_arn;
```