---
title: "Table: aws_codebuild_project - Query AWS CodeBuild Projects using SQL"
description: "Allows users to query AWS CodeBuild Projects and retrieve comprehensive information about each project."
---

# Table: aws_codebuild_project - Query AWS CodeBuild Projects using SQL

The `aws_codebuild_project` table in Steampipe provides information about projects within AWS CodeBuild. This table allows DevOps engineers to query project-specific details, including project ARN, creation date, project name, service role, and other associated metadata. Users can utilize this table to gather insights on projects, such as the status of each project, the source code repository used, the build environment configuration, and more. The schema outlines the various attributes of the CodeBuild project, including the project ARN, creation date, last modified date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codebuild_project` table, you can use the `.inspect aws_codebuild_project` command in Steampipe.

### Key columns:

- `name`: The name of the CodeBuild project. This can be used to join with other tables that contain CodeBuild project names.
- `arn`: The Amazon Resource Number (ARN) of the CodeBuild project. This unique identifier can be used to join with other tables that contain CodeBuild project ARNs.
- `service_role`: The service role associated with the CodeBuild project. This can be used to join with IAM role tables to get more information about the permissions and policies associated with the project.

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

### List private build projects

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
