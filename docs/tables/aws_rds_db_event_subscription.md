# Table: aws_rds_db_event_subscription

Amazon RDS uses the Amazon Simple Notification Service (Amazon SNS) to provide notification when an Amazon RDS event occurs. These notifications can be in any notification form supported by Amazon SNS for an AWS Region, such as an email, a text message, or a call to an HTTP endpoint.

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

### List DB event subscription which are enable

```sql
select
  cust_subscription_id,
  enabled
from
  aws_rds_db_event_subscription
where
  enabled;
```