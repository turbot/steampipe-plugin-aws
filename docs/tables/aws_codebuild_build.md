---
title: "Table: aws_codebuild_build - Query AWS CodeBuild Build using SQL"
description: "Allows users to query AWS CodeBuild Build to retrieve information about AWS CodeBuild projects' builds."
---

# Table: aws_codebuild_build - Query AWS CodeBuild Build using SQL

The `aws_codebuild_build` table in Steampipe provides information about builds in AWS CodeBuild. This table allows DevOps engineers to query build-specific details, including build statuses, source details, build environment, and associated metadata. Users can utilize this table to gather insights on builds, such as build status, source version, the duration of the build, and more. The schema outlines the various attributes of the CodeBuild build, including the build ID, build status, start and end time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codebuild_build` table, you can use the `.inspect aws_codebuild_build` command in Steampipe.

Key columns:

- `id`: The unique identifier for the build. This can be used to join this table with other tables for more detailed information on a specific build.
- `project_name`: The name of the AWS CodeBuild project. This is useful for querying all builds related to a specific project.
- `status`: The current status of the build. This can be helpful in tracking the progress of builds and identifying any issues.

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

### Get environment details of builds

```sql
select
  id,
  environment ->> 'Certificate' as environment_certificate,
  environment ->> 'ComputeType' as environment_compute_type,
  environment ->> 'EnvironmentVariables' as environment_variables,
  environment ->> 'Image' as environment_image,
  environment ->> 'ImagePullCredentialsType' as environment_image_pull_credentials_type,
  environment ->> 'PrivilegedMode' as environment_privileged_mode,
  environment ->> 'RegistryCredential' as environment_registry_credential,
  environment ->> 'Type' as environment_type
from
  aws_codebuild_build;
```

### Get log details of builds

```sql
select
  id,
  logs -> 'S3Logs' ->> 'Status' as s3_log_status,
  logs -> 'S3Logs' ->> 'Location' as s3_log_location,
  logs -> 'S3Logs' ->> 'BucketOwnerAccess' as s3_log_bucket_owner_access,
  logs -> 'S3Logs' ->> 'EncryptionDisabled' as s3_log_encryption_disabled,
  logs ->> 'DeepLink' as deep_link,
  logs ->> 'GroupName' as group_name,
  logs ->> 'S3LogsArn' as s3_logs_arn,
  logs ->> 'S3DeepLink' as s3_deep_link,
  logs ->> 'StreamName' as stream_name,
  logs ->> 'CloudWatchLogsArn' as cloud_watch_logs_arn,
  logs -> 'CloudWatchLogs' ->> 'Status' as cloud_watch_logs_status,
  logs -> 'CloudWatchLogs' ->> 'GroupName' as cloud_watch_logs_group_name,
  logs -> 'CloudWatchLogs' ->> 'StreamName' as cloud_watch_logs_stream_name
from
  aws_codebuild_build;
```

### Get network interface details of builds

```sql
select
  id,
  network_interfaces ->> 'NetworkInterfaceId' as network_interface_id,
  network_interfaces ->> 'SubnetId' as subnet_id,
from
  aws_codebuild_build;
```

### List phase details of builds

```sql
select
  id,
  p ->> 'EndTime' as end_time,
  p ->> 'Contexts' as contexts,
  p ->> 'PhaseType' as phase_type,
  p ->> 'StartTime' as start_time,
  p ->> 'DurationInSeconds' as duration_in_seconds,
  p ->> 'PhaseStatus' as phase_status
from
  aws_codebuild_build,
  jsonb_array_elements(phases) as p;
```

### Get source details of builds

```sql
select
  id,
  source ->> 'Auth' as source_auth,
  source ->> 'BuildStatusConfig' as source_BuildStatusConfig,
  source ->> 'Buildspec' as source_buildspec,
  source ->> 'GitCloneDepth' as source_git_clone_depth,
  source ->> 'GitSubmodulesConfig' as source_git_submodules_config,
  source ->> 'GitCloneDepth' as source_git_clone_depth,
  source ->> 'InsecureSsl' as source_insecure_ssl,
  source ->> 'Location' as source_location,
  source ->> 'ReportBuildStatus' as source_report_build_status,
  source ->> 'SourceIdentifier' as source_identifier,
  source ->> 'Type' as source_type
from
  aws_codebuild_build;
```

### List file system location details of builds

```sql
select
  id,
  f ->> 'Identifier' as file_system_identifier,
  f ->> 'Location' as file_system_location,
  f ->> 'MountOptions' as file_system_mount_options,
  f ->> 'MountPoint' as file_system_mount_point,
  f ->> 'Type' as file_system_type
from
  aws_codebuild_build,
  jsonb_array_elements(file_system_locations) as f;
```