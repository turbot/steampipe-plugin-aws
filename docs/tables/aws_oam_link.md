# Table: aws_oam_link

Amazon CloudWatch Observability Access Manager (OAM) link shares the observability data with the monitoring account. The shared observability data can include metrics in Amazon CloudWatch, logs in Amazon CloudWatch Logs, and traces in AWS X-Ray.

## Example

### Basic info

```sql
select
  id,
  arn,
  sink_arn,
  label,
  resource_types
from
  aws_oam_link;
```

### Get sink details of each link

```sql
select
  l.id,
  l.arn,
  s.name as sink_name,
  l.sink_arn
from
  aws_oam_link as l,
  aws_oam_sink as s;
```

### List links that share data of CloudWatch log group resource type

```sql
select
  id,
  arn,
  label,
  label_template,
  r as resource_type
from
  aws_oam_link,
  jsonb_array_elements_text(resource_types) as r
where
  r = 'AWS::Logs::LogGroup';
```
