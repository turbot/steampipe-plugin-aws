# Table: aws_eventbridge_rule

Amazon EventBridge defines the building and management of event-driven applications by taking care of event ingestion and delivery, security, authorization, and error-handling. Events is serverless, highly-available, and scalable.

## Examples

### EventBridge rule basic info

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


### List of rules which are not enabled

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


### List of targets info which are associated with rule

```sql
select
  name,
  cd ->> 'Arn' as target_arn,
  cd ->> 'RoleArn' as role_arn,
  cd ->> 'Id' as id
from
  aws_eventbridge_rule,
  jsonb_array_elements(targets) as cd;
```


### List of rules which have iam role

```sql
select
  name,
  cd ->> 'RoleArn' as role_arn
from
  aws_eventbridge_rule,
  jsonb_array_elements(targets) as cd
where
  cd ->> 'RoleArn' != 'null';
```

