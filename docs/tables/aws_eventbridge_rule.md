---
title: "Steampipe Table: aws_eventbridge_rule - Query AWS EventBridge Rule using SQL"
description: "Allows users to query AWS EventBridge Rule to access information regarding the EventBridge rules defined within an AWS account."
folder: "EventBridge"
---

# Table: aws_eventbridge_rule - Query AWS EventBridge Rule using SQL

The AWS EventBridge Rule is a component of Amazon EventBridge, a serverless event bus service that makes it easy to connect applications together using data from your own applications, integrated Software-as-a-Service (SaaS) applications, and AWS services. EventBridge delivers a stream of real-time data from event sources and routes that data to targets like AWS Lambda. It primarily ingests, filters, and delivers events so you can build new applications quickly, and get to market faster.

## Table Usage Guide

The `aws_eventbridge_rule` table in Steampipe provides you with information about EventBridge rules within AWS EventBridge. This table allows you, as a DevOps engineer, to query rule-specific details, including the rule name, ARN, state, description, schedule expression, and associated metadata. You can utilize this table to gather insights on rules, such as the rules associated with a specific event bus, the state of the rules (whether they are enabled or disabled), and more. The schema outlines the various attributes of the EventBridge rule for you, including the rule ARN, event bus name, description, state, and associated tags.

## Examples

### Basic info
Gain insights into the status and origins of your AWS EventBridge rules. This query is particularly useful for auditing and maintaining an overview of your EventBridge settings.

```sql+postgres
select
  name,
  arn,
  state,
  created_by,
  event_bus_name
from
  aws_eventbridge_rule;
```

```sql+sqlite
select
  name,
  arn,
  state,
  created_by,
  event_bus_name
from
  aws_eventbridge_rule;
```


### List disabled rules
Determine the areas in which AWS EventBridge rules are not active. This is useful for identifying potential gaps in event-driven workflows or areas where automation may have been turned off.

```sql+postgres
select
  name,
  arn,
  state,
  created_by
from
  aws_eventbridge_rule
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
  aws_eventbridge_rule
where
  state != 'ENABLED';
```


### Get the target information for each rule
This query allows you to identify the target details for each rule in your AWS EventBridge service. It's useful for auditing and understanding the relationships and dependencies between different rules and their targets within your AWS infrastructure.

```sql+postgres
select
  name,
  cd ->> 'Id' as target_id,
  cd ->> 'Arn' as target_arn,
  cd ->> 'RoleArn' as role_arn
from
  aws_eventbridge_rule,
  jsonb_array_elements(targets) as cd;
```

```sql+sqlite
select
  name,
  json_extract(cd.value, '$.Id') as target_id,
  json_extract(cd.value, '$.Arn') as target_arn,
  json_extract(cd.value, '$.RoleArn') as role_arn
from
  aws_eventbridge_rule,
  json_each(targets) as cd;
```