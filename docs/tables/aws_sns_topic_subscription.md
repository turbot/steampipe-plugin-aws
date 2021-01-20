# Table: aws_sns_topic_subscription

The AWS SNS Subscription resource subscribes an endpoint to an SNS topic.

## Examples

### List of subscriptions which are not configured with dead letter queue

```sql
select
  title,
  redrive_policy
from
  aws_sns_topic_subscription
where
  redrive_policy is null;
```


### List of subscriptions which are not configured to filter messages

```sql
select
  title,
  filter_policy
from
  aws_sns_topic_subscription
where
  filter_policy is null;
```


### Subscription count by topic arn

```sql
select
  title,
  count(subscription_arn) as subscription_count
from
  aws_sns_topic_subscription
group by
  title;
```
