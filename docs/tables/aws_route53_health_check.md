---
title: "Table: aws_route53_health_check - Query AWS Route 53 Health Check using SQL"
description: "Allows users to query AWS Route 53 Health Check data, providing information about health checks within AWS Route 53. This includes details such as health check configuration, health check status, and associated metadata."
---

# Table: aws_route53_health_check - Query AWS Route 53 Health Check using SQL

The `aws_route53_health_check` table in Steampipe provides information about health checks within AWS Route 53. This table allows DevOps engineers to query health check-specific details, including health check configuration, health check status, and associated metadata. Users can utilize this table to gather insights on health checks, such as health check configuration, health check status, and more. The schema outlines the various attributes of the health check, including the health check ID, health check version, type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_health_check` table, you can use the `.inspect aws_route53_health_check` command in Steampipe.

**Key columns**:

- `id`: This is the unique identifier of the health check. It can be used to join this table with other tables.
- `health_check_version`: This is the version of the health check. It can be useful for version tracking and management.
- `type`: This column indicates the type of health check, which can be useful for filtering and categorization.

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
