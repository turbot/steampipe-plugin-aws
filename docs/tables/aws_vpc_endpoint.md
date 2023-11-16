---
title: "Table: aws_vpc_endpoint - Query AWS VPC Endpoints using SQL"
description: "Allows users to query AWS VPC Endpoints and retrieve information about each endpoint's configuration, type, status, and related resources such as network interfaces, DNS entries, and security groups."
---

# Table: aws_vpc_endpoint - Query AWS VPC Endpoints using SQL

The `aws_vpc_endpoint` table in Steampipe provides information about VPC Endpoints within Amazon Virtual Private Cloud (VPC). This table allows network administrators and DevOps engineers to query endpoint-specific details, including its service configuration, type (Interface or Gateway), status, and associated resources such as network interfaces, DNS entries, and security groups. Users can utilize this table to gather insights on VPC Endpoints, such as their accessibility, security configuration, and integration with other AWS services. The schema outlines the various attributes of the VPC Endpoint, including the endpoint ID, VPC ID, service name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_endpoint` table, you can use the `.inspect aws_vpc_endpoint` command in Steampipe.

**Key columns**:

- `vpc_endpoint_id`: This is the unique identifier for the VPC endpoint. This can be used to join this table with other tables that contain VPC endpoint information.
- `vpc_id`: This is the ID of the VPC in which the endpoint is located. This can be used to join this table with other tables that contain VPC information.
- `service_name`: This is the name of the AWS service that the endpoint is configured for. This can be used to join this table with other tables that contain AWS service information.

## Examples

### List of VPC endpoint and the corresponding services

```sql
select
  vpc_endpoint_id,
  vpc_id,
  service_name
from
  aws_vpc_endpoint;
```


### Subnet Id count for each VPC endpoints

```sql
select
  vpc_endpoint_id,
  jsonb_array_length(subnet_ids) as subnet_id_count
from
  aws_vpc_endpoint;
```


### Network details for each VPC endpoint

```sql
select
  vpc_endpoint_id,
  vpc_id,
  jsonb_array_elements(subnet_ids) as subnet_ids,
  jsonb_array_elements(network_interface_ids) as network_interface_ids,
  jsonb_array_elements(route_table_ids) as route_table_ids,
  sg ->> 'GroupName' as sg_name
from
  aws_vpc_endpoint
  cross join jsonb_array_elements(groups) as sg;
```


### DNS information for the VPC endpoints

```sql
select
  vpc_endpoint_id,
  private_dns_enabled,
  dns ->> 'DnsName' as dns_name,
  dns ->> 'HostedZoneId' as hosted_zone_id
from
  aws_vpc_endpoint
  cross join jsonb_array_elements(dns_entries) as dns;
```


### VPC endpoint count by VPC ID

```sql
select
  vpc_id,
  count(vpc_endpoint_id) as vpc_endpoint_count
from
  aws_vpc_endpoint
group by
  vpc_id;
```