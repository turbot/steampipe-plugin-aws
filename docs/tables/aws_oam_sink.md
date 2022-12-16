# Table: aws_oam_sink

Amazon CloudWatch Observability Access Manager to create and manage links between source accounts and monitoring accounts by using cloudwatch cross-account observability.

## Example

### Basic info

```sql
select
  name,
  id,
  arn,
  tags,
  title
from
  aws_oam_sink;
```

### Get sink by ID

```sql
select
  name,
  id,
  arn
from
  aws_oam_sink
where
  id = 'hfj44c81-7bdf-3847-r7i3-5dfc61b17483';
```