# Table: aws_sagemaker_model

The AWS SageMaker Model resource to create a model to host at an Amazon SageMaker endpoint.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model;
```

### List models which are isolated

```sql
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model
where
  enable_network_isolation;
```
