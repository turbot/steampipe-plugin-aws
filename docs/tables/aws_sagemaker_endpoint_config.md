# Table: aws_sagemaker_endpoint_config

The AWS Sagemaker Endpoint configuration that Amazon SageMaker hosting services uses to deploy models. In the configuration, you identify one or more models, created using the CreateModel API, to deploy and the resources that you want Amazon SageMaker to provision.

## Examples

### Basic info

```sql
select
  name,
  arn,
  kms_key_id,
  creation_time,
  production_variants,
  tags
from
  aws_sagemaker_endpoint_config;
```

### List endpoint config which do not have KMS Key ID configuration

```sql
select
  name,
  arn,
  kms_key_id
from
  aws_sagemaker_endpoint_config
where
  kms_key_id is null;
```
