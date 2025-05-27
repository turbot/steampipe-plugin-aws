---
title: "Steampipe Table: aws_route53_traffic_policy - Query AWS Route 53 Traffic Policies using SQL"
description: "Allows users to query AWS Route 53 Traffic Policies to retrieve information about each policy's versions, including the policy identifier, name, type, and document. This table also provides data related to the policy's associated metadata."
folder: "Route 53"
---

# Table: aws_route53_traffic_policy - Query AWS Route 53 Traffic Policies using SQL

The AWS Route 53 Traffic Policy is a resource within Amazon's Route 53 service. It allows you to manage how traffic is routed to your application endpoints, enabling you to improve availability and latency. Using various routing types, like simple, failover, or geolocation routing, you can define more complex routes to optimize your application's performance.

## Table Usage Guide

The `aws_route53_traffic_policy` table in Steampipe uses data from AWS Route 53 to provide you with information about each traffic policy's versions. This includes details such as the policy identifier, name, type, and document. Furthermore, it provides you with metadata related to the traffic policy. As a DevOps engineer or other user, you can utilize this table to query and analyze data related to traffic policies, including their versions and associated metadata. The schema outlines the various attributes of the traffic policy for you, including the policy ARN, version, type, document, and more.

## Examples

### Basic Info
Explore the specifics of your AWS Route53 traffic policies, such as their names, IDs, and versions, to gain a better understanding of your network's traffic routing configurations. This can be particularly useful for managing and optimizing your network traffic flow.

```sql+postgres
select
  name,
  id,
  version,
  document,
  region
from 
  aws_route53_traffic_policy;
```

```sql+sqlite
select
  name,
  id,
  version,
  document,
  region
from 
  aws_route53_traffic_policy;
```

### List policies with latest version
Discover the segments that have the most recent versions of policies. This is useful for maintaining up-to-date policy information and ensuring compliance with the latest versions.

```sql+postgres
select 
  name,
  policy.id,
  policy.version, 
  comment 
from 
  aws_route53_traffic_policy policy,
  (select
    id,
    max(version) as version
  from 
    aws_route53_traffic_policy 
  group by 
    id) as latest
where 
  latest.id = policy.id 
  and latest.version = policy.version;
```

```sql+sqlite
select 
  name,
  policy.id,
  policy.version, 
  comment 
from 
  aws_route53_traffic_policy policy
join 
  (select
    id,
    max(version) as version
  from 
    aws_route53_traffic_policy 
  group by 
    id) as latest
on 
  latest.id = policy.id 
  and latest.version = policy.version;
```

### List total policies in each dns type
Assess the distribution of policies across different DNS types to better understand your AWS Route53 traffic policy configuration. This can help in optimizing policy management and identifying potential areas for improvement.

```sql+postgres
select
  document ->> 'RecordType' as dns_type,
  count(id) as "policies"
from
  aws_route53_traffic_policy
group by 
  dns_type;
```

```sql+sqlite
select
  json_extract(document, '$.RecordType') as dns_type,
  count(id) as "policies"
from
  aws_route53_traffic_policy
group by 
  dns_type;
```