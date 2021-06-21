# Table: aws_sagemaker_endpoint_configuration

An endpoint configuration resource is used by Amazon SageMaker hosting services to deploy models. Each configuration defines one or more models, along with other resources, that Sagemaker will provision.

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
  aws_sagemaker_endpoint_configuration;
```

### List unencrypted endpoint configurations

```sql
select
  name,
  arn,
  kms_key_id
from
  aws_sagemaker_endpoint_configuration
where
  kms_key_id is null;
```
