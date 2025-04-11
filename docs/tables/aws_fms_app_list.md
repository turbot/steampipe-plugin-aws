---
title: "Steampipe Table: aws_fms_app_list - Query AWS Firewall Manager Applications using SQL"
description: "Allows users to query AWS Firewall Manager Applications to obtain comprehensive details about each application, including application ID, protocol, source and destination IP ranges, and source and destination ports."
folder: "Firewall Manager (FMS)"
---

# Table: aws_fms_app_list - Query AWS Firewall Manager Applications using SQL

The AWS Firewall Manager Applications is a resource that enables you to centrally manage and deploy AWS WAF rules across your AWS accounts and applications. This service allows you to set up firewall rules once and apply them across your entire infrastructure. It helps to maintain a consistent security posture, even as new resources and accounts are added.

## Table Usage Guide

The `aws_fms_app_list` table in Steampipe provides you with information about applications within AWS Firewall Manager (FMS). This table allows you, as a DevOps engineer, to query application-specific details, including application ID, protocol, source and destination IP ranges, and source and destination ports. You can utilize this table to gather insights on applications, such as their associated protocols, IP ranges, and ports. The schema outlines the various attributes of the application for you, including the application ARN, creation date, attached policies, and associated tags.

## Examples

### Basic info
Explore the creation times of various applications in AWS Firewall Manager to understand their longevity and potential security implications. This can be particularly useful for auditing purposes or identifying outdated applications that may pose a security risk.

```sql+postgres
select
  list_name,
  list_id,
  arn,
  create_time
from
  aws_fms_app_list;
```

```sql+sqlite
select
  list_name,
  list_id,
  arn,
  create_time
from
  aws_fms_app_list;
```

### List of apps created in last 30 days
Discover the segments that have newly created apps within the past month. This can help in tracking recent additions and managing app inventory effectively.

```sql+postgres
select
  list_name,
  list_id,
  arn,
  create_time
from
  aws_fms_app_list
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  list_name,
  list_id,
  arn,
  create_time
from
  aws_fms_app_list
where
  create_time >= datetime('now', '-30 day');
```

### Get application details of each app list
This query is used to gain insights into the applications within each list, including their names and network settings. This could be useful for understanding the structure and organization of your applications, particularly in terms of their network configurations.

```sql+postgres
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

```sql+sqlite
select
  list_name,
  list_id,
  json_extract(a.value, '$.AppName') as app_name,
  json_extract(a.value, '$.Port') as port,
  json_extract(a.value, '$.Protocol') as protocol
from
  aws_fms_app_list
  join json_each(aws_fms_app_list.apps_list, '$.AppsList') as a;
```