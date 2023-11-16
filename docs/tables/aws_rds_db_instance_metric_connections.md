---
title: "Table: aws_rds_db_instance_metric_connections - Query AWS RDS DBInstance Metrics using SQL"
description: "Allows users to query AWS RDS DBInstance Metrics for a comprehensive view of the number of database connections."
---

# Table: aws_rds_db_instance_metric_connections - Query AWS RDS DBInstance Metrics using SQL

The `aws_rds_db_instance_metric_connections` table in Steampipe provides information about the number of database connections to each Amazon RDS DB instance. This table allows DevOps engineers, database administrators, and other technical professionals to query connection-specific details, including the time at which the number of connections was recorded, the maximum number of connections during the period, and the number of data points used for the statistical calculation. Users can utilize this table to monitor and manage their database connections, analyze connection trends, and troubleshoot potential connection issues. The schema outlines the various attributes of the RDS DB instance connections, including the DB instance identifier, timestamp, sample count, average, minimum, and maximum number of connections.

The `aws_rds_db_instance_metric_connections` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_connections` table, you can use the `.inspect aws_rds_db_instance_metric_connections` command in Steampipe.

**Key columns**:

- `db_instance_identifier`: The identifier of the DB instance. This column is useful for joining with other tables that contain DB instance information.
- `timestamp`: The timestamp for the data point in UTC format. This column is important for tracking the historical trend of database connections.
- `maximum`: The maximum number of connections during the period. This column is useful for identifying potential spikes in database connections.

## Examples

### Basic info

```sql
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections
order by
  db_instance_identifier,
  timestamp;
```



### Intervals averaging over 100 connections

```sql
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```