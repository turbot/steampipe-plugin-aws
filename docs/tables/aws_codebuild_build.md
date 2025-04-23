---
title: "Steampipe Table: aws_codebuild_build - Query AWS CodeBuild Build using SQL"
description: "Allows users to query AWS CodeBuild Build to retrieve information about AWS CodeBuild projects' builds."
folder: "CodeBuild"
---

# Table: aws_codebuild_build - Query AWS CodeBuild Build using SQL

AWS CodeBuild is a fully managed continuous integration service that compiles source code, runs tests, and produces software packages that are ready to deploy. It allows you to build and test code with continuous scaling and enables you to pay only for the build time you use. CodeBuild eliminates the need to provision, manage, and scale your own build servers.

## Table Usage Guide

The `aws_codebuild_build` table in Steampipe provides you with information about builds in AWS CodeBuild. This table allows you as a DevOps engineer to query build-specific details, including build statuses, source details, build environment, and associated metadata. You can utilize this table to gather insights on builds, such as build status, source version, the duration of the build, and more. The schema outlines for you the various attributes of the CodeBuild build, including the build ID, build status, start and end time, and associated tags.

## Examples

### Basic info
Explore which AWS CodeBuild projects have been completed and gain insights into their build status, duration, and other related details. This can help in managing and optimizing the build processes in your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that include encrypted build output artifacts, allowing you to focus on the areas where secure data is being used in your AWS CodeBuild projects.

```sql+postgres
select
  arn,
  id,
  encryption_key
from
  aws_codebuild_build
where
  encryption_key is not null;
```

```sql+sqlite
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
Explore which AWS CodeBuild projects have been fully built. This is useful for assessing project progress and identifying any projects that may still be in progress or have yet to begin.

```sql+postgres
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

```sql+sqlite
select
  id,
  arn,
  artifacts,
  build_complete
from
  aws_codebuild_build
where
  build_complete = 1;
```

### List VPC configuration details of builds
Explore the security aspects of your AWS CodeBuild projects by examining the Virtual Private Cloud (VPC) configurations. This can help you understand and manage the security group IDs, subnets, and VPC IDs associated with your builds.

```sql+postgres
select
  id,
  arn,
  vpc_config ->> 'SecurityGroupIds' as security_group_id,
  vpc_config ->> 'Subnets' as subnets,
  vpc_config ->> 'VpcId' as vpc_id
from
  aws_codebuild_build;
```

```sql+sqlite
select
  id,
  arn,
  json_extract(vpc_config, '$.SecurityGroupIds') as security_group_id,
  json_extract(vpc_config, '$.Subnets') as subnets,
  json_extract(vpc_config, '$.VpcId') as vpc_id
from
  aws_codebuild_build;
```

### List artifact details of builds
This query is useful to gain insights into the specific details of artifacts associated with various builds in AWS CodeBuild. It helps in understanding the access level, encryption status, and other crucial aspects of these artifacts, which can aid in better management and security of your build artifacts.

```sql+postgres
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

```sql+sqlite
select
  id,
  arn,
  json_extract(artifacts, '$.ArtifactIdentifier') as artifact_id,
  json_extract(artifacts, '$.BucketOwnerAccess') as bucket_owner_access,
  json_extract(artifacts, '$.EncryptionDisabled') as encryption_disabled,
  json_extract(artifacts, '$.OverrideArtifactName') as override_artifact_name
from
  aws_codebuild_build;
```

### Get environment details of builds
Explore the specific environmental aspects of your builds in AWS CodeBuild. This can help you understand the settings like compute type, image, and credentials used, which can be useful for troubleshooting or optimizing your build processes.

```sql+postgres
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

```sql+sqlite
select
  id,
  json_extract(environment, '$.Certificate') as environment_certificate,
  json_extract(environment, '$.ComputeType') as environment_compute_type,
  json_extract(environment, '$.EnvironmentVariables') as environment_variables,
  json_extract(environment, '$.Image') as environment_image,
  json_extract(environment, '$.ImagePullCredentialsType') as environment_image_pull_credentials_type,
  json_extract(environment, '$.PrivilegedMode') as environment_privileged_mode,
  json_extract(environment, '$.RegistryCredential') as environment_registry_credential,
  json_extract(environment, '$.Type') as environment_type
from
  aws_codebuild_build;
```

