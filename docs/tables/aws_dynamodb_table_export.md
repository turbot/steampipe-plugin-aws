# Table: aws_dynamodb_table_export

Using DynamoDB table export, you can export data from an Amazon DynamoDB table from any time within your point-in-time recovery window to an Amazon S3 bucket.

## Examples

### Basic info

```sql
select
  arn,
  end_time,
  export_format,
  export_status,
  s3_bucket
from
  aws_dynamodb_table_export;
```

### List exports that are not completed

```sql
select
  arn,
  end_time,
  export_format,
  export_status,
  s3_bucket
from
  aws_dynamodb_table_export
where
  export_status <> 'COMPLETED';
```

### List export details from the last 10 days

```sql
select
  arn,
  end_time,
  export_format,
  export_status,
  export_time,
  s3_bucket
from
  aws_dynamodb_table_export
where
  export_time >= now() - interval '10' day;
```