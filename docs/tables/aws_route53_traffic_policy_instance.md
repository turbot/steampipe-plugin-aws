---
title: "Table: aws_route53_traffic_policy_instance - Query AWS Route 53 Traffic Policy Instances using SQL"
description: "Allows users to query AWS Route 53 Traffic Policy Instances, providing detailed information about each instance such as the ID, version, DNS name, and more. This table is useful for gaining insights into the configuration and status of traffic policy instances."
---

# Table: aws_route53_traffic_policy_instance - Query AWS Route 53 Traffic Policy Instances using SQL

The `aws_route53_traffic_policy_instance` table in Steampipe provides information about Traffic Policy Instances within AWS Route 53. This table allows DevOps engineers to query instance-specific details, including the instance ID, version, DNS name, and associated metadata. Users can utilize this table to gather insights on traffic policy instances, such as their configuration, status, and more. The schema outlines the various attributes of the Traffic Policy Instance, including the ID, version, DNS name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_traffic_policy_instance` table, you can use the `.inspect aws_route53_traffic_policy_instance` command in Steampipe.

**Key columns**:

- `id`: The ID of the traffic policy instance. This can be used to link this table with other tables that require a traffic policy instance ID.
- `version`: The version number of the traffic policy that is associated with the current instance. This is useful for tracking changes and updates to the traffic policy.
- `dns_name`: The DNS name that is associated with the traffic policy instance. This can be used to link this table with other tables that require a DNS name.

## Examples

### Basic Info

```sql
select
  name,
  id,
  hosted_zone_id,
  ttl,
  region
from 
  aws_route53_traffic_policy_instance;
```

### List associated hosted zone details for each instance

```sql
select 
  i.name,
  i.id,
  h.id as hosted_zone_id,
  h.name as hosted_zone_name,
  h.caller_reference,
  h.private_zone
from 
  aws_route53_traffic_policy_instance i
  join aws_route53_zone h on i.hosted_zone_id = h.id;
```

### List associated traffic policy details for each instance

```sql
select 
  i.name,
  i.id,
  traffic_policy_id,
  p.name as traffic_policy_name,
  traffic_policy_type,
  traffic_policy_version,
  p.document
from 
  aws_route53_traffic_policy_instance i
  join aws_route53_traffic_policy p on i.traffic_policy_id = p.id 
  and i.traffic_policy_version = p.version;
```

### List instances that failed creation

```sql
select
  name,
  id,
  state,
  hosted_zone_id,
  message as failed_reason
from 
  aws_route53_traffic_policy_instance
where
  state = 'Failed';
```