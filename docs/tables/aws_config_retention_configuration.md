---
title: "Steampipe Table: aws_config_retention_configuration - Query AWS Config Retention Configuration using SQL"
description: "Allows users to query AWS Config Retention Configuration for information about the retention period that AWS Config uses to retain your configuration items."
folder: "Config"
---

# Table: aws_config_retention_configuration - Query AWS Config Retention Configuration using SQL

The AWS Config Retention Configuration is a feature within the AWS Config service that allows you to specify the retention period (in days) for your configuration items. This helps in managing the volume of historical configuration items and reducing storage costs. AWS Config automatically deletes configuration items older than the specified retention period.

## Table Usage Guide

The `aws_config_retention_configuration` table in Steampipe provides you with information about the retention period that AWS Config uses to retain your configuration items. This table allows you, as a DevOps engineer, to query retention period details, including the number of days AWS Config retains the configuration items and whether the retention is permanent. You can utilize this table to gather insights on the retention configurations, such as the duration of retention and whether the retention is set to be permanent. The schema outlines the various attributes of the retention configuration for you, including the name of the retention period and the retention period in days.

## Examples

### Basic info
Explore which AWS Config retention configurations are active and determine the areas in which they are applied. This can help assess the elements within your AWS environment that have specific retention periods for configuration items, facilitating efficient resource management and compliance monitoring.

```sql+postgres
select
  name,
  retention_period_in_days,
  title,
  region
from
  aws_config_retention_configuration;
```

```sql+sqlite
select
  name,
  retention_period_in_days,
  title,
  region
from
  aws_config_retention_configuration;
```

### List retention configuration with the retention period less than 1 year
Discover the segments that have a retention period of less than a year in the AWS configuration. This can be useful to identify and review any potentially risky settings where data might not be retained long enough for compliance or auditing purposes.

```sql+postgres
select
  name,
  retention_period_in_days,
  title
from
  aws_config_retention_configuration
where
  retention_period_in_days < 356;
```

```sql+sqlite
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
Discover the segments that have specific retention configurations in a particular region. This can help in understanding how long configuration data is retained and can aid in better compliance management.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which retention settings of configuration recorders are applied, allowing you to understand how long your AWS Config data is retained in different regions.

```sql+postgres
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

```sql+sqlite
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