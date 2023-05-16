# Table: aws_service_discovery_namespace

AWS Service Discovery Namespace refers to a logical group or container for registering and discovering services within the AWS ecosystem. It is a component of the AWS Service Discovery service that allows you to easily manage service discovery for your applications running on AWS.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  type,
  region
from
  aws_service_discovery_namespace;
```

### List HTTP type namespaces

```sql
select
  name,
  id,
  arn,
  type,
  service_account
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

### List portfolios created in the last 30 days

```sql
name
  display_name,
  id,
  description,
  create_date
from
  aws_service_discovery_namespace
where
  create_date >= now() - interval '30' day;
```

### Get HTTP property for the namespace type HTTP

```sql
select
  name,
  id,
  http_properties ->> 'HttpName' as http_name
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

### Get private DNS property for the namespace type private DNS

```sql
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

### Count namespaces by type

```sql
select
  type,
  count(type)
from
  aws_service_discovery_namespace
group by
  type;
```