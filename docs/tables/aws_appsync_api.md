---
title: "Steampipe Table: aws_appsync_api - Query AWS AppSync APIs using SQL"
description: "Allows users to query AWS AppSync APIs (Event APIs) to retrieve detailed information about each API configuration."
folder: "AppSync"
---

# Table: aws_appsync_api - Query AWS AppSync APIs using SQL

AWS AppSync is a managed service that provides real-time data and offline capabilities for web and mobile applications. The `aws_appsync_api` table provides information about **AppSync Event APIs**, which are different from GraphQL APIs. Event APIs enable real-time messaging and event-driven architectures with WebSocket connections for publishing and subscribing to events.

## Table Usage Guide

The `aws_appsync_api` table in Steampipe provides you with information about **AppSync Event APIs** within AWS. This table allows you, as a DevOps engineer, to query Event API-specific details, including API ID, name, DNS endpoints, event configuration, and associated metadata. You can utilize this table to gather insights on Event APIs, such as their real-time messaging configuration, DNS endpoints, X-Ray tracing settings, and associated tags.


## Examples

### Basic info
Explore the features and settings of your AWS AppSync Event APIs to better understand their configuration, such as creation time, DNS endpoints, and regional distribution. This can help in assessing Event API performance and operational efficiency.

```sql+postgres
select
  name,
  api_id,
  arn,
  created,
  xray_enabled,
  region
from
  aws_appsync_api;
```

```sql+sqlite
select
  name,
  api_id,
  arn,
  created,
  xray_enabled,
  region
from
  aws_appsync_api;
```

### List Event APIs with X-Ray tracing enabled
Find AppSync Event APIs that have X-Ray tracing enabled to understand which Event APIs are being monitored for performance and debugging. This is particularly useful for real-time applications where performance monitoring is critical.

```sql+postgres
select
  name,
  api_id,
  created,
  xray_enabled
from
  aws_appsync_api
where
  xray_enabled = true;
```

```sql+sqlite
select
  name,
  api_id,
  created,
  xray_enabled
from
  aws_appsync_api
where
  xray_enabled = 1;
```

### Get Event API DNS endpoints
Retrieve DNS endpoints for Event APIs to understand the HTTP and WebSocket endpoints available for real-time communication. Event APIs provide both HTTP endpoints for REST operations and WebSocket endpoints for real-time messaging.

```sql+postgres
select
  name,
  api_id,
  dns
from
  aws_appsync_api
where
  dns is not null;
```

```sql+sqlite
select
  name,
  api_id,
  dns
from
  aws_appsync_api
where
  dns is not null;
```

### Get Event API details with tags
Retrieve detailed information about Event APIs including their associated tags for resource management and cost allocation.

```sql+postgres
select
  name,
  api_id,
  owner_contact,
  waf_web_acl_arn,
  tags
from
  aws_appsync_api
where
  tags is not null;
```

```sql+sqlite
select
  name,
  api_id,
  owner_contact,
  waf_web_acl_arn,
  tags
from
  aws_appsync_api
where
  tags is not null;
```

### List recently created Event APIs
Find Event APIs that were created recently to track new deployments and changes in your AppSync Event API environment.

```sql+postgres
select
  name,
  api_id,
  created,
  xray_enabled
from
  aws_appsync_api
where
  created >= now() - interval '30 days'
order by
  created desc;
```

### Get Event API configuration details
Retrieve detailed Event API configuration including event settings for real-time messaging capabilities.

```sql+postgres
select
  name,
  api_id,
  event_config,
  dns
from
  aws_appsync_api
where
  event_config is not null;
```

```sql+sqlite
select
  name,
  api_id,
  event_config,
  dns
from
  aws_appsync_api
where
  event_config is not null;
```

### Get Event API cache configuration
Retrieve cache configuration for Event APIs to understand caching settings and performance optimization.

```sql+postgres
select
  name,
  api_id,
  api_cache
from
  aws_appsync_api
where
  api_cache is not null;
```

```sql+sqlite
select
  name,
  api_id,
  api_cache
from
  aws_appsync_api
where
  api_cache is not null;
```
