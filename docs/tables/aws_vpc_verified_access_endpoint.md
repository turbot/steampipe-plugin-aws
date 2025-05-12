---
title: "Steampipe Table: aws_vpc_verified_access_endpoint - Query AWS VPC Verified Access Endpoint using SQL"
description: "Allows users to query AWS VPC Verified Access Endpoint data, including details about the endpoint configuration, service name, and VPC ID. This information can be used to manage and secure network access to services within an AWS Virtual Private Cloud."
folder: "VPC"
---

# Table: aws_vpc_verified_access_endpoint - Query AWS VPC Verified Access Endpoint using SQL

The AWS VPC Verified Access Endpoint is a feature within Amazon's Virtual Private Cloud (VPC) service. It enables users to verify that the traffic leaving their VPC is coming from Amazon WorkSpaces, a managed, secure Desktop-as-a-Service (DaaS). This helps in meeting compliance requirements by providing an additional layer of security and control over the network traffic.

## Table Usage Guide

The `aws_vpc_verified_access_endpoint` table in Steampipe provides you with information about the verified access endpoints within AWS Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query endpoint-specific details, including the endpoint configuration, service name, and VPC ID. You can utilize this table to gather insights on endpoints, such as endpoint configurations, associated services, and the VPCs they belong to. The schema outlines the various attributes of the VPC verified access endpoint for you, including the endpoint ID, creation timestamp, attached security groups, and associated tags.

## Examples

### Basic info
Explore the status and validation of access endpoints within your AWS VPC to ensure security and compliance. This query allows you to identify potential vulnerabilities by examining the creation time, status, and domain certificate of each endpoint.

```sql+postgres
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

```sql+sqlite
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
Identify instances where AWS VPC verified access endpoints have not been updated in the last 30 days. This can be useful for maintaining system security and efficiency by ensuring outdated endpoints are reviewed and updated as necessary.

```sql+postgres
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

```sql+sqlite
select
  verified_access_endpoint_id,
  creation_time,
  description,
  status_code
from
  aws_vpc_verified_access_endpoint
where
  creation_time <= datetime('now', '-30 day');
```

### List endpoints that are not in active state
Discover the segments that are not currently active within your AWS VPC verified access endpoints. This is useful to identify potential issues or areas of your network that may require attention or maintenance.

```sql+postgres
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

```sql+sqlite
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
  status_code != 'active';
```

### Get group details of each endpoint
Explore the creation times and relationships between various access endpoints and groups within your AWS VPC. This can be beneficial for understanding the timeline and structure of your network security configurations.

```sql+postgres
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

```sql+sqlite
select
  e.verified_access_endpoint_id,
  e.creation_time,
  g.verified_access_group_id,
  g.creation_time as group_create_time
from
  aws_vpc_verified_access_endpoint as e
join
  aws_vpc_verified_access_group as g
on
  e.verified_access_group_id = g.verified_access_group_id;
```

### Get trusted provider details of each endpoint
Explore the trusted provider details for each endpoint to understand their creation times and associated instances. This can be particularly useful in auditing and maintaining the security of your AWS VPC.

```sql+postgres
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

```sql+sqlite
select
  e.verified_access_group_id,
  e.creation_time,
  i.creation_time as instance_create_time,
  i.verified_access_instance_id,
  i.verified_access_trust_providers as verified_access_trust_providers
from
  aws_vpc_verified_access_endpoint as e,
  aws_vpc_verified_access_instance as i
where
  e.verified_access_instance_id = i.verified_access_instance_id;
```

### Count of endpoints per instance
Determine the number of endpoints associated with each instance in your AWS Virtual Private Cloud. This can help assess the complexity and potential security exposure of your network setup.

```sql+postgres
select
  verified_access_instance_id,
  count(verified_access_endpoint_id) as instance_count
from
  aws_vpc_verified_access_endpoint
group by
  verified_access_instance_id;
```

```sql+sqlite
select
  verified_access_instance_id,
  count(verified_access_endpoint_id) as instance_count
from
  aws_vpc_verified_access_endpoint
group by
  verified_access_instance_id;
```

### Get network interface details of each endpoint
This query is useful for gaining insights into the network interface details associated with each verified access endpoint in your AWS VPC. It can help you understand the type of interface, private IP address, and associated public IP, which is beneficial for network management and troubleshooting.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```