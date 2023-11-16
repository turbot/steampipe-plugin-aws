---
title: "Table: aws_route53_traffic_policy - Query AWS Route 53 Traffic Policies using SQL"
description: "Allows users to query AWS Route 53 Traffic Policies to retrieve information about each policy's versions, including the policy identifier, name, type, and document. This table also provides data related to the policy's associated metadata."
---

# Table: aws_route53_traffic_policy - Query AWS Route 53 Traffic Policies using SQL

The `aws_route53_traffic_policy` table in Steampipe uses the data from AWS Route 53 to provide information about each traffic policy's versions. This includes details such as the policy identifier, name, type, and document. Furthermore, it provides metadata related to the traffic policy. DevOps engineers and other users can utilize this table to query and analyze data related to traffic policies, including their versions and associated metadata. The schema outlines the various attributes of the traffic policy, including the policy ARN, version, type, document, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_traffic_policy` table, you can use the `.inspect aws_route53_traffic_policy` command in Steampipe.

**Key columns**:

- `id`: This is the identifier of the traffic policy. It is a crucial column as it uniquely identifies each traffic policy and can be used to join this table with others.
- `name`: This column represents the name of the traffic policy. It is important as it aids in the identification of specific policies when querying or joining tables.
- `type`: This column indicates the type of the traffic policy. It's useful in filtering or sorting the traffic policies based on their types during queries.


## Examples

### Basic Info

```sql
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

```sql
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

### List total policies in each dns type

```sql
select
  document ->> 'RecordType' as dns_type,
  count(id) as "policies"
from
  aws_route53_traffic_policy
group by 
  dns_type;
```