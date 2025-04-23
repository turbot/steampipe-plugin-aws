---
title: "Steampipe Table: aws_servicequotas_service_quota_change_request - Query AWS Service Quotas using SQL"
description: "Allows users to query AWS Service Quotas change requests."
folder: "Service Quotas"
---

# Table: aws_servicequotas_service_quota_change_request - Query AWS Service Quotas using SQL

The AWS Service Quotas allows you to view and manage your quotas for AWS services from a central location. Service Quotas provides a unified view of the quotas for AWS services, and it allows you to request quota increases for the services that support it. It offers the ability to manage quotas for the number of resources that you can create in an AWS account.

## Table Usage Guide

The `aws_servicequotas_service_quota_change_request` table in Steampipe provides you with information about service quota change requests within AWS Service Quotas. This table enables you, as a DevOps engineer, to query request-specific details, including the status of the request, the requested value, and the associated metadata. You can utilize this table to gather insights on service quota change requests, such as pending requests, details of requested increases in service quotas, and more. The schema outlines the various attributes of the service quota change request for you, including the service name, quota name, request id, and associated tags.

## Examples

### Basic info
Identify instances where changes have been requested for AWS service quotas. This is particularly useful for understanding and managing resource allocation and usage limits within your AWS environment.

```sql+postgres
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request;
```

```sql+sqlite
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request;
```

### List denied service quota change requests
Explore which service quota change requests have been denied in your AWS environment. This could be used to identify potential bottlenecks or limitations in your resource allocation.

```sql+postgres
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request
where
  status = 'DENIED';
```

```sql+sqlite
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request
where
  status = 'DENIED';
```

### List service quota change requests for a specific service
Determine the status of requests to change service quotas for a specific service. This is useful for tracking any modifications to service limits, particularly for services that are heavily utilized or critical to operations.

```sql+postgres
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request
where
  service_code = 'athena';
```

```sql+sqlite
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request
where
  service_code = 'athena';
```