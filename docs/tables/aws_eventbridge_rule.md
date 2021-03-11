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