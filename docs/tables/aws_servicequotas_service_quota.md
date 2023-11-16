---
title: "Table: aws_servicequotas_service_quota - Query AWS Service Quotas using SQL"
description: "Allows users to query AWS Service Quotas, providing detailed information about each quota's value, default value, and whether it's adjustable."
---

# Table: aws_servicequotas_service_quota - Query AWS Service Quotas using SQL

The `aws_servicequotas_service_quota` table in Steampipe provides information about service quotas within AWS Service Quotas. This table allows DevOps engineers to query quota-specific details, including quota value, default value, and whether it's adjustable. Users can utilize this table to gather insights on service quotas, such as identifying quotas nearing their limit, understanding default and custom quota values, and determining which quotas can be adjusted. The schema outlines the various attributes of the service quota, including the quota ARN, quota code, service code, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_servicequotas_service_quota` table, you can use the `.inspect aws_servicequotas_service_quota` command in Steampipe.

**Key columns**:

- `quota_arn`: The Amazon Resource Name (ARN) of the service quota. This can be used to join this table with other tables that include ARN information.
- `quota_code`: The code identifier for the service quota. This can be useful for joining with tables that use quota codes, or for filtering specific quotas.
- `service_code`: The code identifier for the service associated with the quota. This can be used to join this table with other AWS service tables or for filtering specific services.

## Examples

### Basic info

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_service_quota;
```

### List global service quotas

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_service_quota
where
  global_quota;
```

### List service quotas for a specific service

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_service_quota
where
  service_code = 'athena';
```