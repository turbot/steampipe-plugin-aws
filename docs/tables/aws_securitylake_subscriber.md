# Table: aws_securitylake_subscriber

Amazon Security Lake is a fully-managed security data lake service. You can use Security Lake to automatically centralize security data from cloud, on-premises, and custom sources into a data lake that's stored in your account. Security Lake helps you analyze security data, so you can get a more complete understanding of your security posture across the entire organization and improve the protection of your workloads, applications, and data.

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