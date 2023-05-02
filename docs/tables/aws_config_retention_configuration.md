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

