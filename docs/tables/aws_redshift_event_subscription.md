---
title: "Table: aws_redshift_event_subscription - Query AWS Redshift Event Subscriptions using SQL"
description: "Allows users to query AWS Redshift Event Subscriptions, providing insights into the subscription's configuration, status, and associated Redshift clusters."
---

# Table: aws_redshift_event_subscription - Query AWS Redshift Event Subscriptions using SQL

The `aws_redshift_event_subscription` table in Steampipe provides information about event subscriptions within Amazon Redshift. This table allows users to query event subscription-specific details, including the subscription's configuration, status, and associated Redshift clusters. Users can utilize this table to gather insights on event subscriptions, such as the event categories subscribed to, the status of the subscription, and the Redshift clusters associated with the subscription. The schema outlines the various attributes of the Redshift event subscription, including the subscription name, enabled status, event categories, source type, severity, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshift_event_subscription` table, you can use the `.inspect aws_redshift_event_subscription` command in Steampipe.

### Key columns:

- `cust_subscription_id`: This is the unique identifier of the event subscription. It can be used to join this table with other tables that contain information about Redshift event subscriptions.
- `enabled`: This column indicates whether the event subscription is enabled or not. It's useful in determining the status of the event subscription.
- `sns_topic_arn`: This is the Amazon Resource Name (ARN) of the Amazon SNS topic used by the event subscription. This can be used to join this table with other tables that contain information about SNS topics.

## Examples

### Basic info

```sql
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

```sql
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


### Get associated source details for each event subscription

```sql
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

```sql
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
