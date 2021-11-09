# Table: aws_eventbridge_bus

Amazon EventBridge is a serverless event bus that makes it easier to build event-driven applications at scale using events generated from your applications, integrated Software-as-a-Service (SaaS) applications, and AWS services. An event bus receives events from a source and routes them to rules associated with that event bus.

## Examples

### Basic info

```sql
select
  name,
  arn,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_eventbridge_bus;
```
