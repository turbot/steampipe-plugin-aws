---
title: "Steampipe Table: aws_codebuild_project - Query AWS CodeBuild Projects using SQL"
description: "Allows users to query AWS CodeBuild Projects and retrieve comprehensive information about each project."
folder: "CodeBuild"
---

# Table: aws_codebuild_project - Query AWS CodeBuild Projects using SQL

The AWS CodeBuild Project is a component of AWS CodeBuild, a fully managed continuous integration service that compiles source code, runs tests, and produces software packages that are ready to deploy. With CodeBuild, you don’t need to provision, manage, and scale your own build servers. It provides prepackaged build environments for popular programming languages and build tools, such as Apache Maven, Gradle, and more.

## Table Usage Guide

The `aws_codebuild_project` table in Steampipe provides you with information about projects within AWS CodeBuild. This table allows you, as a DevOps engineer, to query project-specific details, including project ARN, creation date, project name, service role, and other associated metadata. You can utilize this table to gather insights on projects, such as the status of each project, the source code repository used, the build environment configuration, and more. The schema outlines the various attributes of the CodeBuild project for you, including the project ARN, creation date, last modified date, and associated tags.

## Examples

### Basic info
Explore the features and settings of your AWS CodeBuild projects to better understand their configuration, such as encryption details, build limits, and regional distribution. This can help in assessing project performance, security, and operational efficiency.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which each project's build input details are configured, such as authorization, build status, and source location. This can help in managing and troubleshooting the build process in AWS CodeBuild projects.

```sql+postgres
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

```sql+sqlite
select
  name,
  source_version,
  json_extract(source, '$.Auth') as auth,
  json_extract(source, '$.BuildStatusConfig') as build_status_config,
  json_extract(source, '$.Buildspec') as build_spec,
  json_extract(source, '$.GitCloneDepth') as git_clone_depth,
  json_extract(source, '$.GitSubmodulesConfig') as git_submodules_config,
  json_extract(source, '$.InsecureSsl') as insecure_ssl,
  json_extract(source, '$.Location') as location,
  json_extract(source, '$.ReportBuildStatus') as report_build_status,
  json_extract(source, '$.SourceIdentifier') as source_identifier,
  json_extract(source, '$.Type') as type
from
  aws_codebuild_project;
```


### List projects which are not created within a VPC
Determine the areas in which AWS CodeBuild projects have been created without a Virtual Private Cloud (VPC) configuration. This is useful for identifying potential security risks and ensuring all projects follow best practices for network security.

```sql+postgres
select
  name,
  description,
  vpc_config
from
  aws_codebuild_project
where
  vpc_config is null;
```

```sql+sqlite
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
Identify projects that have disabled logging, allowing you to pinpoint areas where crucial data might not be being recorded for future analysis. This is particularly useful for maintaining project transparency and troubleshooting potential issues.

```sql+postgres
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

```sql+sqlite
select
  name,
  description,
  json_extract(logs_config, '$.CloudWatchLogs.Status') as cloud_watch_logs_status,
  json_extract(logs_config, '$.S3Logs.Status') as s3_logs_status
from
  aws_codebuild_project
where
  json_extract(logs_config, '$.CloudWatchLogs.Status') = 'DISABLED'
  and json_extract(logs_config, '$.S3Logs.Status') = 'DISABLED';
```

### List private build projects
Determine the areas in which your AWS CodeBuild projects are set to private, allowing you to gain insights into your project visibility settings and understand where potential privacy concerns may arise.

```sql+postgres
select
  name,
  arn,
  project_visibility
from
  aws_codebuild_project
where
  project_visibility = 'PRIVATE';
```

```sql+sqlite
select
  name,
  arn,
  project_visibility
from
  aws_codebuild_project
where
  project_visibility = 'PRIVATE';
```