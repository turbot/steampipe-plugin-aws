---
title: "Steampipe Table: aws_mgn_application - Query AWS Migration Service Applications using SQL"
description: "Allows users to query AWS Migration Service Applications to retrieve detailed information about each application."
folder: "Application Migration Service"
---

# Table: aws_mgn_application - Query AWS Migration Service Applications using SQL

The AWS Migration Service (MGN) is designed to minimize downtime during migration. It allows you to quickly lift and shift applications to AWS with minimal changes, enabling rapid migration. The 'application' in AWS MGN represents a group of servers in a single application, which can be migrated together to AWS.

## Table Usage Guide

The `aws_mgn_application` table in Steampipe provides you with information about applications within AWS Migration Service. This table allows you, as a DevOps engineer, to query application-specific details, such as application status, associated servers, and lifecycle details. You can utilize this table to gather insights on applications, such as their current lifecycle state, last updated time, and associated server IDs. The schema outlines the various attributes of the Migration Service application for you, including the application ID, lifecycle, and associated tags.

## Examples

### Basic Info
Explore the basic information of your AWS Migration (MGN) applications to identify their archival status, creation dates, and associated tags. This can help in assessing the current state of your applications and in planning any requisite migration waves.

```sql+postgres
select
  name,
  arn,
  application_id,
  creation_date_time,
  is_archived,
  wave_id,
  tags
from
  aws_mgn_application;
```

```sql+sqlite
select
  name,
  arn,
  application_id,
  creation_date_time,
  is_archived,
  wave_id,
  tags
from
  aws_mgn_application;
```

### List archived applications
Identify instances where applications have been archived within the AWS Migration Service. This is useful for maintaining a clean and organized application environment by keeping track of archived applications.

```sql+postgres
select
  name,
  arn,
  application_id,
  creation_date_time,
  is_archived
from
  aws_mgn_application
where
  is_archived;
```

```sql+sqlite
select
  name,
  arn,
  application_id,
  creation_date_time,
  is_archived
from
  aws_mgn_application
where
  is_archived = 1;
```

### Get aggregated status details for an application
Explore the health and progress status of an application along with the total number of source servers. This can be useful in assessing the overall performance and progress of an application, and in identifying potential issues related to server resources.

```sql+postgres
select
  name,
  application_id,
  application_aggregated_status ->> 'HealthStatus' as health_status,
  application_aggregated_status ->> 'ProgressStatus' as progress_status,
  application_aggregated_status ->> 'TotalSourceServers' as total_source_servers
from
  aws_mgn_application;
```

```sql+sqlite
select
  name,
  application_id,
  json_extract(application_aggregated_status, '$.HealthStatus') as health_status,
  json_extract(application_aggregated_status, '$.ProgressStatus') as progress_status,
  json_extract(application_aggregated_status, '$.TotalSourceServers') as total_source_servers
from
  aws_mgn_application;
```

### List applications older than 30 days
Explore which applications have been in existence for more than 30 days. This can be useful in understanding the longevity of applications and their usage patterns over time.

```sql+postgres
select
  name,
  application_id,
  creation_date_time,
  is_archived,
  wave_id
from
  aws_mgn_application
where
  creation_date_time >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  application_id,
  creation_date_time,
  is_archived,
  wave_id
from
  aws_mgn_application
where
  creation_date_time >= datetime('now','-30 days');
```