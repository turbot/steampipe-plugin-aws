---
title: "Table: aws_transfer_server - Query AWS Transfer for SFTP Servers using SQL"
description: "Allows users to query AWS Transfer for SFTP Servers and retrieve detailed information about SFTP servers in their AWS account."
---

# Table: aws_transfer_server - Query AWS Transfer for SFTP Servers using SQL

The `aws_transfer_server` table in Steampipe provides information about SFTP servers within AWS Transfer for SFTP. This table allows DevOps engineers to query server-specific details, including server configurations, endpoint details, and associated metadata. Users can utilize this table to gather insights on SFTP servers, such as server states, user counts, endpoint types, and more. The schema outlines the various attributes of the SFTP server, including the server ID, ARN, endpoint type, logging role, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_transfer_server` table, you can use the `.inspect aws_transfer_server` command in Steampipe.

**Key columns**:

- `server_id`: The unique ID of the SFTP server. This is a primary identifier that can be used to join this table with other tables.
- `arn`: The Amazon Resource Number (ARN) of the SFTP server. This is a globally unique identifier that can be used for joining with other AWS resource tables.
- `endpoint_type`: The endpoint type of the SFTP server (PUBLIC or VPC_ENDPOINT). This column can provide insights into the network accessibility of the SFTP server.

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
  aws_transfer_server;
```

### Get certificate details for servers

```sql
select
  s.server_id,
  c.certificate_arn,
  c.status as certificate_status,
  c.key_algorithm
from
  aws_transfer_server as s,
  aws_acm_certificate as c
where
  s.certificate = c.certificate_arn;
```