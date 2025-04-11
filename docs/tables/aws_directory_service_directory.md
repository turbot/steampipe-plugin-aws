---
title: "Steampipe Table: aws_directory_service_directory - Query AWS Directory Service Directories using SQL"
description: "Allows users to query AWS Directory Service Directories for information about AWS Managed Microsoft AD, AWS Managed AD, and Simple AD directories."
folder: "Directory Service"
---

# Table: aws_directory_service_directory - Query AWS Directory Service Directories using SQL

The AWS Directory Service provides multiple ways to use Microsoft Active Directory with other AWS services. Directories store information about a network's users, groups, and devices, enabling AWS services and instances to use this information. AWS Directory Service Directories are highly available and scalable, providing a cost-effective way to apply policies and security settings across an AWS environment.

## Table Usage Guide

The `aws_directory_service_directory` table in Steampipe provides you with information about AWS Directory Service Directories. These include AWS Managed Microsoft AD, AWS Managed AD, and Simple AD directories. This table allows you, as a DevOps engineer, to query directory-specific details, including directory ID, type, size, and status, among others. You can utilize this table to gather insights on directories, such as their descriptions, DNS IP addresses, and security group IDs. The schema outlines the various attributes of the Directory Service Directory for you, including its ARN, creation timestamp, alias, and associated tags.

## Examples

### Basic Info
Explore the basic information linked to your AWS Directory Service to better manage and monitor your resources. This can be particularly useful in maintaining security and compliance within your IT infrastructure.

```sql+postgres
select
  name,
  arn,
  directory_id
from
  aws_directory_service_directory;
```

```sql+sqlite
select
  name,
  arn,
  directory_id
from
  aws_directory_service_directory;
```

### List MicrosoftAD type directories
Determine the areas in which MicrosoftAD type directories are being used within your AWS Directory Service. This can help in auditing and managing your AWS resources efficiently.

```sql+postgres
select
  name,
  arn,
  directory_id,
  type
from
  aws_directory_service_directory
where
  type = 'MicrosoftAD';
```

```sql+sqlite
select
  name,
  arn,
  directory_id,
  type
from
  aws_directory_service_directory
where
  type = 'MicrosoftAD';
```

### Get details about the shared directories
Discover the segments that share directories within your network. This query is useful to understand the distribution of shared resources, their status, and the accounts they are shared with, helping you maintain a balanced and secure network.

```sql+postgres
select
  name,
  directory_id,
  sd ->> 'ShareMethod' share_method,
  sd ->> 'ShareStatus' share_status,
  sd ->> 'SharedAccountId' shared_account_id,
  sd ->> 'SharedDirectoryId' shared_directory_id
from
  aws_directory_service_directory,
  jsonb_array_elements(shared_directories) sd;
```

```sql+sqlite
select
  name,
  directory_id,
  json_extract(sd.value, '$.ShareMethod') as share_method,
  json_extract(sd.value, '$.ShareStatus') as share_status,
  json_extract(sd.value, '$.SharedAccountId') as shared_account_id,
  json_extract(sd.value, '$.SharedDirectoryId') as shared_directory_id
from
  aws_directory_service_directory
join
  json_each(shared_directories) as sd;
```

### Get snapshot limit details of each directory
Identify instances where the snapshot limit of each directory in your AWS Directory Service has been reached. This can help manage storage and prevent any potential disruptions due to reaching the limit.

```sql+postgres
select
  name,
  directory_id,
  snapshot_limit ->> 'ManualSnapshotsCurrentCount' as manual_snapshots_current_count,
  snapshot_limit ->> 'ManualSnapshotsLimit' as manual_snapshots_limit,
  snapshot_limit ->> 'ManualSnapshotsLimitReached' as manual_snapshots_limit_reached
from
  aws_directory_service_directory;
```

```sql+sqlite
select
  name,
  directory_id,
  json_extract(snapshot_limit, '$.ManualSnapshotsCurrentCount') as manual_snapshots_current_count,
  json_extract(snapshot_limit, '$.ManualSnapshotsLimit') as manual_snapshots_limit,
  json_extract(snapshot_limit, '$.ManualSnapshotsLimitReached') as manual_snapshots_limit_reached
from
  aws_directory_service_directory;
```

### Get SNS topic details of each directory
Determine the areas in which Simple Notification Service (SNS) topics are linked with each directory in your AWS Directory Service. This can be useful to understand the communication setup and status within your organization's AWS infrastructure.

```sql+postgres
select
  name,
  directory_id,
  e ->> 'CreatedDateTime' as topic_created_date_time,
  e ->> 'Status' as topic_status,
  e ->> 'TopicArn' as topic_arn,
  e ->> 'TopicName' as topic_name
from
  aws_directory_service_directory,
  jsonb_array_elements(event_topics) as e;
```

```sql+sqlite
select
  name,
  directory_id,
  json_extract(e.value, '$.CreatedDateTime') as topic_created_date_time,
  json_extract(e.value, '$.Status') as topic_status,
  json_extract(e.value, '$.TopicArn') as topic_arn,
  json_extract(e.value, '$.TopicName') as topic_name
from
  aws_directory_service_directory
join
  json_each(event_topics) as e;
```