---
title: "Steampipe Table: aws_servicequotas_auto_management_configuration - Query AWS Service Quotas Auto Management Configuration using SQL"
description: "Allows users to query AWS Service Quotas Automatic Management configuration, providing details on opt-in status, notification settings, and exclusion lists."
folder: "Service Quotas"
---

# Table: aws_servicequotas_auto_management_configuration - Query AWS Service Quotas Auto Management Configuration using SQL

AWS Service Quotas Automatic Management monitors your quota utilization and can automatically request increases when quotas approach their limits. This table returns the account-level configuration for Automatic Management per region.

## Table Usage Guide

The `aws_servicequotas_auto_management_configuration` table in Steampipe provides you with the Automatic Management configuration for each AWS region. This is a singleton table that returns one row per region, reflecting the current opt-in status, notification ARN, and any excluded services.

## Examples

### Check auto-management status across all regions
Determine the Automatic Management opt-in status for each region.

```sql+postgres
select
  region,
  opt_in_status,
  opt_in_type,
  opt_in_level
from
  aws_servicequotas_auto_management_configuration;
```

```sql+sqlite
select
  region,
  opt_in_status,
  opt_in_type,
  opt_in_level
from
  aws_servicequotas_auto_management_configuration;
```

### Find regions with auto-management enabled
Identify regions where Automatic Management is active.

```sql+postgres
select
  region,
  opt_in_type,
  notification_arn
from
  aws_servicequotas_auto_management_configuration
where
  opt_in_status = 'ENABLED';
```

```sql+sqlite
select
  region,
  opt_in_type,
  notification_arn
from
  aws_servicequotas_auto_management_configuration
where
  opt_in_status = 'ENABLED';
```

### View exclusion list details
Explore which services are excluded from Automatic Management.

```sql+postgres
select
  region,
  opt_in_status,
  exclusion_list
from
  aws_servicequotas_auto_management_configuration
where
  exclusion_list is not null;
```

```sql+sqlite
select
  region,
  opt_in_status,
  exclusion_list
from
  aws_servicequotas_auto_management_configuration
where
  exclusion_list is not null;
```
