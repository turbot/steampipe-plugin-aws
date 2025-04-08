---
title: "Steampipe Table: aws_redshift_event_subscription - Query AWS Redshift Event Subscriptions using SQL"
description: "Allows users to query AWS Redshift Event Subscriptions, providing insights into the subscription's configuration, status, and associated Redshift clusters."
folder: "Redshift"
---

# Table: aws_redshift_event_subscription - Query AWS Redshift Event Subscriptions using SQL

The AWS Redshift Event Subscription is a feature of Amazon Redshift that allows you to subscribe to events related to your clusters, snapshots, security groups, and parameter groups. This service sends notifications to the Amazon Simple Notification Service (SNS) when specific events occur, enabling you to automate responses to these events. It provides a streamlined process for managing and responding to events in your Amazon Redshift environment.

## Table Usage Guide

The `aws_redshift_event_subscription` table in Steampipe provides you with information about event subscriptions within Amazon Redshift. This table allows you to query event subscription-specific details, including the subscription's configuration, status, and associated Redshift clusters. You can utilize this table to gather insights on event subscriptions, such as the event categories subscribed to, the status of the subscription, and the Redshift clusters associated with the subscription. The schema outlines the various attributes of the Redshift event subscription, including the subscription name, enabled status, event categories, source type, severity, and associated tags.

## Examples

### Basic info
Explore which AWS Redshift event subscriptions are active, by assessing elements like their status and creation time. This can help in managing and monitoring your AWS resources effectively.

```sql+postgres
select
  cust_subscription_id,
  customer_aws_id,
  status,
  sns_topic_arn,
  subscription_creation_time
from
  aws_redshift_event_subscription;
```

```sql+sqlite
select
  cust_subscription_id,
  customer_aws_id,
  status,
  sns_topic_arn,
  subscription_creation_time
from
  aws_redshift_event_subscription;
```


### List disabled event subscriptions
Identify instances where event subscriptions have been disabled in AWS Redshift. This is useful for auditing purposes, ensuring that all necessary subscriptions are active and functioning as expected.

```sql+postgres
select
  cust_subscription_id,
  customer_aws_id,
  status,
  enabled,
  sns_topic_arn,
  subscription_creation_time
from
  aws_redshift_event_subscription
where
  enabled is false;
```

```sql+sqlite
select
  cust_subscription_id,
  customer_aws_id,
  status,
  enabled,
  sns_topic_arn,
  subscription_creation_time
from
  aws_redshift_event_subscription
where
  enabled = 0;
```


### Get associated source details for each event subscription
Determine the areas in which event subscriptions are associated with different sources in your AWS Redshift environment. This can help prioritize and manage events based on their source and severity.

```sql+postgres
select
  cust_subscription_id,
  severity,
  source_type,
  event_categories_list,
  source_ids_list
from
  aws_redshift_event_subscription;
```

```sql+sqlite
select
  cust_subscription_id,
  severity,
  source_type,
  event_categories_list,
  source_ids_list
from
  aws_redshift_event_subscription;
```


### List unencrypted SNS topics associated with each event subscription
Explore which event subscriptions are associated with unencrypted SNS topics. This can help identify potential security risks in your AWS Redshift environment.

```sql+postgres
select
  e.cust_subscription_id,
  e.status,
  s.kms_master_key_id,
  s.topic_arn as arn
from
  aws_redshift_event_subscription as e
  join aws_sns_topic as s on s.topic_arn = e.sns_topic_arn
where
  s.kms_master_key_id is null;
```

```sql+sqlite
select
  e.cust_subscription_id,
  e.status,
  s.kms_master_key_id,
  s.topic_arn as arn
from
  aws_redshift_event_subscription as e
  join aws_sns_topic as s on s.topic_arn = e.sns_topic_arn
where
  s.kms_master_key_id is null;
```