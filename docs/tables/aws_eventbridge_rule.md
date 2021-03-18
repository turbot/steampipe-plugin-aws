# Table: aws_eventbridge_rule

Amazon EventBridge defines the building and management of event-driven applications by taking care of event ingestion and delivery, security, authorization, and error-handling.

## Examples

### Basic info

```sql
select
  name,
  arn,
  state,
  created_by,
  event_bus_name
from
  aws_eventbridge_rule;
```


### List disabled rules

```sql
select
  name,
  arn,
  state,
  created_by
from
  aws_eventbridge_rule
where
  state != 'ENABLED';
```


### Get the target information for each rule

```sql
select
  name,
  cd ->> 'Id' as target_id,
  cd ->> 'Arn' as target_arn,
  cd ->> 'RoleArn' as role_arn
from
  aws_eventbridge_rule,
  jsonb_array_elements(targets) as cd;
```
