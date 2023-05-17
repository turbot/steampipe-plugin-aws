# Table: aws_config_retention_configuration

AWS Config allows you to delete your data by specifying a retention period for your ConfigurationItems. When you specify a retention period, AWS Config retains your ConfigurationItems for that specified period.

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
