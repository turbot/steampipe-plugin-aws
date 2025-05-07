---
title: "Steampipe Table: aws_cloudwatch_event_rule - Query AWS CloudWatch Event Rule using SQL"
description: "Allows users to query AWS CloudWatch Event Rule to access information regarding the event rules defined within an AWS account."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_event_rule - Query AWS CloudWatch Event Rule using SQL

AWS CloudWatch Events delivers a near real-time stream of system events that describe changes in AWS resources. CloudWatch Events allows you to route events to targets like AWS Lambda functions, Amazon SNS topics, Amazon SQS queues, or built-in targets. CloudWatch Events allows you to set up rules to trigger automated actions when an event matches your rule.

## Table Usage Guide

The `aws_cloudwatch_event_rule` table in Steampipe provides you with information about CloudWatch Event rules within AWS CloudWatch Events. This table allows you, as a DevOps engineer, to query rule-specific details, including the rule name, ARN, state, description, schedule expression, and associated metadata. You can utilize this table to gather insights on rules, such as the rules associated with a specific event bus, the state of the rules (whether they are enabled or disabled), and more. The schema outlines the various attributes of the CloudWatch Event rule for you, including the rule ARN, event bus name, description, state, and associated tags.

## Examples

### Basic info
Gain insights into the status and origins of your AWS CloudWatch Events rules. This query is particularly useful for auditing and maintaining an overview of your event settings.

```sql+postgres
select
  name,
  arn,
  state,
  created_by,
  event_bus_name
from
  aws_cloudwatch_event_rule;
```

```sql+sqlite
select
  name,
  arn,
  state,
  created_by,
  event_bus_name
from
  aws_cloudwatch_event_rule;
```

### List disabled rules
Determine the areas in which AWS CloudWatch Events rules are not active. This is useful for identifying potential gaps in event-driven workflows or areas where automation may have been turned off.

```sql+postgres
select
  name,
  arn,
  state,
  created_by
from
  aws_cloudwatch_event_rule
where
  state != 'ENABLED';
```

```sql+sqlite
select
  name,
  arn,
  state,
  created_by
from
  aws_cloudwatch_event_rule
where
  state != 'ENABLED';
```

### Get the target information for each rule
This query allows you to identify the target details for each rule in your AWS CloudWatch Events service. It's useful for auditing and understanding the relationships and dependencies between different rules and their targets within your AWS infrastructure.

```sql+postgres
select
  name,
  cd ->> 'Id' as target_id,
  cd ->> 'Arn' as target_arn,
  cd ->> 'RoleArn' as role_arn
from
  aws_cloudwatch_event_rule,
  jsonb_array_elements(targets) as cd;
```

```sql+sqlite
select
  name,
  json_extract(cd.value, '$.Id') as target_id,
  json_extract(cd.value, '$.Arn') as target_arn,
  json_extract(cd.value, '$.RoleArn') as role_arn
from
  aws_cloudwatch_event_rule,
  json_each(targets) as cd;
```

### List CloudWatch Event rules with schedule expressions
Identify rules that are triggered on a schedule rather than by events, which is useful for understanding time-based automation in your environment.

```sql+postgres
select
  name,
  schedule_expression,
  state,
  description
from
  aws_cloudwatch_event_rule
where
  schedule_expression is not null;
```

```sql+sqlite
select
  name,
  schedule_expression,
  state,
  description
from
  aws_cloudwatch_event_rule
where
  schedule_expression is not null;
```

### Find rules with Lambda function targets
Discover which CloudWatch Event rules are triggering Lambda functions, helping to map out serverless event-driven architectures in your AWS environment.

```sql+postgres
select
  r.name as rule_name,
  t ->> 'Id' as target_id,
  t ->> 'Arn' as lambda_arn
from
  aws_cloudwatch_event_rule as r,
  jsonb_array_elements(r.targets) as t
where
  t ->> 'Arn' like 'arn:aws:lambda:%';
```

```sql+sqlite
select
  r.name as rule_name,
  json_extract(t.value, '$.Id') as target_id,
  json_extract(t.value, '$.Arn') as lambda_arn
from
  aws_cloudwatch_event_rule as r,
  json_each(r.targets) as t
where
  json_extract(t.value, '$.Arn') like 'arn:aws:lambda:%';
```
