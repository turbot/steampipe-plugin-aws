# Table: aws_cloudtrail_import

AWS CloudTrail import of logged trail events from a source S3 bucket to a destination event data store. By default, CloudTrail only imports events contained in the S3 bucket's CloudTrail prefix and the prefixes inside the CloudTrail prefix, and does not check prefixes for other AWS services.

## Examples

### Basic info

```sql
select
  import_id,
  created_timestamp,
  import_status,
  destinations
from
  aws_cloudtrail_import;
```

### List imports that are not completed

```sql
select
  import_id,
  created_timestamp,
  import_source
from
  aws_cloudtrail_import
where
  import_status <> 'COMPLETED';
```

### List imports that are created in last 30 days

```sql
select
  import_id,
  created_timestamp,
  import_status,
  start_event_time,
  end_event_time
from
  aws_cloudtrail_import
where
  created_timestamp >= now() - interval '30' day;
```

### Get import source details of each import

```sql
select
  import_id,
  import_status,
  import_source ->> 'S3BucketAccessRoleArn' as s3_bucket_access_role_arn,
  import_source ->> 'S3BucketRegion' as s3_bucket_region,
  import_source ->> 'S3LocationUri' as s3_location_uri

from
  aws_cloudtrail_import;
```

### Get import statistic of each import

```sql
select
  import_id,
  import_status,
  import_statistics -> 'EventsCompleted' as events_completed,
  import_statistics -> 'FailedEntries' as failed_entries,
  import_statistics -> 'FilesCompleted' as files_completed,
  import_statistics -> 'FilesCompleted' as prefixes_completed,
  import_statistics -> 'PrefixesFound' as PrefixesFound
from
  aws_cloudtrail_import;
```