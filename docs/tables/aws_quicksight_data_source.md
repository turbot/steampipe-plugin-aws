---
title: "Steampipe Table: aws_quicksight_data_source - Query AWS QuickSight Data Sources using SQL"
description: "Allows users to query AWS QuickSight Data Sources, providing details about data source configurations, connections, and status information."
folder: "QuickSight"
---

# Table: aws_quicksight_data_source - Query AWS QuickSight Data Sources using SQL

AWS QuickSight Data Source is a connection to your data that QuickSight uses to create datasets. It contains the connection information, credentials, and other parameters needed to access your data from various sources like Amazon S3, Amazon RDS, or other supported data sources.

## Table Usage Guide

The `aws_quicksight_data_source` table in Steampipe provides you with information about data sources within AWS QuickSight. This table allows you, as a data analyst or administrator, to query data source-specific details, including connection properties, status, and configuration parameters. You can utilize this table to gather insights on data sources, such as their creation time, last update time, connection status, and type.

## Examples

### Basic info
Explore the fundamental details of your QuickSight data sources to understand their configuration and connection types.

```sql+postgres
select
  name,
  data_source_id,
  arn,
  created_time,
  last_updated_time,
  status,
  type
from
  aws_quicksight_data_source;
```

```sql+sqlite
select
  name,
  data_source_id,
  arn,
  created_time,
  last_updated_time,
  status,
  type
from
  aws_quicksight_data_source;
```

### List failed data sources
Identify data sources that have failed to help troubleshoot connection issues.

```sql+postgres
select
  name,
  data_source_id,
  status,
  error_info
from
  aws_quicksight_data_source
where
  status like '%FAILED%';
```

```sql+sqlite
select
  name,
  data_source_id,
  status,
  error_info
from
  aws_quicksight_data_source
where
  status like '%FAILED%';
```

### Get S3 data sources
Analyze the configuration of S3-based data sources to understand your data connections.

```sql+postgres
select
  name,
  data_source_id,
  status,
  data_source_parameters
from
  aws_quicksight_data_source
where
  type = 'S3';
```

```sql+sqlite
select
  name,
  data_source_id,
  status,
  data_source_parameters
from
  aws_quicksight_data_source
where
  type = 'S3';
```

### List data sources with VPC connections
Identify data sources that are configured with VPC connections for network security assessment.

```sql+postgres
select
  name,
  data_source_id,
  vpc_connection_properties
from
  aws_quicksight_data_source
where
  vpc_connection_properties is not null;
```

```sql+sqlite
select
  name,
  data_source_id,
  vpc_connection_properties
from
  aws_quicksight_data_source
where
  vpc_connection_properties is not null;
```

### Get data sources with SSL configuration
Review data sources that have SSL properties configured for secure connections.

```sql+postgres
select
  name,
  data_source_id,
  ssl_properties
from
  aws_quicksight_data_source
where
  ssl_properties is not null;
```

```sql+sqlite
select
  name,
  data_source_id,
  ssl_properties
from
  aws_quicksight_data_source
where
  ssl_properties is not null;
```

### List data sources with alternate parameters
Find data sources that have alternate connection parameters configured.

```sql+postgres
select
  name,
  data_source_id,
  type,
  alternate_data_source_parameters
from
  aws_quicksight_data_source
where
  alternate_data_source_parameters is not null;
```

```sql+sqlite
select
  name,
  data_source_id,
  type,
  alternate_data_source_parameters
from
  aws_quicksight_data_source
where
  alternate_data_source_parameters is not null;
```
