---
title: "Steampipe Table: aws_service_discovery_namespace - Query AWS Cloud Map Service Discovery Namespace using SQL"
description: "Allows users to query AWS Cloud Map Service Discovery Namespace to retrieve details about the namespaces in AWS Cloud Map."
folder: "Service Discovery"
---

# Table: aws_service_discovery_namespace - Query AWS Cloud Map Service Discovery Namespace using SQL

The AWS Cloud Map Service Discovery Namespace is a component of AWS Cloud Map that helps applications to discover services dynamically over the cloud. It allows services to register their instance and utilize naming schema for easy discovery. This aids in maintaining the updated location of services, which is crucial for microservice architectures and serverless applications.

## Table Usage Guide

The `aws_service_discovery_namespace` table in Steampipe provides you with information about AWS Cloud Map Service Discovery Namespaces. This table allows you, as a DevOps engineer, to query namespace-specific details, including namespace type (DNS, HTTP), associated services, and associated metadata. You can utilize this table to gather insights on namespaces, such as the number of services in each namespace, namespace types, and more. The schema outlines the various attributes of the service discovery namespace for you, including the namespace ID, ARN, name, type, and associated tags.

## Examples

### Basic info
Determine the areas in which AWS services are being utilized, by identifying their names, IDs, types, and regions. This is useful for understanding your AWS service usage and distribution across different regions.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  region
from
  aws_service_discovery_namespace;
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  region
from
  aws_service_discovery_namespace;
```

### List private namespaces
Explore which namespaces are private within your AWS Service Discovery to better manage your resources and ensure data security. This is particularly useful in scenarios where you need to isolate specific resources or data within your AWS environment.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  service_count
from
  aws_service_discovery_namespace
where
  type ilike '%private%';
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  service_count
from
  aws_service_discovery_namespace
where
  type like '%private%';
```

### List HTTP type namespaces
Identify instances where the type of AWS Service Discovery Namespace is HTTP. This can help in understanding the distribution of different namespace types and aid in the management of AWS services.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  service_count
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  service_count
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

### List namespaces created in the last 30 days
Discover the segments that have been recently added to the service discovery namespace in AWS within the past month. This can be useful in tracking recent changes or additions to your AWS environment.

```sql+postgres
select
  name,
  id,
  description,
  create_date
from
  aws_service_discovery_namespace
where
  create_date >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  id,
  description,
  create_date
from
  aws_service_discovery_namespace
where
  create_date >= datetime('now', '-30 day');
```

### Get HTTP property details of namespaces
Explore the details of HTTP properties associated with specific namespaces to better understand their configuration and usage in your AWS service discovery setup. This can be useful in diagnosing issues or optimizing resource utilization.

```sql+postgres
select
  name,
  id,
  http_properties ->> 'HttpName' as http_name
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

```sql+sqlite
select
  name,
  id,
  json_extract(http_properties, '$.HttpName') as http_name
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

### Get private DNS property details of namspaces
Determine the areas in which private DNS properties of namespaces are utilized within AWS Service Discovery. This assists in understanding the specific configuration and settings of your private DNS, thus aiding in better resource management and optimization.

```sql+postgres
select
  name,
  id,
  dns_properties ->> 'HostedZoneId' as HostedZoneId,
  dns_properties -> 'SOA' ->> 'TTL' as ttl
from
  aws_service_discovery_namespace
where
  type = 'DNS_PRIVATE';
```

```sql+sqlite
select
  name,
  id,
  json_extract(dns_properties, '$.HostedZoneId') as HostedZoneId,
  json_extract(json_extract(dns_properties, '$.SOA'), '$.TTL') as ttl
from
  aws_service_discovery_namespace
where
  type = 'DNS_PRIVATE';
```

### Count namespaces by type
Analyze the distribution of different types of namespaces within AWS service discovery to understand their usage and prevalence.

```sql+postgres
select
  type,
  count(type)
from
  aws_service_discovery_namespace
group by
  type;
```

```sql+sqlite
select
  type,
  count(type)
from
  aws_service_discovery_namespace
group by
  type;
```