---
title: "Table: aws_cloudtrail_import - Query AWS CloudTrail using SQL"
description: "Allows users to query AWS CloudTrail imports to extract data about imported trail files such as the file name, import time, hash value, and more."
---

# Table: aws_cloudtrail_import - Query AWS CloudTrail using SQL

The `aws_cloudtrail_import` table in Steampipe provides information about imported trail files within AWS CloudTrail. This table allows DevOps engineers to query import-specific details, including the file name, import time, hash value, and more. Users can utilize this table to gather insights on imported trail files, such as their import status, hash type, and hash value. The schema outlines the various attributes of the imported trail file, including the import ID, import time, file name, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudtrail_import` table, you can use the `.inspect aws_cloudtrail_import` command in Steampipe.

**Key columns**:

- `import_id`: This is the unique identifier for each import. It can be used to join this table with other tables to get more detailed information about specific imports.
- `file_name`: This is the name of the imported trail file. It can be used to filter results based on specific file names.
- `import_time`: This is the time when the trail file was imported. It can be used to filter results based on import times.

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