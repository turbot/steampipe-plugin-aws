---
title: "Steampipe Table: aws_application_signals_service_level_objective - Query AWS Application Signals Service Level Objectives using SQL"
description: "Allows users to query AWS Application Signals Service Level Objectives, providing detailed information on SLO configurations, goals, and performance metrics."
---

# Table: aws_application_signals_service_level_objective - Query AWS Application Signals Service Level Objectives using SQL

AWS Application Signals Service Level Objectives (SLOs) help you monitor and manage the reliability and performance of your applications. Each SLO defines performance goals and evaluates the service’s ability to meet those goals over a given period. The `aws_application_signals_service_level_objective` table in Steampipe allows you to query detailed information about SLOs in your AWS environment, including their goals, attainment levels, evaluation types, and performance metrics.

## Table Usage Guide

The `aws_application_signals_service_level_objective` table enables cloud administrators, DevOps engineers, and reliability engineers to monitor the performance of their applications by querying SLOs. You can gather insights into various aspects of the SLO, such as its attainment goal, operation name, performance metrics, and whether it’s request-based or period-based. This table is useful for tracking application reliability, analyzing SLO compliance, and ensuring your application meets the defined performance standards.

## Examples

### Basic SLO information
Retrieve basic information about SLOs, including their name, ARN, and region.

```sql+postgres
select
  name,
  arn,
  region,
  created_time
from
  aws_application_signals_service_level_objective;
```

```sql+sqlite
select
  name,
  arn,
  region,
  created_time
from
  aws_application_signals_service_level_objective;
```

### List SLOs by attainment goal
Identify service level objectives based on their attainment goals.

```sql+postgres
select
  name,
  attainment_goal,
  region,
  account_id
from
  aws_application_signals_service_level_objective
where
  attainment_goal >= 99.0;
```

```sql+sqlite
select
  name,
  attainment_goal,
  region,
  account_id
from
  aws_application_signals_service_level_objective
where
  attainment_goal >= 99.0;
```

### List request-based SLOs
Retrieve all service level objectives that are request-based, which monitor performance based on requests.

```sql+postgres
select
  name,
  arn,
  region,
  evaluation_type,
  request_based_sli
from
  aws_application_signals_service_level_objective
where
  evaluation_type = 'RequestBased';
```

```sql+sqlite
select
  name,
  arn,
  region,
  evaluation_type,
  request_based_sli
from
  aws_application_signals_service_level_objective
where
  evaluation_type = 'RequestBased';
```

### List SLOs by operation name
Identify service level objectives specific to certain operations.

```sql+postgres
select
  name,
  operation_name,
  arn,
  region
from
  aws_application_signals_service_level_objective
where
  operation_name is not null;
```

```sql+sqlite
select
  name,
  operation_name,
  arn,
  region
from
  aws_application_signals_service_level_objective
where
  operation_name is not null;
```

### List service level objectives with select service level indicator details
Retrieve a list of service level objectives (SLOs) along with specific service level indicator (SLI) details, such as comparison operators and metric thresholds.

```sql+postgres
select
  arn,
  name,
  sli::json -> 'ComparisonOperator' as "Must Be",
  sli::json -> 'MetricThreshold' as "Threshold"
  region
from
  aws_application_signals_service_level_objective;
```

```sql+sqlite
select
  arn,
  name,
  json_extract(sli, '$.ComparisonOperator') as "Must Be",
  json_extract(sli, '$.MetricThreshold') as "Threshold"
  region
from
  aws_application_signals_service_level_objective
```
