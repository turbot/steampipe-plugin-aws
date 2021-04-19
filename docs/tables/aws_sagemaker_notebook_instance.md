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


### List notebook instances with unecrypted data on the ML storage volume

```sql
select
  name,
  kms_key_id
from
  aws_sagemaker_notebook_instance
where
  kms_key_id is null;
```


### List notebook instances whose data encrypted with kms customer master keys

```sql
select
  name,
  kms_key_id
from
  aws_sagemaker_notebook_instance
where
  kms_key_id is null;
```


### List notebook instances that are publicly available

```sql
select
  name,
  direct_internet_access
from
  aws_sagemaker_notebook_instance
where
  direct_internet_access = 'Disabled';
```


### List notebook instances which provides root access to the users

```sql
select
  name,
  root_access
from
  aws_sagemaker_notebook_instance
where
  root_access = 'Enabled';
```