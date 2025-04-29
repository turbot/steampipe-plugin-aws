---
title: "Steampipe Table: aws_lightsail_database - Query AWS Lightsail Databases using SQL"
description: "Allows users to query AWS Lightsail Databases for detailed information about database instances, including their configuration, status, and associated metadata."
folder: "Lightsail"
---

# Table: aws_lightsail_database - Query AWS Lightsail Databases using SQL

The AWS Lightsail Database is a managed database service that provides an easy way to set up, operate, and scale a relational database in the cloud. It offers a cost-effective, high-performance database solution with built-in security, backup, and maintenance features.

## Table Usage Guide

The `aws_lightsail_database` table in Steampipe provides you with information about AWS Lightsail database instances. This table allows you as a DevOps engineer to query database-specific details, including the database name, engine type, version, status, backup settings, and associated metadata. You can utilize this table to gather insights on databases, such as their current state, maintenance windows, and backup configurations. The schema outlines the various attributes of the AWS Lightsail database for you, including the database ARN, endpoint information, and associated tags.

## Examples

### Basic info
Explore the basic information about your Lightsail databases, including their names, engine types, and current status. This can help you understand the current state of your databases and identify any that might need attention.

```sql+postgres
select
  name,
  engine,
  engine_version,
  state,
  created_at
from
  aws_lightsail_database;
```

```sql+sqlite
select
  name,
  engine,
  engine_version,
  state,
  created_at
from
  aws_lightsail_database;
```

### List databases with backup settings
Identify databases and their backup configurations to ensure proper data protection.

```sql+postgres
select
  name,
  backup_retention_enabled,
  preferred_backup_window,
  preferred_maintenance_window
from
  aws_lightsail_database;
```

```sql+sqlite
select
  name,
  backup_retention_enabled,
  preferred_backup_window,
  preferred_maintenance_window
from
  aws_lightsail_database;
```

### List databases with endpoint information
View database endpoint details and access information to help with connection management.

```sql+postgres
select
  name,
  master_endpoint,
  master_username,
  publicly_accessible
from
  aws_lightsail_database;
```

```sql+sqlite
select
  name,
  master_endpoint,
  master_username,
  publicly_accessible
from
  aws_lightsail_database;
```

### List databases by state
Analyze the distribution of databases by their current state to understand your database landscape.

```sql+postgres
select
  state,
  count(*) as database_count
from
  aws_lightsail_database
group by
  state
order by
  database_count desc;
```

```sql+sqlite
select
  state,
  count(*) as database_count
from
  aws_lightsail_database
group by
  state
order by
  database_count desc;
```

### List databases with specific tags
Find databases that have specific tags associated with them to help organize and manage your databases based on custom criteria.

```sql+postgres
select
  name,
  engine,
  tags
from
  aws_lightsail_database
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  name,
  engine,
  tags
from
  aws_lightsail_database
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List databases by engine type
Analyze the distribution of databases by their engine type to understand your database technology stack.

```sql+postgres
select
  engine,
  engine_version,
  count(*) as database_count
from
  aws_lightsail_database
group by
  engine,
  engine_version
order by
  database_count desc;
```

```sql+sqlite
select
  engine,
  engine_version,
  count(*) as database_count
from
  aws_lightsail_database
group by
  engine,
  engine_version
order by
  database_count desc;
```

### List databases with high availability
Find databases that have high availability enabled by checking their secondary availability zone.

```sql+postgres
select
  name,
  engine,
  secondary_availability_zone
from
  aws_lightsail_database
where
  secondary_availability_zone is not null;
```

```sql+sqlite
select
  name,
  engine,
  secondary_availability_zone
from
  aws_lightsail_database
where
  secondary_availability_zone is not null;
```

### List databases with backup retention enabled
Find databases that have backup retention enabled to ensure proper data protection.

```sql+postgres
select
  name,
  backup_retention_enabled,
  preferred_backup_window
from
  aws_lightsail_database
where
  backup_retention_enabled = true;
```

```sql+sqlite
select
  name,
  backup_retention_enabled,
  preferred_backup_window
from
  aws_lightsail_database
where
  backup_retention_enabled = 1;
``` 