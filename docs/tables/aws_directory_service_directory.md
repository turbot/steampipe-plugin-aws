# Table: aws_directory_service_directory

AWS Directory Service for Microsoft Active Directory, also known as AWS Managed Microsoft Active Directory (AD), enables your directory-aware workloads and AWS resources to use managed Active Directory (AD) in AWS.

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