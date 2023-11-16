---
title: "Table: aws_sagemaker_training_job - Query AWS SageMaker Training Jobs using SQL"
description: "Allows users to query AWS SageMaker Training Jobs to retrieve information about individual training jobs."
---

# Table: aws_sagemaker_training_job - Query AWS SageMaker Training Jobs using SQL

The `aws_sagemaker_training_job` table in Steampipe provides information about training jobs within AWS SageMaker. This table allows data scientists, machine learning engineers, and DevOps engineers to query job-specific details, including the configuration of the training job, status, performance metrics, and associated metadata. Users can utilize this table to monitor the progress of training jobs, verify configuration settings, analyze performance metrics, and more. The schema outlines the various attributes of the training job, including the job name, creation time, training time, billable time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sagemaker_training_job` table, you can use the `.inspect aws_sagemaker_training_job` command in Steampipe.

**Key columns**:

- `training_job_name`: The name of the training job. This is the primary identifier of the training job and can be used to join with other tables that reference AWS SageMaker training jobs.
- `training_job_arn`: The AWS Resource Name (ARN) of the training job. This unique identifier is important for joining with other AWS tables that use ARN as a reference.
- `creation_time`: The time when the training job was created. This can be useful for tracking the progress and duration of training jobs.

## Examples

### Basic info

```sql
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

```sql
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

### List failed training jobs

```sql
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
