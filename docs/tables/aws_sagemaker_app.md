# Table: aws_sagemaker_app

Amazon SageMaker app represents an application that supports the reading and execution experience of the user’s notebooks, terminals, and consoles. A user may have multiple Apps active simultaneously.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_app
where 
  user_profile_name = 'testprofile';
```

### List apps that failed to create

```sql
select
  name,
  arn,
  creation_time,
  status,
  failure_reason
from
  aws_sagemaker_app
where 
  user_profile_name = 'testprofile'
  and status = 'Failed';
```