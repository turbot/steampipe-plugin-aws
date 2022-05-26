# Table: aws_route53_health_check

Amazon Route 53 health checks monitor the health and performance of your web applications, web servers, and other resources.

## Examples

### Basic Info

```sql
select
  akas,
  id,
  health_check_version,
  health_check_config
from 
  aws_route53_health_check;
```

### List cloud watch configuration for health checks with monitoring enabled 

```sql
select
  id,
  health_check_version,
  cloud_watch_alarm_configuration ->> 'ComparisonOperator' as cloud_watch_comparison_operator,
  cloud_watch_alarm_configuration ->> 'Dimensions' as cloud_watch_dimensions,
  cloud_watch_alarm_configuration ->> 'EvaluationPeriods' as cloud_watch_evaluation_periods,
  cloud_watch_alarm_configuration ->> 'MetricName' as cloud_watch_metric_name,
  cloud_watch_alarm_configuration ->> 'Period' as cloud_watch_period,
  cloud_watch_alarm_configuration ->> 'Statistic' as cloud_watch_statistic,
  cloud_watch_alarm_configuration ->> 'Threshold' as cloud_watch_threshold
from 
  aws_route53_health_check
where
  cloud_watch_alarm_configuration is not null;
```

### List health checks created by another service

```sql
select
  id,
  health_check_version,
  linked_service_description,
  linked_service_principal
from 
  aws_route53_health_check
where
  linked_service_description is not null;
```

### List disabled health checks

```sql
select
  id,
  health_check_version,
  health_check_config ->> 'Disabled' as disabled
from 
  aws_route53_health_check 
where
  cast(health_check_config ->> 'Disabled' as boolean);
```

### List health checks configuration details

```sql
select
  id,
  health_check_version,
  health_check_config ->> 'FullyQualifiedDomainName' as fully_qualified_domain_name,
  health_check_config ->> 'IPAddress' as ip_address,
  health_check_config ->> 'Port' as port,
  health_check_config ->> 'Type' as type,
  health_check_config ->> 'RequestInterval' as request_interval
from 
  aws_route53_health_check;
```

### List health checks where CloudWatch alarm is configured

```sql
select
  id,
  health_check_version,
  health_check_config ->> 'FullyQualifiedDomainName' as fully_qualified_domain_name,
  health_check_config ->> 'IPAddress' as ip_address,
  health_check_config ->> 'Port' as port,
  health_check_config ->> 'Type' as type,
  health_check_config ->> 'RequestInterval' as request_interval,
  health_check_config ->> 'AlarmIdentifier' as alarm_identifier
from 
  aws_route53_health_check
where
  health_check_config ->> 'AlarmIdentifier' is not null;
```

### List details of failed health checks

```sql
select
  id,
  health_check_version,
  hc ->> 'IPAddress' as ip_address,
  hc ->> 'Region' as region,
  hc-> 'StatusReport' as status_report
from 
  aws_route53_health_check,
  jsonb_array_elements(health_check_status) hc 
where 
  hc-> 'StatusReport' ->> 'Status' not like '%Success%';
```
