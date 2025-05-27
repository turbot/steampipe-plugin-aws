---
title: "Steampipe Table: aws_route53_traffic_policy_instance - Query AWS Route 53 Traffic Policy Instances using SQL"
description: "Allows users to query AWS Route 53 Traffic Policy Instances, providing detailed information about each instance such as the ID, version, DNS name, and more. This table is useful for gaining insights into the configuration and status of traffic policy instances."
folder: "Route 53"
---

# Table: aws_route53_traffic_policy_instance - Query AWS Route 53 Traffic Policy Instances using SQL

The AWS Route 53 Traffic Policy Instance is a component of Amazon's scalable and highly available Domain Name System (DNS) web service. This resource allows you to manage complex DNS settings with reusable traffic policies and simplifies the process of routing internet traffic to your applications. It works by using DNS to direct incoming requests based on the routing rules defined in your traffic policy.

## Table Usage Guide

The `aws_route53_traffic_policy_instance` table in Steampipe provides you with information about Traffic Policy Instances within AWS Route 53. This table allows you, as a DevOps engineer, to query instance-specific details, including the instance ID, version, DNS name, and associated metadata. You can utilize this table to gather insights on traffic policy instances, such as their configuration, status, and more. The schema outlines the various attributes of the Traffic Policy Instance for you, including the ID, version, DNS name, and associated tags.

## Examples

### Basic Info
Explore which traffic policy instances are currently active within your AWS Route53 service. This can help you manage your DNS configurations more effectively and ensure optimal routing of your network traffic.

```sql+postgres
select
  name,
  id,
  hosted_zone_id,
  ttl,
  region
from 
  aws_route53_traffic_policy_instance;
```

```sql+sqlite
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
The query is designed to provide insights into the relationship between different instances and their associated hosted zones. This can be useful for understanding how traffic policies are being applied across various zones in a network, which can aid in network management and troubleshooting.

```sql+postgres
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

```sql+sqlite
select 
  i.name,
  i.id,
  h.id as hosted_zone_id,
  h.name as hosted_zone_name,
  h.caller_reference,
  h.private_zone
from 
  aws_route53_traffic_policy_instance as i
  join aws_route53_zone as h on i.hosted_zone_id = h.id;
```

### List associated traffic policy details for each instance
Explore the relationship between traffic policy instances and their associated traffic policies in AWS Route 53. This query can be used to understand how each instance is linked to a specific policy, providing insights into the configuration and management of traffic flow within your network.

```sql+postgres
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

```sql+sqlite
select 
  i.name,
  i.id,
  traffic_policy_id,
  p.name as traffic_policy_name,
  traffic_policy_type,
  traffic_policy_version,
  p.document
from 
  aws_route53_traffic_policy_instance as i
  join aws_route53_traffic_policy as p on i.traffic_policy_id = p.id 
  and i.traffic_policy_version = p.version;
```

### List instances that failed creation
Determine the areas in which creation of certain instances failed. This allows you to pinpoint specific locations where issues have occurred, aiding in troubleshooting and resolution.

```sql+postgres
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

```sql+sqlite
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