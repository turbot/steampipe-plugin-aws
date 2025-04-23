---
title: "Steampipe Table: aws_securitylake_subscriber - Query AWS Security Lake Subscriber using SQL"
description: "Allows users to query AWS Security Lake Subscriber data, providing information about each subscriber's details in the AWS Security Lake service. This includes subscriber status, endpoint type, and subscription creation time."
folder: "Security Lake"
---

# Table: aws_securitylake_subscriber - Query AWS Security Lake Subscriber using SQL

The AWS Security Lake Subscriber is a component of AWS Lake Formation, a service that makes it easy to set up, secure, and manage your data lake. It helps in subscribing to data access events in the data lake, enabling granular control over who has access to specific data. It allows monitoring, auditing, and receiving notifications about specific activities in your AWS data lake.

## Table Usage Guide

The `aws_securitylake_subscriber` table in Steampipe provides you with information about subscribers within the AWS Security Lake service. This table allows you, as a DevOps engineer, security analyst, or other technical professional, to query subscriber-specific details, including the subscriber's status, endpoint type, and subscription creation time. You can utilize this table to gather insights on subscribers, such as their current status, the type of endpoint they are subscribed to, and when they were created. The schema outlines the various attributes of the AWS Security Lake Subscriber for you, including the subscriber ARN, endpoint type, status, and creation time.

## Examples

### Basic info
Determine the areas in which subscribers are interacting with your AWS security system. This query provides insights into the creation time, subscription endpoints, and associated roles of each subscriber, allowing you to better understand your user base and their access levels.

```sql+postgres
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

```sql+sqlite
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
Identify subscribers who have been part of the system for over a month. This is useful to understand user retention and recognize long-term subscribers.

```sql+postgres
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

```sql+sqlite
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
  created_at <= datetime(created_at, '-30 day');
```

### Get IAM role details for each subscriber
Analyze the access policies of each subscriber's IAM role. This helps in auditing security and access controls.

```sql+postgres
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

```sql+sqlite
select
  s.subscriber_name,
  s.subscription_id,
  r.arn,
  r.inline_policies,
  r.attached_policy_arns,
  r.assume_role_policy
from
  aws_securitylake_subscriber as s
join
  aws_iam_role as r
on
  s.role_arn = r.arn;
```

### Get S3 bucket details for each subscriber
Review the configuration of the S3 bucket linked to each subscriber. This aids in ensuring proper storage setup and security measures.

```sql+postgres
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

```sql+sqlite
select
  s.subscriber_name,
  s.subscription_id,
  b.arn,
  b.event_notification_configuration,
  b.server_side_encryption_configuration,
  b.acl
from
  aws_securitylake_subscriber as s
join
  aws_s3_bucket as b
on
  s.s3_bucket_arn = b.arn;
```

### List subscribers that are not active
Discover subscribers who aren't currently active. This can assist in identifying potential issues or areas for improvement in user engagement.

```sql+postgres
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

```sql+sqlite
select
  subscriber_name,
  created_at,
  subscription_status,
  s3_bucket_arn,
  sns_arn
from
  aws_securitylake_subscriber
where
  subscription_status != 'ACTIVE';
```