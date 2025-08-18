---
title: "Steampipe Table: aws_transfer_connector - Query AWS Transfer Connectors using SQL"
description: "Allows users to query AWS Transfer Connectors and retrieve detailed information about AS2 and SFTP connectors in their AWS account."
folder: "Transfer Family"
---

# Table: aws_transfer_connector - Query AWS Transfer Connectors using SQL

AWS Transfer Connector is a resource that enables you to send files using either the AS2 or SFTP protocol. Connectors are used to transfer files between your AWS Transfer Family server and external partners or systems.

## Table Usage Guide

The `aws_transfer_connector` table in Steampipe provides you with information about AS2 and SFTP connectors within AWS Transfer Family. This table allows you, as a DevOps engineer, to query connector-specific details, including connector configurations, endpoint details, and associated metadata. You can utilize this table to gather insights on connectors, such as connector states, protocol configurations, egress IP addresses, and more. The schema outlines the various attributes of the connector for you, including the connector ID, ARN, URL, access role, and associated tags.

## Examples

### Basic info
Explore which AWS transfer connectors are being used, identifying their associated URLs and the types of access roles and security policies they utilize. This allows for better management and configuration of file transfer processes.

```sql+postgres
select
  connector_id,
  arn,
  url,
  access_role,
  logging_role,
  security_policy_name
from
  aws_transfer_connector;
```

```sql+sqlite
select
  connector_id,
  arn,
  url,
  access_role,
  logging_role,
  security_policy_name
from
  aws_transfer_connector;
```

### List connectors with AS2 configuration
Identify connectors that are configured for AS2 protocol, which is commonly used for secure business-to-business file transfers. This helps in understanding which connectors are set up for AS2-based file exchanges.

```sql+postgres
select
  connector_id,
  arn,
  url,
  as2_config
from
  aws_transfer_connector
where
  as2_config is not null;
```

```sql+sqlite
select
  connector_id,
  arn,
  url,
  as2_config
from
  aws_transfer_connector
where
  as2_config is not null;
```

### List connectors with SFTP configuration
Identify connectors that are configured for SFTP protocol, which is commonly used for secure file transfers over SSH. This helps in understanding which connectors are set up for SFTP-based file exchanges.

```sql+postgres
select
  connector_id,
  arn,
  url,
  sftp_config
from
  aws_transfer_connector
where
  sftp_config is not null;
```

```sql+sqlite
select
  connector_id,
  arn,
  url,
  sftp_config
from
  aws_transfer_connector
where
  sftp_config is not null;
```

### List connectors with their egress IP addresses
Explore which connectors have assigned egress IP addresses for outbound connections. This information is useful for network security and firewall configuration.

```sql+postgres
select
  connector_id,
  arn,
  service_managed_egress_ip_addresses
from
  aws_transfer_connector
where
  service_managed_egress_ip_addresses is not null;
```

```sql+sqlite
select
  connector_id,
  arn,
  service_managed_egress_ip_addresses
from
  aws_transfer_connector
where
  service_managed_egress_ip_addresses is not null;
```

### Find connectors by tag
Discover connectors that have specific tags assigned, such as environment tags. This helps in organizing and filtering connectors based on their purpose or environment.

```sql+postgres
select
  connector_id,
  arn,
  tags
from
  aws_transfer_connector
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  connector_id,
  arn,
  tags
from
  aws_transfer_connector
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List connectors with specific security policy
Identify connectors that are using a particular security policy. This helps in ensuring consistent security configurations across your connectors.

```sql+postgres
select
  connector_id,
  arn,
  security_policy_name
from
  aws_transfer_connector
where
  security_policy_name = 'TransferSecurityPolicy-2020-06';
```

```sql+sqlite
select
  connector_id,
  arn,
  security_policy_name
from
  aws_transfer_connector
where
  security_policy_name = 'TransferSecurityPolicy-2020-06';
```

### List connectors with access role information
Explore connectors that have access roles configured, which are essential for proper authentication and authorization. This helps in understanding the IAM setup for your connectors.

```sql+postgres
select
  connector_id,
  arn,
  access_role,
  logging_role
from
  aws_transfer_connector
where
  access_role is not null;
```

```sql+sqlite
select
  connector_id,
  arn,
  access_role,
  logging_role
from
  aws_transfer_connector
where
  access_role is not null;
```

### List connectors with their full configuration details
Get a comprehensive view of all connector configurations including their protocols, roles, security policies, and network settings. This is useful for auditing and compliance purposes.

```sql+postgres
select
  connector_id,
  arn,
  url,
  access_role,
  logging_role,
  security_policy_name,
  service_managed_egress_ip_addresses,
  as2_config,
  sftp_config,
  tags
from
  aws_transfer_connector;
```

```sql+sqlite
select
  connector_id,
  arn,
  url,
  access_role,
  logging_role,
  security_policy_name,
  service_managed_egress_ip_addresses,
  as2_config,
  sftp_config,
  tags
from
  aws_transfer_connector;
```

### Find connectors by URL pattern
Search for connectors based on their endpoint URL patterns. This helps in identifying connectors that use specific protocols or partner endpoints.

```sql+postgres
select
  connector_id,
  arn,
  url
from
  aws_transfer_connector
where
  url like '%sftp%';
```

```sql+sqlite
select
  connector_id,
  arn,
  url
from
  aws_transfer_connector
where
  url like '%sftp%';
```