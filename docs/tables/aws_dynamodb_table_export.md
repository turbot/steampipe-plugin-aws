---
title: "Steampipe Table: aws_dynamodb_table_export - Query AWS DynamoDB Table Export using SQL"
description: "Allows users to query AWS DynamoDB Table Exports, providing detailed information on the exports of DynamoDB tables including the export time, status, and the exported data format."
folder: "DynamoDB"
---

# Table: aws_dynamodb_table_export - Query AWS DynamoDB Table Export using SQL

The AWS DynamoDB Table Export is a feature within the AWS DynamoDB service that allows users to export data from their DynamoDB tables into an Amazon S3 bucket. This operation provides a SQL-compatible export of your DynamoDB data, enabling comprehensive data analysis and large scale exports without impacting the performance of your applications. The exported data can be in one of the following formats: Amazon Ion or DynamoDB JSON.

## Table Usage Guide

The `aws_dynamodb_table_export` table in Steampipe provides you with information about the exports of DynamoDB tables within AWS DynamoDB. This table allows you, as a DevOps engineer, to query export-specific details, including the export time, the status of the export, and the format of the exported data. You can utilize this table to gather insights on exports, such as the time of the last export, the status of ongoing exports, and the format of previously exported data. The schema outlines the various attributes of the DynamoDB table export for you, including the export ARN, export time, export status, and the exported data format.

## Examples

### Basic info
Explore the status of your AWS DynamoDB table exports to understand when they ended and their respective formats. This can be useful in managing data exports and ensuring they are successfully stored in the correct S3 bucket.

```sql+postgres
select
  arn,
  end_time,
  export_format,
  export_status,
  s3_bucket
from
  aws_dynamodb_table_export;
```

```sql+sqlite
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
Identify instances where the export process from AWS DynamoDB tables is still ongoing. This is useful to monitor the progress of data exports and ensure they are completing as expected.

```sql+postgres
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

```sql+sqlite
select
  arn,
  end_time,
  export_format,
  export_status,
  s3_bucket
from
  aws_dynamodb_table_export
where
  export_status != 'COMPLETED';
```

### List export details from the last 10 days
Explore the details of your recent AWS DynamoDB table exports to ensure they've been completed successfully and sent to the correct S3 bucket. This is particularly useful for maintaining data integrity and tracking export activities over the past 10 days.

```sql+postgres
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

```sql+sqlite
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
  export_time >= datetime('now', '-10 day');
```