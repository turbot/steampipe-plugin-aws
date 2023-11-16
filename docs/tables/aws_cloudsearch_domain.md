---
title: "Table: aws_cloudsearch_domain - Query AWS CloudSearch Domain using SQL"
description: "Allows users to query AWS CloudSearch Domain to retrieve detailed information about each search domain configured within an AWS account."
---

# Table: aws_cloudsearch_domain - Query AWS CloudSearch Domain using SQL

The `aws_cloudsearch_domain` table in Steampipe provides information about each search domain configured within an AWS account. This table allows DevOps engineers, data analysts, and other technical professionals to query domain-specific details, including the domain's ARN, creation date, domain ID, and associated metadata. Users can utilize this table to gather insights on domains, such as the status, endpoint, and whether the domain requires signing. The schema outlines the various attributes of the CloudSearch domain, including the domain ARN, creation date, document count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudsearch_domain` table, you can use the `.inspect aws_cloudsearch_domain` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the search domain. This is a unique identifier and can be used to join this table with other AWS tables.
- `domain_id`: The unique identifier for the domain. This can be used to cross-reference with other tables or data sources that track CloudSearch domains.
- `domain_name`: The name of the search domain. This can be useful when joining with other tables that use domain names as a reference.

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