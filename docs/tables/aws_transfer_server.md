# Table: aws_transfer_server

Transfer Family is a fully managed service that enables the transfer of files over the File Transfer Protocol (FTP), File Transfer Protocol over SSL (FTPS), or Secure Shell (SSH) File Transfer Protocol (SFTP) directly into and out of Amazon Simple Storage Service (Amazon S3) or Amazon EFS.

## Examples

### Basic info

```sql
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type
from
  aws_transfer_server;
```
### List servers that are currently OFFLINE

```sql
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type,
  state
from
  aws_transfer_server
where
  state = 'OFFLINE';
```

### Sort servers descending by user count

```sql
select
  server_id,
  user_count
from
  aws_transfer_server
order by
  user_count desc;
```

### List workflows on upload event

```sql
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type,
  workflow_details ->> 'OnUpload' as on_upload_workflow
from
  aws_transfer_server;
```

### List structured destination CloudWatch groups

```sql
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type,
  structured_log_destinations
from
  aws_transfer_server
```
