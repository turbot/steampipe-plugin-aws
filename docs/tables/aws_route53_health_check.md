---
title: "Steampipe Table: aws_route53_health_check - Query AWS Route 53 Health Check using SQL"
description: "Allows users to query AWS Route 53 Health Check data, providing information about health checks within AWS Route 53. This includes details such as health check configuration, health check status, and associated metadata."
folder: "Health"
---

# Table: aws_route53_health_check - Query AWS Route 53 Health Check using SQL

The AWS Route 53 Health Check is a feature of Amazon Route 53 service that helps to monitor the health and performance of your web applications, web servers, and other resources. It sends requests to your application, server, or other resource to verify that it's reachable, available, and functional. If the endpoint fails to respond, Route 53 can reroute traffic to healthy resources.

## Table Usage Guide

The `aws_route53_health_check` table in Steampipe provides you with information about health checks within AWS Route 53. This table allows you, as a DevOps engineer, to query health check-specific details, including health check configuration, health check status, and associated metadata. You can utilize this table to gather insights on health checks, such as health check configuration, health check status, and more. The schema outlines the various attributes of the health check for you, including the health check ID, health check version, type, and associated tags.

## Examples

### Basic Info
Explore which health checks are currently configured in your AWS Route53 setup. This can help you understand the health status of your resources and identify any potential issues or irregularities in your configuration.

```sql+postgres
select
  akas,
  id,
  health_check_version,
  health_check_config
from 
  aws_route53_health_check;
```

```sql+sqlite
select
  akas,
  id,
  health_check_version,
  health_check_config
from 
  aws_route53_health_check;
```

### List cloud watch configuration for health checks with monitoring enabled 
Explore the configurations of health checks that have monitoring enabled within your cloud environment. This can help you understand and manage your cloud health checks more effectively.

```sql+postgres
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

```sql+sqlite
select
  id,
  health_check_version,
  json_extract(cloud_watch_alarm_configuration, '$.ComparisonOperator') as cloud_watch_comparison_operator,
  json_extract(cloud_watch_alarm_configuration, '$.Dimensions') as cloud_watch_dimensions,
  json_extract(cloud_watch_alarm_configuration, '$.EvaluationPeriods') as cloud_watch_evaluation_periods,
  json_extract(cloud_watch_alarm_configuration, '$.MetricName') as cloud_watch_metric_name,
  json_extract(cloud_watch_alarm_configuration, '$.Period') as cloud_watch_period,
  json_extract(cloud_watch_alarm_configuration, '$.Statistic') as cloud_watch_statistic,
  json_extract(cloud_watch_alarm_configuration, '$.Threshold') as cloud_watch_threshold
from 
  aws_route53_health_check
where
  cloud_watch_alarm_configuration is not null;
```

### List health checks created by another service
Discover the segments that are monitored by health checks initiated by another service, enabling you to assess the elements within your system that are under external supervision. This can be particularly useful in identifying potential dependencies or points of failure in your infrastructure.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which health checks have been disabled within your AWS Route53 service. This can be useful for troubleshooting or auditing purposes, ensuring that all necessary health checks are active and functioning as expected.

```sql+postgres
select
  id,
  health_check_version,
  health_check_config ->> 'Disabled' as disabled
from 
  aws_route53_health_check 
where
  cast(health_check_config ->> 'Disabled' as boolean);
```

```sql+sqlite
select
  aws_route53_health_check.id,
  health_check_version,
  json_extract(health_check_config, '$.Disabled') as disabled
from 
  aws_route53_health_check 
where
  json_extract(health_check_config, '$.Disabled') = 'true';
```

### List health checks configuration details
Explore the configuration details of health checks to understand the status and performance of your domain and IP address. This can help in monitoring network connectivity and diagnosing any potential issues.

```sql+postgres
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

```sql+sqlite
select
  id,
  health_check_version,
  json_extract(health_check_config, '$.FullyQualifiedDomainName') as fully_qualified_domain_name,
  json_extract(health_check_config, '$.IPAddress') as ip_address,
  json_extract(health_check_config, '$.Port') as port,
  json_extract(health_check_config, '$.Type') as type,
  json_extract(health_check_config, '$.RequestInterval') as request_interval
from 
  aws_route53_health_check;
```

### List health checks where CloudWatch alarm is configured
Discover the segments of your network where health checks have been configured with a CloudWatch alarm. This can assist in identifying potential areas of concern and ensuring that your network is properly monitored.

```sql+postgres
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

```sql+sqlite
select
  id,
  health_check_version,
  json_extract(health_check_config, '$.FullyQualifiedDomainName') as fully_qualified_domain_name,
  json_extract(health_check_config, '$.IPAddress') as ip_address,
  json_extract(health_check_config, '$.Port') as port,
  json_extract(health_check_config, '$.Type') as type,
  json_extract(health_check_config, '$.RequestInterval') as request_interval,
  json_extract(health_check_config, '$.AlarmIdentifier') as alarm_identifier
from 
  aws_route53_health_check
where
  json_extract(health_check_config, '$.AlarmIdentifier') is not null;
```

### List details of failed health checks
Discover the segments that have failed health checks within your AWS Route53 service. This allows you to pinpoint specific areas of concern and take necessary corrective action.

```sql+postgres
select
  r.id,
  r.health_check_version,
  hc ->> 'IPAddress' as ip_address,
  hc ->> 'Region' as region,
  hc-> 'StatusReport' as status_report
from 
  aws_route53_health_check as r,
  jsonb_array_elements(health_check_status) hc 
where 
  hc-> 'StatusReport' ->> 'Status' not like '%Success%';
```

```sql+sqlite
select
  r.id,
  r.health_check_version,
  json_extract(hc.value, '$.IPAddress') as ip_address,
  json_extract(hc.value, '$.Region') as region,
  json_extract(hc.value, '$.StatusReport') as status_report
from 
  aws_route53_health_check as r,
  json_each(health_check_status) as hc
where 
  json_extract(json_extract(hc.value, '$.StatusReport'), '$.Status') not like '%Success%';
```