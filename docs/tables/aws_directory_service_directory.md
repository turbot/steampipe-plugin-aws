---
title: "Table: aws_directory_service_directory - Query AWS Directory Service Directories using SQL"
description: "Allows users to query AWS Directory Service Directories for information about AWS Managed Microsoft AD, AWS Managed AD, and Simple AD directories."
---

# Table: aws_directory_service_directory - Query AWS Directory Service Directories using SQL

The `aws_directory_service_directory` table in Steampipe provides information about AWS Directory Service Directories. These include AWS Managed Microsoft AD, AWS Managed AD, and Simple AD directories. This table allows DevOps engineers to query directory-specific details, including directory ID, type, size, and status, among others. Users can utilize this table to gather insights on directories, such as their descriptions, DNS IP addresses, and security group IDs. The schema outlines the various attributes of the Directory Service Directory, including its ARN, creation timestamp, alias, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_directory_service_directory` table, you can use the `.inspect aws_directory_service_directory` command in Steampipe.

**Key columns**:

- `directory_id`: This is the identifier of the directory. It can be used to join this table with other tables that reference the directory ID.
- `arn`: The Amazon Resource Name (ARN) of the directory. This can be used in joining with any table that utilizes ARNs.
- `size`: The size of the directory. This information could be useful when joining with tables that contain resource usage or allocation data.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  directory_id
from
  aws_directory_service_directory;
```

### List MicrosoftAD type directories

```sql
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

```sql
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

### Get snapshot limit details of each directory

```sql
select
  name,
  directory_id,
  snapshot_limit ->> 'ManualSnapshotsCurrentCount' as manual_snapshots_current_count,
  snapshot_limit ->> 'ManualSnapshotsLimit' as manual_snapshots_limit,
  snapshot_limit ->> 'ManualSnapshotsLimitReached' as manual_snapshots_limit_reached
from
  aws_directory_service_directory;
```

### Get SNS topic details of each directory

```sql
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