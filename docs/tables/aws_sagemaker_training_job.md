# Table: aws_sagemaker_training_job

A Training job helps to train a model in Amazon SageMaker.

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

### Get details of ML compute instances and ML storage volumes associated with training job

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
