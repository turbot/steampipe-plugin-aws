---
title: "Table: aws_servicequotas_service_quota_change_request - Query AWS Service Quotas using SQL"
description: "Allows users to query AWS Service Quotas change requests."
---

# Table: aws_servicequotas_service_quota_change_request - Query AWS Service Quotas using SQL

The `aws_servicequotas_service_quota_change_request` table in Steampipe provides information about service quota change requests within AWS Service Quotas. This table allows DevOps engineers to query request-specific details, including the status of the request, the requested value, and the associated metadata. Users can utilize this table to gather insights on service quota change requests, such as pending requests, details of requested increases in service quotas, and more. The schema outlines the various attributes of the service quota change request, including the service name, quota name, request id, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_servicequotas_service_quota_change_request` table, you can use the `.inspect aws_servicequotas_service_quota_change_request` command in Steampipe.

### Key columns:

- `service_code`: The service code corresponds to the service name for which the quota applies. This is a key column as it can be used to join this table with other tables that contain service-specific information.
- `quota_name`: This column represents the name of the specific quota within the service for which the change request has been made. It is useful in understanding the specific resources for which quota changes are requested.
- `request_id`: This is the unique identifier for the quota change request. It is important as it can be used to track the status of specific quota change requests.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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