# Table: aws_codebuild_project

An AWS CodeBuild project configures how CodeBuild builds your source code. For example, it tells CodeBuild where to get the source code and which build environment to use.

## Examples

### Basic info

```sql
select
  name,
  description,
  encryption_key,
  concurrent_build_limit,
  source_version,
  service_role,
  created,
  last_modified,
  region
from
  aws_codebuild_project;
```


### Get the build input details for each project

```sql
select
  name,
  source_version,
  source ->> 'Auth' as auth,
  source ->> 'BuildStatusConfig' as build_status_config,
  source ->> 'Buildspec' as build_spec,
  source ->> 'GitCloneDepth' as git_clone_depth,
  source ->> 'GitSubmodulesConfig' as git_submodules_config,
  source ->> 'InsecureSsl' as insecure_ssl,
  source ->> 'Location' as location,
  source ->> 'ReportBuildStatus' as report_build_status,
  source ->> 'SourceIdentifier' as source_identifier,
  source ->> 'Type' as type
from
  aws_codebuild_project;
```


### List projects which are not created within a VPC

```sql
select
  name,
  description,
  vpc_config
from
  aws_codebuild_project
where
  vpc_config is null;
```


### List projects that do not have logging enabled

```sql
select
  name,
  description,
  logs_config -> 'CloudWatchLogs' ->> 'Status' as cloud_watch_logs_status,
  logs_config -> 'S3Logs' ->> 'Status' as s3_logs_status
from
  aws_codebuild_project
where
  logs_config -> 'CloudWatchLogs' ->> 'Status' = 'DISABLED'
  and logs_config -> 'S3Logs' ->> 'Status' = 'DISABLED';
```

### List projects which are private

```sql
select
  name,
  arn,
  project_visibility
from
  aws_codebuild_project
where
  project_visibility = 'PRIVATE';
```
