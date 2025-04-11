---
title: "Steampipe Table: aws_codepipeline_pipeline - Query AWS CodePipeline Pipeline using SQL"
description: "Allows users to query AWS CodePipeline Pipeline data, including pipeline names, statuses, stages, and associated metadata."
folder: "CodePipeline"
---

# Table: aws_codepipeline_pipeline - Query AWS CodePipeline Pipeline using SQL

The AWS CodePipeline is a fully managed continuous delivery service that helps you automate your release pipelines for fast and reliable application and infrastructure updates. CodePipeline automates the build, test, and deploy phases of your release process every time there is a code change, based on the release model you define. This enables you to rapidly and reliably deliver features and updates.

## Table Usage Guide

The `aws_codepipeline_pipeline` table in Steampipe provides you with information about pipelines within AWS CodePipeline. This table allows you, as a DevOps engineer, to query pipeline-specific details, including pipeline names, statuses, stages, and associated metadata. You can utilize this table to gather insights on pipelines, such as pipeline execution history, pipeline settings, and more. The schema outlines the various attributes of the pipeline for you, including the pipeline ARN, creation date, stages, and associated tags.

## Examples

### Basic info
Discover the segments that are part of the AWS CodePipeline service. This information can be useful for auditing, tracking resource usage, and understanding your overall AWS environment.

```sql+postgres
select
  name,
  arn,
  tags_src,
  region,
  account_id
from
  aws_codepipeline_pipeline;
```

```sql+sqlite
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
Discover the segments that have unencrypted pipelines in the AWS CodePipeline service to enhance your security measures. This helps in identifying potential security risks and taking necessary actions to protect your data.

```sql+postgres
select
  name,
  arn,
  encryption_key
from
  aws_codepipeline_pipeline
where
  encryption_key is null;
```

```sql+sqlite
select
  name,
  arn,
  encryption_key
from
  aws_codepipeline_pipeline
where
  encryption_key is null;
```