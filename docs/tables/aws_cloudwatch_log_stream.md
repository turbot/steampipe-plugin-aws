# Table: aws_cloudwatch_log_stream

A log stream is a sequence of log events that share the same source. A log stream is generally intended to represent the sequence of events coming from the application instance or resource being monitored.

**Important notes:**
- To enhance performance, it is recommended to utilize the optional qualifiers `name`, `log_stream_name_prefix`, `descending`, and `order_by` for result set limitation.
- It's important to note that the columns `name` and `log_stream_name_prefix` cannot be specified together. If both are included as query parameters in the `where` clause, the `name` parameter value will be overridden by the `log_stream_name_prefix` parameter value in the input.
- The value of the `order_by` column can be either `LogStreamName` or `LastEventTime`. If the value is `LogStreamName`, the results are ordered by log stream name. If the value is `LastEventTime`, the results are ordered by the event time. The default value is LogStreamName. If you order the results by event time, you cannot specify the logStreamNamePrefix parameter. LastEventTimestamp represents the time of the most recent log event in the log stream in CloudWatch Logs. This number is expressed as the number of milliseconds after Jan 1, 1970 00:00:00 UTC. lastEventTimestamp updates on an eventual consistency basis. It typically updates in less than an hour from ingestion, but in rare situations might take longer.
- If the `descending` key column value is true, results are returned in descending order. If the value is to false, results are returned in ascending order. The default value is false.

## Examples

### Basic info

```sql
select
  name,
  log_group_name,
  region
from
  aws_cloudwatch_log_stream;
```

### Count of log streams per log group

```sql
select
  log_group_name,
  count(*) as log_stream_count
from
  aws_cloudwatch_log_stream
group by
  log_group_name;
```
