# Table: aws_codebuild_project

AWS CodeBuild is a fully managed build service in the cloud. CodeBuild compiles your source code, runs unit tests, and produces artifacts that are ready to deploy. CodeBuild eliminates the need to provision, manage, and scale your own build servers. It provides prepackaged build environments for the most popular programming languages and build tools, such as Apache Maven, Gradle, and more.

## Examples

### Basic info

```sql
select
  arn,
  id,
  complete,
  build_timeout_in_minutes,
  build_groups,
  batch_build_status,
  encryption_key,
  end_time,
  region
from
  aws_codebuild_build;
```


### List VPC configuration that CodeBuild accesses

```sql
select
  arn,
  id,
  vpc_config
from
  aws_codebuild_build
where
  vpc_config is not null;
```


### List encrypted build output artifacts

```sql
select
  arn,
  id,
  encryption_key
from
  aws_codebuild_build
where
  encryption_key is not null;
```


### List complete build

```sql
select
  id,
  arn,
  artifacts
from
  aws_codebuild_build
where
  complete;
```

### List VPC configuration details of build 

```
select
  id,
  arn,
  vpc_config ->> 'SecurityGroupIds' as security_groups,
  vpc_config ->> 'Subnets' as subnets,
  vpc_config ->> 'VpcId' as vpc_id
from
  aws_codebuild_build;

```