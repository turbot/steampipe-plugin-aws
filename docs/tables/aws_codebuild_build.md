# Table: aws_codebuild_project

AWS CodeBuild is a fully managed build service in the cloud. CodeBuild compiles your source code, runs unit tests, and produces artifacts that are ready to deploy. CodeBuild eliminates the need to provision, manage, and scale your own build servers. It provides prepackaged build environments for the most popular programming languages and build tools, such as Apache Maven, Gradle, and more.

## Examples

### Basic info

```sql
select
  arn,
  id,
  build_complete,
  timeout_in_minutes,
  project_name,
  build_status,
  encryption_key,
  end_time,
  region
from
  aws_codebuild_build;
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

### List complete builds

```sql
select
  id,
  arn,
  artifacts,
  build_complete
from
  aws_codebuild_build
where
  build_complete;
```

### List VPC configuration details of builds 

```sql
select
  id,
  arn,
  vpc_config ->> 'SecurityGroupIds' as security_group_id,
  vpc_config ->> 'Subnets' as subnets,
  vpc_config ->> 'VpcId' as vpc_id
from
  aws_codebuild_build;
```

### List artifact details of builds 

```sql
select
  id,
  arn,
  artifacts ->> 'ArtifactIdentifier' as artifact_id,
  artifacts ->> 'BucketOwnerAccess' as bucket_owner_access,
  artifacts ->> 'EncryptionDisabled' as encryption_disabled,
  artifacts ->> 'OverrideArtifactName' as override_artifact_name
from
  aws_codebuild_build;
```
