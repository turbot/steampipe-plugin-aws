---
title: "Steampipe Table: aws_rolesanywhere_trust_anchor - Query AWS Roles Anywhere Trust Anchors using SQL"
description: "Allows users to query Roles Anywhere for detailed information about the Trust Anchor configurations."
folder: "Roles Anywhere"
---

# Table: aws_rolesanywhere_trust_anchor - Query AWS Roles Anywhere Trust Anchors using SQL

AWS Roles Anywhere enables trusted entities outside of an AWS account to carry out operations within that account by using certificates for IAM Role assumption. Trust Anchors contain the certificate bundle information used for those Role session creations.

## Table Usage Guide

The `aws_rolesanywhere_trust_anchor` table in Steampipe provides you with information about Trust Anchors with Roles Anywhere. This table allows you, as a DevOps engineer, to query Anchor-specific details, including certificate bundles, expiry notification settings, and associated metadata. You can utilize this table to gather insights on Trust Anchors, such as Anchors configured to use self-signed bundles vs ACM bundles, certificates not configured to notify on expiration, and more. The schema outlines the various attributes of the Trust Anchor for you, including the ARN, certificate data, and create and update times.

## Examples

### List enabled Trust Anchors.
Determine the Trust Anchors that are currently enabled. 
This can be useful to determine which certificates can access the account.

```sql+postgres
select
  arn,
  source_type, 
  source_data
from
  aws_rolesanywhere_trust_anchor
where
  enabled;
```

```sql+sqlite
select
  arn,
  source_type, 
  source_data
from
  aws_rolesanywhere_trust_anchor
where
  enabled;
```

### List Trust Anchors not using an ACM Private CA.
Determine the Trust Anchors that are not configured to use a Certificate Manager (ACM) Private Certificate Authority (PCA).
This can be useful to determine which Anchors are using an uploaded certificate bundle.

```sql+postgres
select
  arn
from
  aws_rolesanywhere_trust_anchor
where
  source_type <> 'AWS_ACM_PCA';
```

```sql+sqlite
select
  arn
from
  aws_rolesanywhere_trust_anchor
where
  source_type <> 'AWS_ACM_PCA';
```

### List enabled Trust Anchor notifications
Determine the notification events that are enabled for each Trust Anchor.
This can be useful to determine which expiry events will trigger a notification.

```sql+postgres
select
  arn, 
  notification_setting -> 'Event' as event
from
  aws_rolesanywhere_trust_anchor as anchor,
  jsonb_array_elements(notification_settings) as notification_setting
where
  notification_setting -> 'Enabled' = 'true'
```

```sql+sqlite
select
  arn, 
  json_extract(notification_setting, '$.Event') as event
from
  aws_rolesanywhere_trust_anchor as anchor,
  json_each(notification_settings) as notification_setting
where
  json_extract(notification_setting, '$.Enabled') = 'true'
```
