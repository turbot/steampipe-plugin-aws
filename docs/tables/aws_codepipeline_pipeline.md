---
title: "Table: aws_codepipeline_pipeline - Query AWS CodePipeline Pipeline using SQL"
description: "Allows users to query AWS CodePipeline Pipeline data, including pipeline names, statuses, stages, and associated metadata."
---

# Table: aws_codepipeline_pipeline - Query AWS CodePipeline Pipeline using SQL

The `aws_codepipeline_pipeline` table in Steampipe provides information about pipelines within AWS CodePipeline. This table allows DevOps engineers to query pipeline-specific details, including pipeline names, statuses, stages, and associated metadata. Users can utilize this table to gather insights on pipelines, such as pipeline execution history, pipeline settings, and more. The schema outlines the various attributes of the pipeline, including the pipeline ARN, creation date, stages, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codepipeline_pipeline` table, you can use the `.inspect aws_codepipeline_pipeline` command in Steampipe.

**Key columns**:

- `name`: This is the name of the pipeline. It is a unique identifier and can be used to join this table with other tables that contain pipeline information.
- `arn`: This is the Amazon Resource Name (ARN) of the pipeline. It is a unique identifier across all AWS resources and can be used to join this table with other tables that contain AWS resource information.
- `created`: This is the date and time the pipeline was created. This information can be used to track the lifecycle of pipelines.

## Examples

### Basic info

```sql
select
  name,
  arn,
  tags_src,
  region,
  account_id
from
  aws_codepipeline_pipeline;
```

### List unencrypted pipelines

```sql
select
  name,
  arn,
  encryption_key
from
  aws_codepipeline_pipeline
where
  encryption_key is null;
```
