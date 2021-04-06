# Table: aws_redshift_event_subscription

Amazon Redshift Event Subscription notifies when an event occurs for a cluster, snapshot, security group or parameter group.

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
