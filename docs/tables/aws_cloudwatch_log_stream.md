# Table: aws_cloudwatch_log_stream

A log stream is a sequence of log events that share the same source. A log stream is generally intended to represent the sequence of events coming from the application instance or resource being monitored.

## Examples

### Log stream basic info

```sql
select
  name,
  arn,
  log_group_name
from
  aws_cloudwatch_log_stream;
```
