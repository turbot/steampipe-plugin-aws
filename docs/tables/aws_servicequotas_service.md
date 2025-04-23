---
title: "Steampipe Table: aws_servicequotas_service - Query AWS Service Quota Services using SQL"
description: "Allows users to query AWS Service Quotas services, providing detailed information about each service's code and name."
folder: "Service Quotas"
---

# Table: aws_servicequotas_service - Query AWS Service Quota Services using SQL

The AWS Service Quotas is a service that allows you to view and manage your quotas for AWS services from a central location. AWS services that integrate with the Service Quotas service each have a code and name.

## Table Usage Guide

The `aws_servicequotas_service` table in Steampipe provides you with information about services within AWS Service Quotas. With the service's service code, you can then query service quotas for that specific service, including default values, 

## Examples

### Basic info
Explore the services supported by the Service Quotas service.

```sql+postgres
select distinct
  service_code,
  service_name
from
  aws_servicequotas_service;
```

```sql+sqlite
select distinct
  service_code,
  service_name
from
  aws_servicequotas_service;
```

### Get a service code for a a specific service
Get the service code for a specific service.

```sql+postgres
select distinct
  service_code
from
  aws_servicequotas_service
where
  service_name = 'AWS CloudTrail';
```

```sql+sqlite
select distinct
  service_code
from
  aws_servicequotas_service
where
  service_name = 'AWS CloudTrail';
```
