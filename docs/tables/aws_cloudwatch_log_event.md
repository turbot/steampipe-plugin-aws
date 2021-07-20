# Table: aws_cloudwatch_log_event

This table list events from a Cloudwatch log group.

**Important Notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
- This table supports optional quals. Queries with optional quals are optimised to used aws cloudwatch filters. Optional quals is supported for below columns:
  - `log_stream_name`
  - `filter`
  - `region`
  - `timestamp`

**Similar Tables**

- [aws_vpc_flow_log_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log_event) - To get details of VPC flow log events from Cloudwatch logs.
- [aws_cloudtrail_trail_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_trail_event) - To get details of CloudTrail events from Cloudwatch logs.

## Examples

### Basic info

```sql
select
  log_group_name,
  log_stream_name,
  event_id,
  timestamp,
  ingestion_time,
  message
from
  aws_cloudwatch_log_event
where
  log_group_name = '/aws/lambda/myfunction';
```
