---
title: "Table: aws_sns_topic_subscription - Query AWS Simple Notification Service (SNS) Topic Subscriptions using SQL"
description: "Allows users to query AWS SNS Topic Subscriptions to obtain detailed information about each subscription, including subscription ARN, owner, protocol, endpoint, and more."
---

# Table: aws_sns_topic_subscription - Query AWS Simple Notification Service (SNS) Topic Subscriptions using SQL

The `aws_sns_topic_subscription` table in Steampipe provides information about topic subscriptions within AWS Simple Notification Service (SNS). This table allows DevOps engineers to query subscription-specific details, including subscription ARN, owner, protocol, endpoint, and more. Users can utilize this table to gather insights on subscriptions, such as subscription status, delivery policy, raw message delivery, and more. The schema outlines the various attributes of the SNS topic subscription, including the subscription ARN, topic ARN, owner, protocol, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sns_topic_subscription` table, you can use the `.inspect aws_sns_topic_subscription` command in Steampipe.

### Key columns:

- `subscription_arn`: The ARN of the subscription. This is a unique identifier and can be used to join this table with other tables.
- `topic_arn`: The ARN of the topic associated with the subscription. This can be used to join with the `aws_sns_topic` table for more detailed topic information.
- `owner`: The AWS account ID of the subscription's owner. This can be used to join with `aws_iam_user' or `aws_iam_role` tables for more detailed user or role information.

## Examples

### List of subscriptions which are not configured with dead letter queue

```sql
select
  title,
  redrive_policy
from
  aws_sns_topic_subscription
where
  redrive_policy is null;
```


### List of subscriptions which are not configured to filter messages

```sql
select
  title,
  filter_policy
from
  aws_sns_topic_subscription
where
  filter_policy is null;
```


### Subscription count by topic arn

```sql
select
  title,
  count(subscription_arn) as subscription_count
from
  aws_sns_topic_subscription
group by
  title;
```
