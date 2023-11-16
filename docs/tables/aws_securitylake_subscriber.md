---
title: "Table: aws_securitylake_subscriber - Query AWS Security Lake Subscriber using SQL"
description: "Allows users to query AWS Security Lake Subscriber data, providing information about each subscriber's details in the AWS Security Lake service. This includes subscriber status, endpoint type, and subscription creation time."
---

# Table: aws_securitylake_subscriber - Query AWS Security Lake Subscriber using SQL

The `aws_securitylake_subscriber` table in Steampipe provides information about subscribers within the AWS Security Lake service. This table allows DevOps engineers, security analysts, and other technical professionals to query subscriber-specific details, including the subscriber's status, endpoint type, and subscription creation time. Users can utilize this table to gather insights on subscribers, such as their current status, the type of endpoint they are subscribed to, and when they were created. The schema outlines the various attributes of the AWS Security Lake Subscriber, including the subscriber ARN, endpoint type, status, and creation time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securitylake_subscriber` table, you can use the `.inspect aws_securitylake_subscriber` command in Steampipe.

**Key columns**:

- `account_id`: The AWS Account ID in which the Security Lake Subscriber resides. This is useful for correlating with other tables that contain account-level data.
- `subscriber_arn`: The ARN of the Security Lake Subscriber. This is a unique identifier and can be used to join with other tables that reference the subscriber.
- `creation_time`: The time when the subscription was created. This can be useful for tracking the lifecycle of subscriptions.

## Examples

### Basic info

```sql
select
  subscriber_name,
  subscription_id,
  created_at,
  role_arn,
  s3_bucket_arn,
  subscription_endpoint
from
  aws_securitylake_subscriber;
```

### List subscribers older than 30 days

```sql
select
  subscriber_name,
  subscription_id,
  created_at,
  role_arn,
  s3_bucket_arn,
  subscription_endpoint
from
  aws_securitylake_subscriber
where
  created_at <= created_at - interval '30' day;
```

## Get IAM role details for each subscriber

```sql
select
  s.subscriber_name,
  s.subscription_id,
  r.arn,
  r.inline_policies,
  r.attached_policy_arns,
  r.assume_role_policy
from
  aws_securitylake_subscriber as s,
  aws_iam_role as r
where
  s.role_arn = r.arn;
```

## Get S3 bucket details for each subscriber

```sql
select
  s.subscriber_name,
  s.subscription_id,
  b.arn,
  b.event_notification_configuration,
  b.server_side_encryption_configuration,
  b.acl
from
  aws_securitylake_subscriber as s,
  aws_s3_bucket as b
where
  s.s3_bucket_arn = b.arn;
```

## List subscribers that are not active

```sql
select
  subscriber_name,
  created_at,
  subscription_status,
  s3_bucket_arn,
  sns_arn
from
  aws_securitylake_subscriber
where
  subscription_status <> 'ACTIVE';
```