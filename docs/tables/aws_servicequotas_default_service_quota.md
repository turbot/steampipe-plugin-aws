---
title: "Table: aws_servicequotas_default_service_quota - Query AWS Service Quotas using SQL"
description: "Allows users to query AWS Service Quotas Default Service Quota to retrieve information about the default values of service quotas for AWS services."
---

# Table: aws_servicequotas_default_service_quota - Query AWS Service Quotas using SQL

The `aws_servicequotas_default_service_quota` table in Steampipe provides information about the default values of service quotas for AWS services. This table allows DevOps engineers to query quota-specific details, including quota names, quota codes, and the default values. Users can utilize this table to gather insights on service quotas, such as understanding the default limits for each AWS service, identifying services that might require a quota increase, and more. The schema outlines the various attributes of the service quota, including the quota ARN, service name, quota code, quota name, and the default value.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_servicequotas_default_service_quota` table, you can use the `.inspect aws_servicequotas_default_service_quota` command in Steampipe.

**Key columns**:

- `quota_arn`: The Amazon Resource Name (ARN) of the service quota. This can be used to join with other tables that use ARN as a unique identifier.
- `service_name`: The name of the AWS service. This is useful to join with other tables that provide service-specific information.
- `quota_code`: The unique code for the service quota. This can be used to join with other tables that provide quota-specific information.

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
  aws_servicequotas_default_service_quota;
```

### List global default service quotas

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_default_service_quota
where
  global_quota;
```

### List default service quotas for a specific service

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_default_service_quota
where
  service_code = 'athena';
```