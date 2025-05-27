---
title: "Steampipe Table: aws_sns_subscription - Query AWS Simple Notification Service (SNS) Topic Subscriptions using SQL"
description: "Allows users to query AWS SNS Topic Subscriptions to obtain detailed information about each subscription, including subscription ARN, owner, protocol, endpoint, and more."
folder: "SNS"
---

# Table: aws_sns_subscription - Query AWS Simple Notification Service (SNS) Topic Subscriptions using SQL

The AWS Simple Notification Service (SNS) Topic Subscriptions allow you to manage and handle messages published to topics. Subscriptions define the endpoints to which messages will be delivered, allowing for the decoupling of microservices, distributed systems, and serverless applications. AWS SNS Topic Subscriptions support a variety of protocols including HTTP, HTTPS, Email, Email-JSON, SQS, Application, Lambda, and SMS.

## Table Usage Guide

The `aws_sns_subscription` table in Steampipe provides you with information about topic subscriptions within AWS Simple Notification Service (SNS). This table allows you, as a DevOps engineer, to query subscription-specific details, including subscription ARN, owner, protocol, endpoint, and more. You can utilize this table to gather insights on subscriptions, such as subscription status, delivery policy, raw message delivery, and more. The schema outlines the various attributes of the SNS topic subscription for you, including the subscription ARN, topic ARN, owner, protocol, and associated tags.

## Examples

### List of subscriptions which are not configured with dead letter queue
Determine the areas in which AWS SNS Topic subscriptions lack a configured dead letter queue. This is useful for identifying potential points of failure in message delivery, as messages could be lost if the subscription service is unavailable and there is no dead letter queue set up.

```sql+postgres
select
  title,
  redrive_policy
from
  aws_sns_subscription
where
  redrive_policy is null;
```

```sql+sqlite
select
  title,
  redrive_policy
from
  aws_sns_subscription
where
  redrive_policy is null;
```

### List of subscriptions which are not configured to filter messages
Determine the areas in which subscriptions are not set up to filter messages. This is beneficial for identifying potential inefficiencies or areas of improvement within your notification system.

```sql+postgres
select
  title,
  filter_policy
from
  aws_sns_subscription
where
  filter_policy is null;
```

```sql+sqlite
select
  title,
  filter_policy
from
  aws_sns_subscription
where
  filter_policy is null;
```

### List subscription count by topic arn
Determine the areas in which your AWS SNS topics are gaining the most traction by analyzing the number of subscriptions each topic has. This can help prioritize content creation and resource allocation for popular topics.

```sql+postgres
select
  title,
  count(subscription_arn) as subscription_count
from
  aws_sns_subscription
group by
  title;
```

```sql+sqlite
select
  title,
  count(subscription_arn) as subscription_count
from
  aws_sns_subscription
group by
  title;
```