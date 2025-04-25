---
title: "Steampipe Table: aws_shield_subscription - Query AWS Shield Advanced Subscription using SQL"
description: "Allow users to query their AWS Shield Advanced Subscription details, such as the start and end dateof the subscription or the status of the proactive engagement of the Shield Response Team."
folder: "Shield"
---

# Table: aws_shield_subscription - Query AWS Shield Advanced Subscription using SQL

AWS Shield Advanced is a DDoS protection service from AWS. For a monthly fee, Shield Advanced will protect your AWS resources against Distributed Denial of Service (DDoS) attacks.

## Table Usage Guide

The `aws_shield_subscription` table in Steampipe allows you to query the current status of the AWS Shield Advanced Subscription of your account. This table provides you with insights into the start and end date of your subscription, the subscription limits or the status of the proactive engagement of the Shield Response Team. For more information about the individual fields, please refer to the [AWS Shield Advanced API documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeSubscription.html).

## Examples

### Basic info

```sql+postgres
select
  subscription_state,
  start_time,
  end_time,
  auto_renew,
  proactive_engagement_status
from
  aws_shield_subscription;
```

```sql+sqlite
select
  subscription_state,
  start_time,
  end_time,
  auto_renew,
  proactive_engagement_status
from
  aws_shield_subscription;
```

### Check if the subscription is active and proactive engagement is enabled

```sql+postgres
select
  subscription_state,
  proactive_engagement_status
from
  aws_shield_subscription
where
  subscription_state = 'ACTIVE'
  and proactive_engagement_status = 'ENABLED';
```

```sql+sqlite
select
  subscription_state,
  proactive_engagement_status
from
  aws_shield_subscription
where
  subscription_state = 'ACTIVE'
  and proactive_engagement_status = 'ENABLED';
```
