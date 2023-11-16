---
title: "Table: aws_fms_app_list - Query AWS Firewall Manager Applications using SQL"
description: "Allows users to query AWS Firewall Manager Applications to obtain comprehensive details about each application, including application ID, protocol, source and destination IP ranges, and source and destination ports."
---

# Table: aws_fms_app_list - Query AWS Firewall Manager Applications using SQL

The `aws_fms_app_list` table in Steampipe provides information about applications within AWS Firewall Manager (FMS). This table allows DevOps engineers to query application-specific details, including application ID, protocol, source and destination IP ranges, and source and destination ports. Users can utilize this table to gather insights on applications, such as their associated protocols, IP ranges, and ports. The schema outlines the various attributes of the application, including the application ARN, creation date, attached policies, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_fms_app_list` table, you can use the `.inspect aws_fms_app_list` command in Steampipe.

### Key columns:

- `app_id`: The unique identifier for an application. This can be used to join this table with other tables to get more detailed information about the application.
- `protocol`: The protocol used by the application. This can provide insights into the type of traffic the application is handling.
- `source_ip`: The source IP range for the application. This can be used to identify the origin of traffic for the application.


## Examples

### Basic info

```sql
select
  list_name,
  list_id,
  arn,
  create_time,
  creation_time
from
  aws_fms_app_list;
```

### List of apps created in last 30 days

```sql
select
  list_name,
  list_id,
  arn,
  create_time,
  creation_time
from
  aws_fms_app_list
where
  create_time >= now() - interval '30' day;
```

### Get application details of each app list

```sql
select
  list_name,
  list_id,
  a ->> 'AppName' as app_name,
  a ->> 'Port' as port,
  a ->> 'Protocol' as protocol
from
  aws_fms_app_list,
  jsonb_array_elements(apps_list -> 'AppsList') as a;

```