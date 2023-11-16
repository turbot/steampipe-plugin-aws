---
title: "Table: aws_config_retention_configuration - Query AWS Config Retention Configuration using SQL"
description: "Allows users to query AWS Config Retention Configuration for information about the retention period that AWS Config uses to retain your configuration items."
---

# Table: aws_config_retention_configuration - Query AWS Config Retention Configuration using SQL

The `aws_config_retention_configuration` table in Steampipe provides information about the retention period that AWS Config uses to retain your configuration items. This table allows DevOps engineers to query retention period details, including the number of days AWS Config retains the configuration items and whether the retention is permanent. Users can utilize this table to gather insights on the retention configurations, such as the duration of retention and whether the retention is set to be permanent. The schema outlines the various attributes of the retention configuration, including the name of the retention period and the retention period in days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_config_retention_configuration` table, you can use the `.inspect aws_config_retention_configuration` command in Steampipe.

**Key columns**:

- `name`: The name of the retention configuration. This is useful for identifying the specific configuration.
- `retention_period_in_days`: The number of days AWS Config retains the configuration items. This is important for understanding how long the configuration items are retained.
- `retention_period_in_days`: Indicates whether the retention is permanent. This is crucial for determining whether the retention is set to be permanent or not.

## Examples

### Basic info

```sql
select
  name,
  retention_period_in_days,
  title,
  region
from
  aws_config_retention_configuration;
```

### List retention configuration with the retention period less than 1 year

```sql
select
  name,
  retention_period_in_days,
  title
from
  aws_config_retention_configuration
where
  retention_period_in_days < 356;
```

### List retention configuration by region

```sql
select
  name,
  retention_period_in_days,
  title,
  region
from
  aws_config_retention_configuration
where
  region = 'us-east-1';
```

### List retention configuration settings of config recorders

```sql
select
  c.title as configuration_recorder,
  r.name as retention_configuration_name,
  r.retention_period_in_days,
  r.region
from
  aws_config_retention_configuration as r
  left join aws_config_configuration_recorder as c
on
  r.region = c.region;
```
