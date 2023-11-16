---
title: "Table: aws_rds_db_event_subscription - Query AWS RDS DB Event Subscriptions using SQL"
description: "Allows users to query AWS RDS DB Event Subscriptions to retrieve information about all event subscriptions for RDS DB instances."
---

# Table: aws_rds_db_event_subscription - Query AWS RDS DB Event Subscriptions using SQL

The `aws_rds_db_event_subscription` table in Steampipe provides information about event subscriptions within Amazon RDS. This table allows DevOps engineers to query event subscription-specific details, including the associated RDS DB instances, the types of events the subscription applies to, and the notification methods for those events. Users can utilize this table to monitor the status of their RDS DB instances, manage event notifications, and ensure all event subscriptions are properly configured. The schema outlines the various attributes of the event subscription, including the subscription name, ARN, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_event_subscription` table, you can use the `.inspect aws_rds_db_event_subscription` command in Steampipe.

### Key columns:

- `cust_subscription_id`: This is the customer-defined name of the event subscription. It can be used to join this table with other tables that also contain this identifier.
- `arn`: The Amazon Resource Name (ARN) of the event subscription. It provides a unique identifier for the subscription and can be used to join this table with other tables that reference the same ARN.
- `sns_topic_arn`: This is the ARN of the Amazon SNS topic that the event notification is sent to. It can be used to join this table with other tables that reference the same SNS topic.

## Examples

### Basic info

```sql
select
  cust_subscription_id,
  customer_aws_id,
  arn,
  status,
  enabled
from
  aws_rds_db_event_subscription;
```

### List enabled DB event subscription

```sql
select
  cust_subscription_id,
  enabled
from
  aws_rds_db_event_subscription
where
  enabled;
```