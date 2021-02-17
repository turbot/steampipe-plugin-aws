# Table: aws_cloudwatch_event_rule

Amazon CloudWatch Events delivers a near real-time stream of system events that describe changes in Amazon Web Services (AWS) resources.

## Examples

### Basic info

```sql
select
  name,
  description,
  event_bus_name,
  state,
  region
from
  aws_cloudwatch_event_rule;
```


### List of cloudwatch event rules that routes IAM events to target

```sql
select
  name,
  jsonb_array_elements_text(event_pattern -> 'detail' -> 'eventName') as event_name
from
  aws_cloudwatch_event_rule,
  jsonb_array_elements_text(event_pattern -> 'source') as source
where
  source = 'aws.iam';
```


### Get the schedule expression of each event rules

```sql
select
  name,
  schedule_expression
from
  aws_cloudwatch_event_rule
where
  schedule_expression is not null;
```


### Count of events rule per region

```sql
select
  region,
  count(*) as event_rule_count
from
  aws_cloudwatch_event_rule
group by
  region;
```