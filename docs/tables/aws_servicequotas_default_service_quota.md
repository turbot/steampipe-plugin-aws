---
title: "Steampipe Table: aws_servicequotas_default_service_quota - Query AWS Service Quotas using SQL"
description: "Allows users to query AWS Service Quotas Default Service Quota to retrieve information about the default values of service quotas for AWS services."
folder: "Service Quotas"
---

# Table: aws_servicequotas_default_service_quota - Query AWS Service Quotas using SQL

The AWS Service Quotas allows you to view and manage your quotas for AWS services from a central location. Quotas, also referred to as limits, are the maximum number of specific resources that you can create in your AWS account. Default Service Quota is the default limit set by AWS for a resource within a specific region.

## Table Usage Guide

The `aws_servicequotas_default_service_quota` table in Steampipe provides you with information about the default values of service quotas for AWS services. This table allows you, as a DevOps engineer, to query quota-specific details, including quota names, quota codes, and the default values. You can utilize this table to gather insights on service quotas, such as understanding the default limits for each AWS service, identifying services that might require a quota increase, and more. The schema outlines the various attributes of the service quota for you, including the quota ARN, service name, quota code, quota name, and the default value.

## Examples

### Basic info
Identify the default service quotas in your AWS environment to understand your current usage and potential scalability. This can help in planning resource utilization and in avoiding service disruptions due to quota limitations.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are set as global default service quotas in your AWS services. This can help you gain insights into your resource usage and manage your AWS services more effectively.

```sql+postgres
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

```sql+sqlite
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
  global_quota = 1;
```

### List default service quotas for a specific service
Explore the default service quotas for a specific service to better manage resource usage and avoid hitting limits unexpectedly. This is useful for keeping track of your resource allocation and ensuring your operations run smoothly.

```sql+postgres
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

```sql+sqlite
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