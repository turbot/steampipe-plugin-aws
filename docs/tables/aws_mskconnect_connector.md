---
title: "Steampipe Table: aws_mskconnect_connector - Query AWS MSK Connect Connectors using SQL"
description: "Allows users to query AWS MSK Connect connectors to retrieve detailed information about each connector configuration, state, and associated metadata."
folder: "MSK"
---

# Table: aws_mskconnect_connector - Query AWS MSK Connect Connectors using SQL

AWS MSK Connect is a fully managed service that makes it easy to deploy and run Apache Kafka Connect applications on AWS. It allows you to easily move data between Apache Kafka clusters and other data stores using pre-built connectors. The `aws_mskconnect_connector` table provides information about MSK Connect connectors, including their configuration, state, and metadata.

## Table Usage Guide

The `aws_mskconnect_connector` table in Steampipe provides you with information about MSK Connect connectors within AWS. This table allows you, as a DevOps engineer, to query connector-specific details, including connector ARN, name, state, configuration, and associated metadata. You can utilize this table to gather insights on connectors, such as their current state, Kafka cluster connections, plugin configurations, and more. The schema outlines the various attributes of the MSK Connect connector for you, including the connector name, ARN, state, capacity settings, and associated tags.

## Examples

### Basic info
Explore the status and details of your AWS MSK Connect connectors to understand their configuration and operational state. This can be useful for auditing purposes or to identify potential issues with connector setup.

```sql+postgres
select
  connector_name,
  arn,
  connector_state,
  creation_time,
  current_version,
  kafka_connect_version
from
  aws_mskconnect_connector;
```

```sql+sqlite
select
  connector_name,
  arn,
  connector_state,
  creation_time,
  current_version,
  kafka_connect_version
from
  aws_mskconnect_connector;
```

### List active connectors
Identify MSK Connect connectors that are currently active and operational.

```sql+postgres
select
  connector_name,
  arn,
  connector_state,
  creation_time
from
  aws_mskconnect_connector
where
  connector_state = 'RUNNING';
```

```sql+sqlite
select
  connector_name,
  arn,
  connector_state,
  creation_time
from
  aws_mskconnect_connector
where
  connector_state = 'RUNNING';
```

### Find connectors with specific Kafka Connect versions
Identify connectors using specific versions of Kafka Connect for version management and compatibility assessment.

```sql+postgres
select
  connector_name,
  arn,
  kafka_connect_version,
  connector_state
from
  aws_mskconnect_connector
where
  kafka_connect_version = '2.7.1';
```

```sql+sqlite
select
  connector_name,
  arn,
  kafka_connect_version,
  connector_state
from
  aws_mskconnect_connector
where
  kafka_connect_version = '2.7.1';
```

### Get connector configuration details
Retrieve detailed configuration information for connectors including their settings and parameters.

```sql+postgres
select
  connector_name,
  arn,
  connector_configuration,
  capacity
from
  aws_mskconnect_connector
where
  connector_configuration is not null;
```

```sql+sqlite
select
  connector_name,
  arn,
  connector_configuration,
  capacity
from
  aws_mskconnect_connector
where
  connector_configuration is not null;
```

### List connectors with CloudWatch log delivery enabled
Identify connectors that have log delivery to CloudWatch Logs configured for monitoring and debugging.

```sql+postgres
select
  connector_name,
  arn,
  log_delivery,
  connector_state
from
  aws_mskconnect_connector
where
  log_delivery is not null;
```

```sql+sqlite
select
  connector_name,
  arn,
  log_delivery,
  connector_state
from
  aws_mskconnect_connector
where
  log_delivery is not null;
```

### Find connectors with specific authentication types
Identify connectors using specific authentication methods for Kafka cluster connections.

```sql+postgres
select
  connector_name,
  arn,
  kafka_cluster_client_authentication,
  connector_state
from
  aws_mskconnect_connector
where
  kafka_cluster_client_authentication is not null;
```

```sql+sqlite
select
  connector_name,
  arn,
  kafka_cluster_client_authentication,
  connector_state
from
  aws_mskconnect_connector
where
  kafka_cluster_client_authentication is not null;
```

### List recently created connectors
Find connectors that were created recently to track new deployments and changes in your data pipeline environment.

```sql+postgres
select
  connector_name,
  arn,
  creation_time,
  connector_state
from
  aws_mskconnect_connector
where
  creation_time >= now() - interval '30 days'
order by
  creation_time desc;
```

```sql+sqlite
select
  connector_name,
  arn,
  creation_time,
  connector_state
from
  aws_mskconnect_connector
where
  creation_time >= datetime('now', '-30 days')
order by
  creation_time desc;
```

### Find connectors with state issues
Identify connectors that may have state issues or failed operations for troubleshooting.

```sql+postgres
select
  connector_name,
  arn,
  connector_state,
  state_description
from
  aws_mskconnect_connector
where
  connector_state not in ('RUNNING', 'CREATING')
  and state_description is not null;
```

```sql+sqlite
select
  connector_name,
  arn,
  connector_state,
  state_description
from
  aws_mskconnect_connector
where
  connector_state not in ('RUNNING', 'CREATING')
  and state_description is not null;
```

### Check connector capacity settings
Analyze the capacity settings of connectors to understand their resource allocation and scaling configuration.

```sql+postgres
select
  connector_name,
  arn,
  capacity,
  connector_state
from
  aws_mskconnect_connector
where
  capacity is not null;
```

```sql+sqlite
select
  connector_name,
  arn,
  capacity,
  connector_state
from
  aws_mskconnect_connector
where
  capacity is not null;
```

### List connectors by service execution role
Identify connectors by their IAM service execution roles for security and access control analysis.

```sql+postgres
select
  connector_name,
  arn,
  service_execution_role_arn,
  connector_state
from
  aws_mskconnect_connector
where
  service_execution_role_arn is not null;
```

```sql+sqlite
select
  connector_name,
  arn,
  service_execution_role_arn,
  connector_state
from
  aws_mskconnect_connector
where
  service_execution_role_arn is not null;
```

### Monitor connector state distribution
Analyze the distribution of connector states across your MSK Connect environment for operational insights.

```sql+postgres
select
  connector_state,
  count(*) as connector_count
from
  aws_mskconnect_connector
group by
  connector_state
order by
  connector_count desc;
```

```sql+sqlite
select
  connector_state,
  count(*) as connector_count
from
  aws_mskconnect_connector
group by
  connector_state
order by
  connector_count desc;
```
