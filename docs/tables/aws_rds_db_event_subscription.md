---
title: "Steampipe Table: aws_rds_db_event_subscription - Query AWS RDS DB Event Subscriptions using SQL"
description: "Allows users to query AWS RDS DB Event Subscriptions to retrieve information about all event subscriptions for RDS DB instances."
folder: "RDS"
---

# Table: aws_rds_db_event_subscription - Query AWS RDS DB Event Subscriptions using SQL

The AWS RDS DB Event Subscription is a feature of Amazon Relational Database Service (RDS) that allows you to receive notifications when specific database events occur. These events can include failovers, backups, configurations changes, and more. By creating an event subscription, you can ensure that you are promptly informed about changes that could impact your database operations.

## Table Usage Guide

The `aws_rds_db_event_subscription` table in Steampipe provides you with information about event subscriptions within Amazon RDS. This table allows you, as a DevOps engineer, to query event subscription-specific details, including the associated RDS DB instances, the types of events the subscription applies to, and the notification methods for those events. You can utilize this table to monitor the status of your RDS DB instances, manage event notifications, and ensure all event subscriptions are properly configured. The schema outlines the various attributes of the event subscription, including the subscription name, ARN, status, and associated tags.

## Examples

### Basic info
Explore the status and activation of your Amazon RDS event subscriptions. This can be helpful to ensure all necessary subscriptions are active and functioning as expected.

```sql+postgres
select
  cust_subscription_id,
  customer_aws_id,
  arn,
  status,
  enabled
from
  aws_rds_db_event_subscription;
```

```sql+sqlite
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
Explore which database event subscriptions are currently active in your AWS RDS service. This can help you manage your database events more effectively by identifying which subscriptions are currently receiving and processing events.

```sql+postgres
select
  cust_subscription_id,
  enabled
from
  aws_rds_db_event_subscription
where
  enabled;
```

```sql+sqlite
select
  cust_subscription_id,
  enabled
from
  aws_rds_db_event_subscription
where
  enabled = 1;
```