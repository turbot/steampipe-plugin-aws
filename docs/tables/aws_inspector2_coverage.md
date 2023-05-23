# Table: aws_inspector2_coverage

AWS Inspector coverage refers to the extent to which the AWS resources in your environment are assessed for security vulnerabilities and compliance with security best practices by the AWS Inspector service.

## Examples

### Basic info

```sql
select
  source_account_id,
  resource_id,
  resource_type,
  scan_type,
  scan_status_reason,
  scan_status_code
from
  aws_inspector2_coverage;
```

### List coverages that are inactive

```sql
select
  source_account_id,
  resource_id,
  resource_type,
  scan_status_code
from
  aws_inspector2_coverage
where
  scan_status_code = 'INACTIVE';
```

### List EC2 resource type coverage

```sql
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  ec2_ami_id,
  ec2_platform
from
  aws_inspector2_coverage
where
  resource_type = 'AWS_EC2_INSTANCE';
```

### List coverages by EC2 instance tags

```sql
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  ec2_ami_id,
  ec2_platform,
  ec2_instance_tags
from
  aws_inspector2_coverage
where
  ec2_inatance_tags = '{"foo": "bar", "foo1": "bar1"}';
```

### List coverages by lambda function tags

```sql
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  lambda_function_name,
  lambda_function_runtime,
  lambda_function_tags
from
  aws_inspector2_coverage
where
  lambda_function_tags = '{"foo": "bar", "foo1": "bar1"}';
```

### List coverage of a package scan

```sql
select
  source_account_id,
  resource_id,
  resource_type,
  scan_type
from
  aws_inspector2_coverage
where
  scan_type = 'PACKAGE';
```

### Get ECR repository details of each coverage

```sql
select
  c.resource_id,
  c.resource_type,
  c.ecr_repository_name,
  r.registry_id,
  r.repository_uri,
  r.encryption_configuration
from
  aws_inspector2_coverage as c,
  aws_ecr_repository as r
where
  r.repository_name = c.ecr_repository_name
and
  c.resource_type = 'AWS_ECR_REPOSITORY';
```