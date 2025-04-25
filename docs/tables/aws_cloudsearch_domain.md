---
title: "Steampipe Table: aws_cloudsearch_domain - Query AWS CloudSearch Domain using SQL"
description: "Allows users to query AWS CloudSearch Domain to retrieve detailed information about each search domain configured within an AWS account."
folder: "CloudSearch"
---

# Table: aws_cloudsearch_domain - Query AWS CloudSearch Domain using SQL

The AWS CloudSearch Domain is a component of AWS CloudSearch, a fully-managed service that makes it easy to set up, manage, and scale a search solution for your website or application. AWS CloudSearch features include indexing of data, running search queries, and updating the search index. It provides a high level of flexibility and scalability, allowing you to search large collections of data efficiently.

## Table Usage Guide

The `aws_cloudsearch_domain` table in Steampipe provides you with information about each search domain configured within your AWS account. This table allows you, as a DevOps engineer, data analyst, or other technical professional, to query domain-specific details, including the domain's ARN, creation date, domain ID, and associated metadata. You can utilize this table to gather insights on domains, such as the status, endpoint, and whether the domain requires signing. The schema outlines the various attributes of the CloudSearch domain for you, including the domain ARN, creation date, document count, and associated tags.

## Examples

### Basic info
Explore the basic information about your AWS CloudSearch domains, such as when they were created and the type and count of search instances. This can help manage resources and assess the capacity and usage of your search domains.

```sql+postgres
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

```sql+sqlite
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
Identify instances where specific domains are linked to a certain type of search instance. This can be useful to understand the spread and usage of different search instances across your domains.

```sql+postgres
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

```sql+sqlite
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
Explore the limits set for each domain in your AWS CloudSearch to understand how it may impact the performance and availability of your search service. This can help in optimizing the search service configuration for better resource management.

```sql+postgres
select
  domain_name,
  domain_id,
  search_service ->> 'Endpoint' as search_service_endpoint,
  limits ->> 'MaximumPartitionCount' as maximum_partition_count,
  limits ->> 'MaximumReplicationCount' as maximum_replication_count
from
  aws_cloudsearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(search_service, '$.Endpoint') as search_service_endpoint,
  json_extract(limits, '$.MaximumPartitionCount') as maximum_partition_count,
  json_extract(limits, '$.MaximumReplicationCount') as maximum_replication_count
from
  aws_cloudsearch_domain;
```