### Get log details of builds
Gain insights into the status and location of your build logs. This query is useful for identifying potential issues with log storage and accessibility, such as encryption status and bucket owner access.

```sql+postgres
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

```sql+sqlite
select
  id,
  json_extract(logs, '$.S3Logs.Status') as s3_log_status,
  json_extract(logs, '$.S3Logs.Location') as s3_log_location,
  json_extract(logs, '$.S3Logs.BucketOwnerAccess') as s3_log_bucket_owner_access,
  json_extract(logs, '$.S3Logs.EncryptionDisabled') as s3_log_encryption_disabled,
  json_extract(logs, '$.DeepLink') as deep_link,
  json_extract(logs, '$.GroupName') as group_name,
  json_extract(logs, '$.S3LogsArn') as s3_logs_arn,
  json_extract(logs, '$.S3DeepLink') as s3_deep_link,
  json_extract(logs, '$.StreamName') as stream_name,
  json_extract(logs, '$.CloudWatchLogsArn') as cloud_watch_logs_arn,
  json_extract(logs, '$
```

### Get network interface details of builds
Explore the network configurations of your AWS CodeBuild projects. This allows you to assess the network interface and subnet details, which can be crucial for understanding your project's networking setup and troubleshooting connectivity issues.

```sql+postgres
select
  id,
  network_interface ->> 'NetworkInterfaceId' as network_interface_id,
  network_interface ->> 'SubnetId' as subnet_id
from
  aws_codebuild_build;
```

```sql+sqlite
select
  id,
  json_extract(network_interface, '$.NetworkInterfaceId') as network_interface_id,
  json_extract(network_interface, '$.SubnetId') as subnet_id
from
  aws_codebuild_build;
```

### List phase details of builds
Explore the progress of your build processes by examining the start and end times, duration, and status of each phase. This can help you identify potential bottlenecks or inefficiencies in your build process.

```sql+postgres
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

```sql+sqlite
select
  aws_codebuild_build.id,
  json_extract(p, '$.EndTime') as end_time,
  json_extract(p, '$.Contexts') as contexts,
  json_extract(p, '$.PhaseType') as phase_type,
  json_extract(p, '$.StartTime') as start_time,
  json_extract(p, '$.DurationInSeconds') as duration_in_seconds,
  json_extract(p, '$.PhaseStatus') as phase_status
from
  aws_codebuild_build,
  json_each(phases) as p;
```

### Get source details of builds
Determine the areas in which the source details of various builds can be analyzed for security and performance. This is beneficial for understanding the build configurations and identifying potential areas of improvement.

```sql+postgres
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

```sql+sqlite
select
  id,
  json_extract(source, '$.Auth') as source_auth,
  json_extract(source, '$.BuildStatusConfig') as source_BuildStatusConfig,
  json_extract(source, '$.Buildspec') as source_buildspec,
  json_extract(source, '$.GitCloneDepth') as source_git_clone_depth,
  json_extract(source, '$.GitSubmodulesConfig') as source_git_submodules_config,
  json_extract(source, '$.GitCloneDepth') as source_git_clone_depth,
  json_extract(source, '$.InsecureSsl') as source_insecure_ssl,
  json_extract(source, '$.Location') as source_location,
  json_extract(source, '$.ReportBuildStatus') as source_report_build_status,
  json_extract(source, '$.SourceIdentifier') as source_identifier,
  json_extract(source, '$.Type') as source_type
from
  aws_codebuild_build;
```

### List file system location details of builds
Explore the specific details of file system locations used in different builds. This can help in understanding the organization of builds and making improvements in the build process.

```sql+postgres
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

```sql+sqlite
select
  aws_codebuild_build.id,
  json_extract(f.value, '$.Identifier') as file_system_identifier,
  json_extract(f.value, '$.Location') as file_system_location,
  json_extract(f.value, '$.MountOptions') as file_system_mount_options,
  json_extract(f.value, '$.MountPoint') as file_system_mount_point,
  json_extract(f.value, '$.Type') as file_system_type
from
  aws_codebuild_build,
  json_each(file_system_locations) as f;
```