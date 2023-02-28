# Table: aws_cloudsearch_domain

Amazon CloudSearch is a managed service in the AWS Cloud that makes it simple and cost-effective to set up, manage, and scale a search solution for your website or application.

## Examples

### Basic info

```sql
select
  domain_name,
  domain_id,
  arn,
  created,
  search_instance_type,
  search_instance_count
from
  aws_cloudsearch_domain;
```

### List domains by instance type

```sql
select
  domain_name,
  domain_id,
  arn,
  created,
  search_instance_type
from
  aws_cloudsearch_domain
where
  search_instance_type = 'search.small';
```

### Get limit details for each domain

```sql
select
  domain_name,
  domain_id,
  search_service ->> 'Endpoint' as search_service_endpoint,
  limits ->> 'MaximumPartitionCount' as maximum_partition_count,
  limits ->> 'MaximumReplicationCount' as maximum_replication_count
from
  aws_cloudsearch_domain;
```