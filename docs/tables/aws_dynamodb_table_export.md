---
title: "Table: aws_dynamodb_table_export - Query AWS DynamoDB Table Export using SQL"
description: "Allows users to query AWS DynamoDB Table Exports, providing detailed information on the exports of DynamoDB tables including the export time, status, and the exported data format."
---

# Table: aws_dynamodb_table_export - Query AWS DynamoDB Table Export using SQL

The `aws_dynamodb_table_export` table in Steampipe provides information about the exports of DynamoDB tables within AWS DynamoDB. This table allows DevOps engineers to query export-specific details, including the export time, the status of the export, and the format of the exported data. Users can utilize this table to gather insights on exports, such as the time of the last export, the status of ongoing exports, and the format of previously exported data. The schema outlines the various attributes of the DynamoDB table export, including the export ARN, export time, export status, and the exported data format.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dynamodb_table_export` table, you can use the `.inspect aws_dynamodb_table_export` command in Steampipe.

### Key columns:

- `table_arn`: The Amazon Resource Number (ARN) of the DynamoDB table. This is a unique identifier that can be used to join this table with other tables to get more detailed information about the DynamoDB table.
- `export_arn`: The Amazon Resource Number (ARN) of the export. This unique identifier can be used to join this table with other tables to get more detailed information about the export.
- `export_status`: The status of the export. This can be used to monitor ongoing exports or to check the status of past exports.

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