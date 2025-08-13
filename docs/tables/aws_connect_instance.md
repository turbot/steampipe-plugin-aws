---
title: "Steampipe Table: aws_connect_instance - Query AWS Connect Instances using SQL"
description: "Allows users to query AWS Connect instances to retrieve detailed information about each contact center instance configuration."
folder: "Connect"
---

# Table: aws_connect_instance - Query AWS Connect Instances using SQL

AWS Connect is a cloud contact center service that makes it easy for any business to deliver better customer service at lower cost. The `aws_connect_instance` table provides information about Connect instances, including their configuration, status, and metadata.

## Table Usage Guide

The `aws_connect_instance` table in Steampipe provides you with information about Connect instances within AWS. This table allows you, as a DevOps engineer, to query instance-specific details, including instance ID, alias, status, identity management type, and associated metadata. You can utilize this table to gather insights on Connect instances, such as their current status, call configuration, access URLs, and associated tags.

## Examples

### Basic info
Explore the features and settings of your AWS Connect instances to better understand their configuration, such as status, identity management type, and regional distribution. This can help in assessing instance performance and operational efficiency.

```sql+postgres
select
  instance_alias,
  id,
  arn,
  instance_status,
  identity_management_type,
  region
from
  aws_connect_instance;
```

```sql+sqlite
select
  instance_alias,
  id,
  arn,
  instance_status,
  identity_management_type,
  region
from
  aws_connect_instance;
```

### List active Connect instances
Find Connect instances that are currently active to understand which contact centers are operational.

```sql+postgres
select
  instance_alias,
  id,
  instance_status,
  created_time
from
  aws_connect_instance
where
  instance_status = 'ACTIVE';
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_status,
  created_time
from
  aws_connect_instance
where
  instance_status = 'ACTIVE';
```

### Find instances with both inbound and outbound calls enabled
Identify Connect instances that have both inbound and outbound calling capabilities enabled for comprehensive contact center functionality.

```sql+postgres
select
  instance_alias,
  id,
  inbound_calls_enabled,
  outbound_calls_enabled,
  instance_access_url
from
  aws_connect_instance
where
  inbound_calls_enabled = true
  and outbound_calls_enabled = true;
```

```sql+sqlite
select
  instance_alias,
  id,
  inbound_calls_enabled,
  outbound_calls_enabled,
  instance_access_url
from
  aws_connect_instance
where
  inbound_calls_enabled = 1
  and outbound_calls_enabled = 1;
```

### Get instance details with tags
Retrieve detailed information about Connect instances including their associated tags for resource management and cost allocation.

```sql+postgres
select
  instance_alias,
  id,
  service_role,
  tags
from
  aws_connect_instance
where
  tags is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  service_role,
  tags
from
  aws_connect_instance
where
  tags is not null;
```

### List recently created instances
Find Connect instances that were created recently to track new deployments and changes in your contact center environment.

```sql+postgres
select
  instance_alias,
  id,
  created_time,
  instance_status
from
  aws_connect_instance
where
  created_time >= now() - interval '30 days'
order by
  created_time desc;
```

```sql+sqlite
select
  instance_alias,
  id,
  created_time,
  instance_status
from
  aws_connect_instance
where
  created_time >= datetime('now', '-30 days')
order by
  created_time desc;
```

### Find instances with status issues
Identify Connect instances that may have status issues or failed creation attempts for troubleshooting.

```sql+postgres
select
  instance_alias,
  id,
  instance_status,
  status_reason
from
  aws_connect_instance
where
  status_reason is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_status,
  status_reason
from
  aws_connect_instance
where
  status_reason is not null;
```

### Get instance access URLs
Retrieve access URLs for Connect instances to understand how contact center users can access the admin website.

```sql+postgres
select
  instance_alias,
  id,
  instance_access_url,
  instance_status
from
  aws_connect_instance
where
  instance_access_url is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_access_url,
  instance_status
from
  aws_connect_instance
where
  instance_access_url is not null;
```

### Find instances with specific attributes enabled
Identify Connect instances with specific features enabled, such as Contact Lens, early media, or multi-party conference capabilities.

```sql+postgres
select
  instance_alias,
  id,
  contact_lens,
  early_media,
  multi_party_conference,
  enhanced_contact_monitoring
from
  aws_connect_instance
where
  contact_lens = 'true'
  or early_media = 'true'
  or multi_party_conference = 'true';
```

```sql+sqlite
select
  instance_alias,
  id,
  contact_lens,
  early_media,
  multi_party_conference,
  enhanced_contact_monitoring
from
  aws_connect_instance
where
  contact_lens = 'true'
  or early_media = 'true'
  or multi_party_conference = 'true';
```

### Get all instance attributes
Retrieve all instance attributes as a JSON object to understand the complete configuration of Connect instances.

```sql+postgres
select
  instance_alias,
  id,
  instance_attributes
from
  aws_connect_instance
where
  instance_attributes is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_attributes
from
  aws_connect_instance
where
  instance_attributes is not null;
```

### Find instances with contact lens enabled
Identify Connect instances that have Contact Lens enabled for advanced analytics and insights.

```sql+postgres
select
  instance_alias,
  id,
  contact_lens,
  enhanced_contact_monitoring,
  enhanced_chat_monitoring
from
  aws_connect_instance
where
  contact_lens = 'true';
```

```sql+sqlite
select
  instance_alias,
  id,
  contact_lens,
  enhanced_contact_monitoring,
  enhanced_chat_monitoring
from
  aws_connect_instance
where
  contact_lens = 'true';
```

### Find instances with bot management capabilities
Identify Connect instances that have bot management and analytics features enabled.

```sql+postgres
select
  instance_alias,
  id,
  bot_management,
  enable_bot_analytics_and_transcripts,
  automated_interaction_log
from
  aws_connect_instance
where
  bot_management = 'true'
  or enable_bot_analytics_and_transcripts = 'true';
```

```sql+sqlite
select
  instance_alias,
  id,
  bot_management,
  enable_bot_analytics_and_transcripts,
  automated_interaction_log
from
  aws_connect_instance
where
  bot_management = 'true'
  or enable_bot_analytics_and_transcripts = 'true';
```

### Find instances with advanced features
Identify Connect instances with advanced features like multi-party conference, forecasting, and high volume outbound.

```sql+postgres
select
  instance_alias,
  id,
  multi_party_conference,
  multi_party_chat_conference,
  forecasting_planning_scheduling,
  high_volume_outbound
from
  aws_connect_instance
where
  multi_party_conference = 'true'
  or forecasting_planning_scheduling = 'true'
  or high_volume_outbound = 'true';
```

```sql+sqlite
select
  instance_alias,
  id,
  multi_party_conference,
  multi_party_chat_conference,
  forecasting_planning_scheduling,
  high_volume_outbound
from
  aws_connect_instance
where
  multi_party_conference = 'true'
  or forecasting_planning_scheduling = 'true'
  or high_volume_outbound = 'true';
```

### Check instance feature configuration
Review the complete feature configuration of Connect instances to understand their capabilities.

```sql+postgres
select
  instance_alias,
  id,
  contact_lens,
  early_media,
  contactflow_logs,
  max_package,
  use_custom_tts_voices,
  auto_resolve_best_voices
from
  aws_connect_instance
order by
  instance_alias;
```

```sql+sqlite
select
  instance_alias,
  id,
  contact_lens,
  early_media,
  contactflow_logs,
  max_package,
  use_custom_tts_voices,
  auto_resolve_best_voices
from
  aws_connect_instance
order by
  instance_alias;
```
