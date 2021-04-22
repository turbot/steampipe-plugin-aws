# Table: aws_sagemaker_notebook_instance

An Amazon SageMaker notebook instance is a machine learning (ML) compute instance running the Jupyter Notebook App.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  instance_type,
  notebook_instance_status
from
  aws_sagemaker_notebook_instance;
```


### List notebook instances that do not have encryption at rest enabled

```sql
select
  name,
  kms_key_id
from
  aws_sagemaker_notebook_instance
where
  kms_key_id is null;
```


### List publicly available notebook instances

```sql
select
  name,
  direct_internet_access
from
  aws_sagemaker_notebook_instance
where
  direct_internet_access = 'Disabled';
```


### List notebook instances that allow root access

```sql
select
  name,
  root_access
from
  aws_sagemaker_notebook_instance
where
  root_access = 'Enabled';
```
