---
title: "Steampipe Table: aws_codebuild_report_group - Query AWS CodeBuild Report Groups using SQL"
description: "Allows users to query AWS CodeBuild Report Groups to retrieve detailed information about each report group configuration."
folder: "CodeBuild"
---

# Table: aws_codebuild_report_group - Query AWS CodeBuild Report Groups using SQL

AWS CodeBuild Report Groups are used to collect and store test reports from your build projects. They provide a centralized location for storing test results, code coverage reports, and other build artifacts. Report groups help you track the quality of your code over time and provide insights into your build and test processes.

## Table Usage Guide

The `aws_codebuild_report_group` table in Steampipe provides you with information about report groups within AWS CodeBuild. This table allows you, as a DevOps engineer, to query report group-specific details, including report group ARN, creation date, type, export configuration, and associated metadata. You can utilize this table to gather insights on report groups, such as their configuration, export settings, and associated tags.

## Examples

### Basic info
Explore the features and settings of your AWS CodeBuild report groups to better understand their configuration, such as export settings, type, and regional distribution. This can help in assessing report group performance and operational efficiency.

```sql+postgres
select
  name,
  arn,
  type,
  created,
  last_modified,
  region
from
  aws_codebuild_report_group;
```

```sql+sqlite
select
  name,
  arn,
  type,
  created,
  last_modified,
  region
from
  aws_codebuild_report_group;
```

### Get export configuration details for each report group
Determine the export configuration for each report group to understand how test reports are being stored and exported. This can help in managing and optimizing the reporting process in AWS CodeBuild.

```sql+postgres
select
  name,
  type,
  export_config ->> 'ExportConfigType' as export_config_type,
  export_config -> 'S3Destination' as s3_destination,
  export_config -> 'S3Destination'->>'Bucket' as s3_bucket,
  export_config -> 'S3Destination'->>'EncryptionKey' as encryption_key,
  export_config -> 'S3Destination'->>'EncryptionDisabled' as encryption_disabled,
  export_config -> 'S3Destination'->>'Packaging' as packaging
from
  aws_codebuild_report_group;
```

```sql+sqlite
select
  name,
  type,
  json_extract(export_config, '$.ExportConfigType') as export_config_type,
  json_extract(export_config, '$.S3Destination') as s3_destination,
  json_extract(export_config, '$.S3Destination.Bucket') as s3_bucket,
  json_extract(export_config, '$.S3Destination.EncryptionKey') as encryption_key,
  json_extract(export_config, '$.S3Destination.EncryptionDisabled') as encryption_disabled,
  json_extract(export_config, '$.S3Destination.Packaging') as packaging
from
  aws_codebuild_report_group;
```

### List report groups with specific types
Find report groups with specific types to understand the different kinds of reports being generated and stored in your AWS CodeBuild environment.

```sql+postgres
select
  name,
  arn,
  type,
  created,
  last_modified
from
  aws_codebuild_report_group
where
  type = 'TEST';
```

```sql+sqlite
select
  name,
  arn,
  type,
  created,
  last_modified
from
  aws_codebuild_report_group
where
  type = 'TEST';
```

### Find report groups without encryption
Identify report groups that don't have encryption enabled for their export configuration, which could pose security risks. This is useful for compliance audits and security assessments.

```sql+postgres
select
  name,
  type,
  export_config -> 'S3Destination'->>'EncryptionDisabled' as encryption_disabled
from
  aws_codebuild_report_group
where
  export_config -> 'S3Destination'->>'EncryptionDisabled' = 'true';
```

```sql+sqlite
select
  name,
  type,
  json_extract(export_config, '$.S3Destination.EncryptionDisabled') as encryption_disabled
from
  aws_codebuild_report_group
where
  json_extract(export_config, '$.S3Destination.EncryptionDisabled') = 'true';
```
