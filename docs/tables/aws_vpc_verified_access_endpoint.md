---
title: "Table: aws_vpc_verified_access_endpoint - Query AWS VPC Verified Access Endpoint using SQL"
description: "Allows users to query AWS VPC Verified Access Endpoint data, including details about the endpoint configuration, service name, and VPC ID. This information can be used to manage and secure network access to services within an AWS Virtual Private Cloud."
---

# Table: aws_vpc_verified_access_endpoint - Query AWS VPC Verified Access Endpoint using SQL

The `aws_vpc_verified_access_endpoint` table in Steampipe provides information about the verified access endpoints within AWS Virtual Private Cloud (VPC). This table allows DevOps engineers to query endpoint-specific details, including the endpoint configuration, service name, and VPC ID. Users can utilize this table to gather insights on endpoints, such as endpoint configurations, associated services, and the VPCs they belong to. The schema outlines the various attributes of the VPC verified access endpoint, including the endpoint ID, creation timestamp, attached security groups, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_verified_access_endpoint` table, you can use the `.inspect aws_vpc_verified_access_endpoint` command in Steampipe.

**Key columns**:

- `vpc_endpoint_id`: The ID of the VPC endpoint. This can be used to join this table with other tables that contain VPC endpoint information.
- `vpc_id`: The ID of the VPC where the endpoint is located. This can be used to join this table with other tables that contain VPC information.
- `service_name`: The name of the service to which the endpoint is connected. This can be used to join this table with other tables that contain service information.

## Examples

### Basic info

```sql
select
  verified_access_endpoint_id,
  verified_access_instance_id,
  verified_access_group_id,
  creation_time,
  verified_access_instance_id,
  domain_certificate_arn,
  device_validation_domain,
  status_code
from
  aws_vpc_verified_access_endpoint;
```

### List endpoints older than 30 days

```sql
select
  verified_access_endpoint_id,
  creation_time,
  description,
  status_code
from
  aws_vpc_verified_access_endpoint
where
  creation_time <= now() - interval '30' day;
```

### List endpoints that are not in active state

```sql
select
  verified_access_endpoint_id,
  status_code,
  creation_time,
  deletion_time,
  description,
  device_validation_domain
from
  aws_vpc_verified_access_endpoint
where
  status_code <> 'active';
```

### Get group details of each endpoint

```sql
select
  e.verified_access_endpoint_id,
  e.creation_time,
  g.verified_access_group_id,
  g.creation_time as group_create_time
from
  aws_vpc_verified_access_endpoint as e,
  aws_vpc_verified_access_group as g
where
  e.verified_access_group_id = g.verified_access_group_id;
```

### Get trusted provider details of each endpoint

```sql
select
  e.verified_access_group_id,
  e.creation_time,
  i.creation_time as instance_create_time,
  i.verified_access_instance_id,
  jsonb_pretty(i.verified_access_trust_providers) as verified_access_trust_providers
from
  aws_vpc_verified_access_endpoint as e,
  aws_vpc_verified_access_instance as i
where
  e.verified_access_instance_id = i.verified_access_instance_id;
```

### Count of endpoints per instance

```sql
select
  verified_access_instance_id,
  count(verified_access_endpoint_id) as instance_count
from
  aws_vpc_verified_access_endpoint
group by
  verified_access_instance_id;
```

### Get network interface details of each endpoint

```sql
select
  e.verified_access_endpoint_id,
  i.network_interface_id,
  i.interface_type,
  i.private_ip_address,
  i.association_public_ip,
  jsonb_pretty(i.groups) as security_groups
from
  aws_vpc_verified_access_endpoint as e,
  aws_ec2_network_interface as i
where
  e.network_interface_options ->> 'NetworkInterfaceId' = i.network_interface_id;
```