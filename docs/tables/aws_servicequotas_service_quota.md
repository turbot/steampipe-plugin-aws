---
title: "Steampipe Table: aws_servicequotas_service_quota - Query AWS Service Quotas using SQL"
description: "Allows users to query AWS Service Quotas, providing detailed information about each quota's value, default value, and whether it's adjustable."
folder: "Service Quotas"
---

# Table: aws_servicequotas_service_quota - Query AWS Service Quotas using SQL

The AWS Service Quotas is a service that allows you to view and manage your quotas for AWS services from a central location. Quotas, also referred to as limits, are the maximum number of resources that you can create in your AWS account. With Service Quotas, you can easily see the quotas for the AWS services that are most relevant to you, request quota increases, and track the status of your requests.

## Table Usage Guide

The `aws_servicequotas_service_quota` table in Steampipe provides you with information about service quotas within AWS Service Quotas. This table allows you, as a DevOps engineer, to query quota-specific details, including quota value, default value, and whether it's adjustable. You can utilize this table to gather insights on service quotas, such as identifying quotas nearing their limit, understanding default and custom quota values, and determining which quotas can be adjusted. The schema outlines the various attributes of the service quota for you, including the quota ARN, quota code, service code, and associated metadata.

## Examples

### Basic info
Analyze the settings to understand the quotas set for various AWS services. This can help manage resources more effectively by identifying services that are nearing their quota limits, thus preventing service disruptions due to exceeded quotas.

```sql+postgres
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

```sql+sqlite
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
Explore which global service quotas are in place in your AWS environment. This can help manage service usage and avoid hitting service limit caps, thus ensuring smooth operation of your applications.

```sql+postgres
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

```sql+sqlite
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
  global_quota = 1;
```

### List service quotas for a specific service
Explore the set limits for a particular service to manage resource usage effectively. This can be particularly useful when planning resource allocation or troubleshooting resource limitations.

```sql+postgres
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

```sql+sqlite
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