---
title: "Table: aws_glue_dev_endpoint - Query AWS Glue Development Endpoints using SQL"
description: "Allows users to query AWS Glue Development Endpoints to retrieve detailed information about individual endpoints, their configurations, and related metadata."
---

# Table: aws_glue_dev_endpoint - Query AWS Glue Development Endpoints using SQL

The `aws_glue_dev_endpoint` table in Steampipe provides comprehensive information about Development Endpoints within AWS Glue. This table allows developers and data engineers to query endpoint-specific details, including the endpoint status, security configurations, associated subnet ID, VPC ID, and much more. Users can utilize this table to analyze and manage their Glue Development Endpoints, such as identifying endpoints with specific security configurations, verifying endpoint statuses, and understanding the network configurations of the endpoints. The schema outlines the various attributes of the Glue Development Endpoint, including the endpoint name, role ARN, public key, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_dev_endpoint` table, you can use the `.inspect aws_glue_dev_endpoint` command in Steampipe.

### Key columns:

- `name`: The name of the development endpoint. This is the unique identifier for the endpoint and can be used to join this table with other tables that reference Glue Development Endpoints.
- `role_arn`: The ARN of the IAM role used in the creation of the development endpoint. This can be useful for joining with IAM tables to gather more information about the permissions and policies associated with the endpoint.
- `vpc_id`: The ID of the VPC associated with the development endpoint. This can be used to join with VPC tables to understand the network configurations and security settings of the VPC where the endpoint resides.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

### List dev endpoints older than 30 days

```sql
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

### Get subnet details attached to a particular dev endpoint

```sql
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

### Get extra jars s3 bucket details for a dev endpoint 

```sql
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