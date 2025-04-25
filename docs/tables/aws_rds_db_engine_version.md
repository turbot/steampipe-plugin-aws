---
title: "Steampipe Table: aws_rds_db_engine_version - Query AWS RDS DB Engine Versions using SQL"
description: "Enables users to query AWS RDS DB Engine Versions to retrieve detailed information on various database engine versions supported by Amazon RDS."
folder: "RDS"
---

# Table: aws_rds_db_engine_version - Query AWS RDS DB Engine Versions using SQL

Amazon Relational Database Service (RDS) supports various database engines, allowing you to run different versions of these engines. Understanding the specific features, capabilities, and limitations of each engine version is crucial for database administration and optimization.

The `aws_rds_db_engine_version` table in Steampipe provides comprehensive information on the different database engine versions available in Amazon RDS. This includes details such as the engine type, version number, status (e.g., whether the version is still supported), and specific attributes like supported feature names, character sets, and whether the engine version supports read replicas or global databases.

Utilizing this table, database administrators and DevOps engineers can make informed decisions regarding database engine upgrades, compatibility checks, and feature utilization. The schema includes attributes critical for understanding each engine version's capabilities and limitations.

**Important notes:**
- For improved performance, this table supports optional quals. Queries with optional quals are optimised to use RDS Engine Version filters. Optional quals are supported for the following columns:
  - `engine`
  - `engine_version`
  - `db_parameter_group_family`
  - `list_supported_character_sets`
  - `list_supported_timezones`
  - `engine_mode`
  - `default_only`
  - `status`

## Examples

### List all available engine versions
This query provides a list of all database engine versions available in your AWS environment, including their status and major version number. This is useful for identifying potential upgrades or assessing version availability.

```sql+postgres
select
  engine,
  engine_version,
  db_engine_version_description,
  status,
  major_engine_version
from
  aws_rds_db_engine_version;
```

```sql+sqlite
select
  engine,
  engine_version,
  db_engine_version_description,
  status,
  major_engine_version
from
  aws_rds_db_engine_version;
```

### Identify engines supporting read replicas
Discover which database engine versions support read replicas. This information is crucial for planning high availability and read scalability.

```sql+postgres
select
  engine,
  engine_version,
  supports_read_replica
from
  aws_rds_db_engine_version
where
  supports_read_replica;
```

```sql+sqlite
select
  engine,
  engine_version,
  supports_read_replica
from
  aws_rds_db_engine_version
where
  supports_read_replica;
```

### Engines with deprecation status
Find out which engine versions are deprecated. This can help in planning for necessary upgrades to maintain support and compatibility.

```sql+postgres
select
  engine,
  engine_version,
  status
from
  aws_rds_db_engine_version
where
  status = 'deprecated';
```

```sql+sqlite
select
  engine,
  engine_version,
  status
from
  aws_rds_db_engine_version
where
  status = 'deprecated';
```

### Supported features by engine version
List the features supported by a specific engine version. Adjust `'specific_engine_version'` to your engine version of interest. This query aids in understanding the capabilities of a given engine version.

```sql+postgres
select
  engine,
  engine_version,
  supported_feature_names
from
  aws_rds_db_engine_version
where
  engine_version = 'specific_engine_version';
```

```sql+sqlite
select
  engine,
  engine_version,
  supported_feature_names
from
  aws_rds_db_engine_version
where
  engine_version = 'specific_engine_version';
```

### List default engine version of engines
List only default engine version of engines.

```sql+postgres
select
  engine,
  engine_version,
  create_time,
  status,
  db_engine_media_type,
  default_only
from
  aws_rds_db_engine_version
where
  default_only;
```

```sql+sqlite
select
  engine,
  engine_version,
  create_time,
  status,
  db_engine_media_type,
  default_only
from
  aws_rds_db_engine_version
where
  default_only;
```
