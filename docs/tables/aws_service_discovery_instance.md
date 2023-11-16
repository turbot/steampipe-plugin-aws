---
title: "Table: aws_service_discovery_instance - Query AWS Cloud Map Service Instances using SQL"
description: "Allows users to query AWS Cloud Map Service Instances and retrieve detailed information about each instance associated with a specified service. This information includes the instance ID, instance attributes, and the health status of the instance."
---

# Table: aws_service_discovery_instance - Query AWS Cloud Map Service Instances using SQL

The `aws_service_discovery_instance` table in Steampipe provides information about service instances within AWS Cloud Map. This table allows DevOps engineers to query instance-specific details, including instance ID, attributes, and health status. Users can utilize this table to gather insights on instances, such as instances with specific attributes, health status of instances, and more. The schema outlines the various attributes of the service instance, including the instance ID, service ID, attributes, and health status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_service_discovery_instance` table, you can use the `.inspect aws_service_discovery_instance` command in Steampipe.

**Key columns**:

- `instance_id`: The ID of the instance that you want to get information about. This can be used to join this table with other tables that contain instance-specific information.
- `service_id`: The ID of the service that the instance is associated with. This can be used to join this table with other tables that contain service-specific information.
- `health_status`: The health status of the service instance. This is important as it can be used to monitor the health of the instances and take necessary action if required.

## Examples

### Basic info

```sql
select
  id,
  service_id,
  ec2_instance_id,
  attributes
from
  aws_service_discovery_instance;
```

### List instances that are unhealthy

```sql
select
  id,
  service_id,
  init_health_status
from
  aws_service_discovery_instance
where
  init_health_status = 'UNHEALTHY';
```

### Count instances by service

```sql
select
  service_id,
  count(id)
from
  aws_service_discovery_instance
group by
  service_id;
```

### Get service details of each instance

```sql
select
  i.id,
  i.service_id,
  s.name as service_name,
  s.create_date as service_create_date,
  s.namespace_id,
  s.type
from
  aws_service_discovery_instance as i,
  aws_service_discovery_service as s
where
  s.id = i.service_id;
```

### Get EC2 instance details of each service discovery instance

```sql
select
  i.id,
  i.service_id,
  i.ec2_instance_id,
  ei.instance_type,
  ei.instance_state,
  ei.launch_time
from
  aws_service_discovery_instance as i,
  aws_ec2_instance as ei
where
  i.ec2_instance_id is not null
and
  ei.instance_id = i.ec2_instance_id;
```

### Get the IP address configuration of service discovery instances 

```sql
select
  id,
  service_id,
  ec2_instance_id,
  instance_ipv4,
  instance_ipv6,
  instance_port
from
  aws_service_discovery_instance;
```
