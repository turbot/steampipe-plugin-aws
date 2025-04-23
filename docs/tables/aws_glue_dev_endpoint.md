---
title: "Steampipe Table: aws_glue_dev_endpoint - Query AWS Glue Development Endpoints using SQL"
description: "Allows users to query AWS Glue Development Endpoints to retrieve detailed information about individual endpoints, their configurations, and related metadata."
folder: "Glue"
---

# Table: aws_glue_dev_endpoint - Query AWS Glue Development Endpoints using SQL

The AWS Glue Development Endpoints are interactive programming interfaces for AWS Glue. They provide a development environment to learn, write, and test scripts that extract, transform, and load data. Using these endpoints, you can debug and test your ETL scripts before deploying them.

## Table Usage Guide

The `aws_glue_dev_endpoint` table in Steampipe provides you with comprehensive information about Development Endpoints within AWS Glue. This table allows you, as a developer or data engineer, to query endpoint-specific details, including the endpoint status, security configurations, associated subnet ID, VPC ID, and much more. You can utilize this table to analyze and manage your Glue Development Endpoints, such as identifying endpoints with specific security configurations, verifying endpoint statuses, and understanding the network configurations of the endpoints. The schema outlines the various attributes of the Glue Development Endpoint for you, including the endpoint name, role ARN, public key, creation time, and associated tags.

## Examples

### Basic info
Explore the status and availability of your AWS Glue development endpoints, including their creation timestamps, versions, and addresses. This can help you monitor the health and performance of your endpoints, ensuring they are functioning optimally and are up-to-date.

```sql+postgres
select
  endpoint_name,
  status,
  availability_zone,
  created_timestamp,
  extra_jars_s3_path,
  glue_version,
  private_address,
  public_address
from
  aws_glue_dev_endpoint;
```

```sql+sqlite
select
  endpoint_name,
  status,
  availability_zone,
  created_timestamp,
  extra_jars_s3_path,
  glue_version,
  private_address,
  public_address
from
  aws_glue_dev_endpoint;
```

### List dev endpoints that are not in ready state
Determine the areas in which development endpoints are not yet ready for use. This can aid in identifying potential issues or bottlenecks in the system.

```sql+postgres
select
  endpoint_name,
  status,
  created_timestamp,
  extra_jars_s3_path,
  glue_version,
  private_address,
  public_address
from
  aws_glue_dev_endpoint
where
  status <> 'READY'; 
```

```sql+sqlite
select
  endpoint_name,
  status,
  created_timestamp,
  extra_jars_s3_path,
  glue_version,
  private_address,
  public_address
from
  aws_glue_dev_endpoint
where
  status <> 'READY';
```

### List dev endpoints updated in the last 30 days
Discover the segments that have seen recent modifications in your development endpoints. This is particularly useful to track changes and stay updated with the latest modifications made within the past month.

```sql+postgres
select
  title,
  arn,
  status,
  glue_version,
  last_modified_timestamp
from
  aws_glue_dev_endpoint
where
   last_modified_timestamp >= now() - interval '30' day;
```

```sql+sqlite
select
  title,
  arn,
  status,
  glue_version,
  last_modified_timestamp
from
  aws_glue_dev_endpoint
where
   last_modified_timestamp >= datetime('now','-30 day');
```

### List dev endpoints older than 30 days
Determine the areas in which development endpoints have been active for more than 30 days. This can be useful for understanding long-term usage patterns and identifying potential areas for optimization or resource reallocation.

```sql+postgres
select
  endpoint_name,
  arn,
  status,
  glue_version,
  created_timestamp
from
  aws_glue_dev_endpoint
where
   created_timestamp >= now() - interval '30' day;
```

```sql+sqlite
select
  endpoint_name,
  arn,
  status,
  glue_version,
  created_timestamp
from
  aws_glue_dev_endpoint
where
   created_timestamp >= datetime('now','-30 day');
```

### Get subnet details attached to a particular dev endpoint
Explore the specifics of a particular development endpoint, such as the availability zone and IP address count, to gain insights into its configuration and status. This is particularly useful for managing network resources and optimizing system performance.

```sql+postgres
select
  e.endpoint_name,
  s.availability_zone,
  s.available_ip_address_count,
  s.cidr_block,
  s.default_for_az,
  s.map_customer_owned_ip_on_launch,
  s.map_public_ip_on_launch,
  s.state
from
  aws_glue_dev_endpoint as e,
  aws_vpc_subnet as s
where
  e.endpoint_name = 'test5'
and
  e.subnet_id = s.subnet_id;
```

```sql+sqlite
select
  e.endpoint_name,
  s.availability_zone,
  s.available_ip_address_count,
  s.cidr_block,
  s.default_for_az,
  s.map_customer_owned_ip_on_launch,
  s.map_public_ip_on_launch,
  s.state
from
  aws_glue_dev_endpoint as e
join
  aws_vpc_subnet as s
on
  e.subnet_id = s.subnet_id
where
  e.endpoint_name = 'test5';
```

### Get extra jars s3 bucket details for a dev endpoint 
Determine the configuration details of specific S3 buckets that are linked to a development endpoint in AWS Glue. This is useful for assessing the versioning status, policy, and object lock configuration of these buckets, aiding in security and management tasks.

```sql+postgres
select
  e.endpoint_name,
  split_part(j, '/', '3') as extra_jars_s3_bucket,
  b.versioning_enabled,
  b.policy,
  b.object_lock_configuration,
  b.restrict_public_buckets,
  b.policy
from
  aws_glue_dev_endpoint as e,
  aws_s3_bucket as b,
  unnest (string_to_array(e.extra_jars_s3_path, ',')) as j
where
  b.name = split_part(j, '/', '3')
and
  e.endpoint_name = 'test34';
```

```sql+sqlite
Error: SQLite does not support the unnest, split_part, or string_to_array functions.
```