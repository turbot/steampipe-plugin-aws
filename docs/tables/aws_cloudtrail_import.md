---
title: "Steampipe Table: aws_cloudtrail_import - Query AWS CloudTrail using SQL"
description: "Allows users to query AWS CloudTrail imports to extract data about imported trail files such as the file name, import time, hash value, and more."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_import - Query AWS CloudTrail using SQL

AWS CloudTrail is a service that enables governance, compliance, operational auditing, and risk auditing of your AWS account. It allows you to log, continuously monitor, and retain account activity related to actions across your AWS infrastructure. CloudTrail provides event history of your AWS account activity, including actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services.

## Table Usage Guide

The `aws_cloudtrail_import` table in Steampipe provides you with information about imported trail files within AWS CloudTrail. This table allows you, as a DevOps engineer, to query import-specific details, including the file name, import time, hash value, and more. You can utilize this table to gather insights on imported trail files, such as their import status, hash type, and hash value. The schema outlines the various attributes of the imported trail file for you, including the import ID, import time, file name, and associated metadata.

## Examples

### Basic info
Explore which AWS CloudTrail imports have been created and their current status to understand where the data is being sent. This can help in assessing the flow of data and ensuring it is reaching the intended destinations.

```sql+postgres
select
  import_id,
  created_timestamp,
  import_status,
  destinations
from
  aws_cloudtrail_import;
```

```sql+sqlite
select
  import_id,
  created_timestamp,
  import_status,
  destinations
from
  aws_cloudtrail_import;
```

### List imports that are not completed
Identify instances where CloudTrail imports are still in progress. This is useful for tracking the progress of data import tasks and identifying any potential issues or delays.

```sql+postgres
select
  import_id,
  created_timestamp,
  import_source
from
  aws_cloudtrail_import
where
  import_status <> 'COMPLETED';
```

```sql+sqlite
select
  import_id,
  created_timestamp,
  import_source
from
  aws_cloudtrail_import
where
  import_status != 'COMPLETED';
```

### List imports that are created in last 30 days
Identify recent imports within the last 30 days to track their status and duration. This is useful for understanding recent activity and ensuring timely data retrieval.

```sql+postgres
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

```sql+sqlite
select
  import_id,
  created_timestamp,
  import_status,
  start_event_time,
  end_event_time
from
  aws_cloudtrail_import
where
  created_timestamp >= datetime('now', '-30 day');
```

### Get import source details of each import
Identify the origins of each import by examining the access role, region, and URI of the S3 bucket used. This can be useful for auditing purposes or to troubleshoot issues related to specific imports.

```sql+postgres
select
  import_id,
  import_status,
  import_source ->> 'S3BucketAccessRoleArn' as s3_bucket_access_role_arn,
  import_source ->> 'S3BucketRegion' as s3_bucket_region,
  import_source ->> 'S3LocationUri' as s3_location_uri

from
  aws_cloudtrail_import;
```

```sql+sqlite
select
  import_id,
  import_status,
  json_extract(import_source, '$.S3BucketAccessRoleArn') as s3_bucket_access_role_arn,
  json_extract(import_source, '$.S3BucketRegion') as s3_bucket_region,
  json_extract(import_source, '$.S3LocationUri') as s3_location_uri

from
  aws_cloudtrail_import;
```

### Get import statistic of each import
Gain insights into the performance of each import operation by assessing the number of completed events, failed entries, and completed files. This is useful for monitoring the efficiency and reliability of data import processes.

```sql+postgres
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

```sql+sqlite
select
  import_id,
  import_status,
  json_extract(import_statistics, '$.EventsCompleted') as events_completed,
  json_extract(import_statistics, '$.FailedEntries') as failed_entries,
  json_extract(import_statistics, '$.FilesCompleted') as files_completed,
  json_extract(import_statistics, '$.FilesCompleted') as prefixes_completed,
  json_extract(import_statistics, '$.PrefixesFound') as PrefixesFound
from
  aws_cloudtrail_import;
```