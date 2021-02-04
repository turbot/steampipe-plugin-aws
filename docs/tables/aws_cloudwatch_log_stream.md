# Table: aws_cloudwatch_log_stream

A log stream is a sequence of log events that share the same source. A log stream is generally intended to represent the sequence of events coming from the application instance or resource being monitored.

## Examples

### Log stream basic info

```sql
select
  name,
  log_group_name,
  region
from
  aws_cloudwatch_log_stream;
```


### Count of log stream per log group

```sql
select
  log_group_name,
  count(*) as log_stream_count
from
  aws_cloudwatch_log_stream
group by
  log_group_name;
```