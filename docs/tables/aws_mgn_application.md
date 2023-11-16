---
title: "Table: aws_mgn_application - Query AWS Migration Service Applications using SQL"
description: "Allows users to query AWS Migration Service Applications to retrieve detailed information about each application."
---

# Table: aws_mgn_application - Query AWS Migration Service Applications using SQL

The `aws_mgn_application` table in Steampipe provides information about applications within AWS Migration Service. This table allows DevOps engineers to query application-specific details, such as application status, associated servers, and lifecycle details. Users can utilize this table to gather insights on applications, such as their current lifecycle state, last updated time, and associated server IDs. The schema outlines the various attributes of the Migration Service application, including the application ID, lifecycle, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_mgn_application` table, you can use the `.inspect aws_mgn_application` command in Steampipe.

**Key columns**:

- `application_id`: The unique identifier for the application. This can be used to join this table with others that also contain application_id.
- `lifecycle`: The lifecycle state of the application. This is useful for understanding the current status of the application.
- `last_updated_date_time`: The last updated date and time of the application. This is important for tracking changes and updates to the application.

## Examples

### Basic Info

```sql
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

```sql
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

### Get aggregated status details for an application

```sql
select
  name,
  application_id,
  application_aggregated_status ->> 'HealthStatus' as health_status,
  application_aggregated_status ->> 'ProgressStatus' as progress_status,
  application_aggregated_status ->> 'TotalSourceServers' as total_source_servers
from
  aws_mgn_application;
```

### List applications older than 30 days

```sql
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
