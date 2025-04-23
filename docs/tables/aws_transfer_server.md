---
title: "Steampipe Table: aws_transfer_server - Query AWS Transfer for SFTP Servers using SQL"
description: "Allows users to query AWS Transfer for SFTP Servers and retrieve detailed information about SFTP servers in their AWS account."
folder: "Transfer Family"
---

# Table: aws_transfer_server - Query AWS Transfer for SFTP Servers using SQL

The AWS Transfer for SFTP service provides a secure way to transfer files into and out of AWS S3 buckets using the Secure Shell (SSH) File Transfer Protocol (SFTP). It integrates with existing authentication systems, and provides DNS routing with Amazon Route 53. This simplifies the migration of file transfer workflows to AWS, while protecting business processes and data.

## Table Usage Guide

The `aws_transfer_server` table in Steampipe provides you with information about SFTP servers within AWS Transfer for SFTP. This table allows you, as a DevOps engineer, to query server-specific details, including server configurations, endpoint details, and associated metadata. You can utilize this table to gather insights on SFTP servers, such as server states, user counts, endpoint types, and more. The schema outlines the various attributes of the SFTP server for you, including the server ID, ARN, endpoint type, logging role, and associated tags.

## Examples

### Basic info
Explore which AWS transfer servers are being used, identifying their associated domains and the types of identity providers and endpoints they utilize. This allows for better management and configuration of data transfer processes.

```sql+postgres
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type
from
  aws_transfer_server;
```

```sql+sqlite
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type
from
  aws_transfer_server;
```
### List servers that are currently OFFLINE
Identify instances where AWS Transfer servers are currently offline. This query can be useful to quickly pinpoint servers that may need attention or troubleshooting.

```sql+postgres
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

```sql+sqlite
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
Analyze your servers to understand which ones are most heavily utilized by users. This aids in resource allocation and capacity planning by highlighting servers with the highest user count.

```sql+postgres
select
  server_id,
  user_count
from
  aws_transfer_server
order by
  user_count desc;
```

```sql+sqlite
select
  server_id,
  user_count
from
  aws_transfer_server
order by
  user_count desc;
```

### List workflows on upload event
Discover the segments that have workflows triggered by an upload event in your AWS Transfer Servers. This can be useful for understanding and managing the actions that take place when data is uploaded to your servers.

```sql+postgres
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type,
  workflow_details ->> 'OnUpload' as on_upload_workflow
from
  aws_transfer_server;
```

```sql+sqlite
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type,
  json_extract(workflow_details, '$.OnUpload') as on_upload_workflow
from
  aws_transfer_server;
```

### List structured destination CloudWatch groups
Explore which AWS Transfer servers are configured to send structured logs to specific destinations. This can help in managing and monitoring file transfers more effectively.

```sql+postgres
select
  server_id,
  domain,
  identity_provider_type,
  endpoint_type,
  structured_log_destinations
from
  aws_transfer_server;
```

```sql+sqlite
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
Explore which servers have specific certificate details and statuses to understand their key algorithms. This is useful in assessing security configurations and ensuring proper server-certificate associations.

```sql+postgres
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

```sql+sqlite
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