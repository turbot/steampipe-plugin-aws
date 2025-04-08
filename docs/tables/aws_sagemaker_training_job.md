---
title: "Steampipe Table: aws_sagemaker_training_job - Query AWS SageMaker Training Jobs using SQL"
description: "Allows users to query AWS SageMaker Training Jobs to retrieve information about individual training jobs."
folder: "SageMaker"
---

# Table: aws_sagemaker_training_job - Query AWS SageMaker Training Jobs using SQL

The AWS SageMaker Training Jobs are part of the Amazon SageMaker service, which provides developers and data scientists with the ability to build, train, and deploy machine learning (ML) models quickly. Training jobs in SageMaker are tasks that have a start and end time, in which a specified algorithm is used to train a model with provided data. It offers a flexible, end-to-end solution to handle raw data, feature engineering, training, and model deployment.

## Table Usage Guide

The `aws_sagemaker_training_job` table in Steampipe provides you with information about training jobs within AWS SageMaker. This table allows you, whether you're a data scientist, machine learning engineer, or DevOps engineer, to query job-specific details, including the configuration of the training job, status, performance metrics, and associated metadata. You can utilize this table to monitor the progress of training jobs, verify configuration settings, analyze performance metrics, and more. The schema outlines the various attributes of the training job for you, including the job name, creation time, training time, billable time, and associated tags.

## Examples

### Basic info
Explore which AWS Sagemaker training jobs are active or inactive, along with their respective creation and last modified times. This can be useful for monitoring job status and understanding the timeline of your machine learning workflows.

```sql+postgres
select
  name,
  arn,
  training_job_status,
  creation_time,
  last_modified_time
from
  aws_sagemaker_training_job;
```

```sql+sqlite
select
  name,
  arn,
  training_job_status,
  creation_time,
  last_modified_time
from
  aws_sagemaker_training_job;
```

### Get details of associated ML compute instances and storage volumes for each training job
Explore the configuration of your machine learning compute instances and storage volumes for each training job to better understand the resources being utilized. This can be useful for optimizing costs and resources in your AWS SageMaker training jobs.

```sql+postgres
select
  name,
  arn,
  resource_config ->> 'InstanceType' as instance_type,
  resource_config ->> 'InstanceCount' as instance_count,
  resource_config ->> 'VolumeKmsKeyId' as volume_kms_id,
  resource_config ->> 'VolumeSizeInGB' as volume_size
from
  aws_sagemaker_training_job;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(resource_config, '$.InstanceType') as instance_type,
  json_extract(resource_config, '$.InstanceCount') as instance_count,
  json_extract(resource_config, '$.VolumeKmsKeyId') as volume_kms_id,
  json_extract(resource_config, '$.VolumeSizeInGB') as volume_size
from
  aws_sagemaker_training_job;
```

### List failed training jobs
Identify instances where training jobs have failed in the AWS SageMaker service. This can be useful in troubleshooting and understanding the reasons for failure, thus enabling effective measures to rectify the issues.

```sql+postgres
select
  name,
  arn,
  training_job_status,
  failure_reason
from
  aws_sagemaker_training_job
where
  training_job_status = 'Failed';
```

```sql+sqlite
select
  name,
  arn,
  training_job_status,
  failure_reason
from
  aws_sagemaker_training_job
where
  training_job_status = 'Failed';
```