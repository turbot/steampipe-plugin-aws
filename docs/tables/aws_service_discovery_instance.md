---
title: "Steampipe Table: aws_service_discovery_instance - Query AWS Cloud Map Service Instances using SQL"
description: "Allows users to query AWS Cloud Map Service Instances and retrieve detailed information about each instance associated with a specified service. This information includes the instance ID, instance attributes, and the health status of the instance."
folder: "Service Discovery"
---

# Table: aws_service_discovery_instance - Query AWS Cloud Map Service Instances using SQL

The AWS Cloud Map Service Instance is a component of AWS Cloud Map that allows you to register any application component, such as databases, queues, microservices, or other services, and manage their location. It provides a unified view of operational data for all your services and ensures they are connected in a reliable and scalable manner. With AWS Cloud Map, you can define custom names for your application resources, and it maintains the updated location of these dynamically changing resources.

## Table Usage Guide

The `aws_service_discovery_instance` table in Steampipe provides you with information about service instances within AWS Cloud Map. This table allows you, as a DevOps engineer, to query instance-specific details, including instance ID, attributes, and health status. You can utilize this table to gather insights on instances, such as instances with specific attributes, health status of instances, and more. The schema outlines the various attributes of the service instance for you, including the instance ID, service ID, attributes, and health status.

## Examples

### Basic info
Explore the relationship between different services and instances in your AWS environment. This query helps in tracking the associations and attributes of services, which is beneficial for managing resources and troubleshooting issues.

```sql+postgres
select
  id,
  service_id,
  ec2_instance_id,
  attributes
from
  aws_service_discovery_instance;
```

```sql+sqlite
select
  id,
  service_id,
  ec2_instance_id,
  attributes
from
  aws_service_discovery_instance;
```

### List instances that are unhealthy
Identify instances where the initial health status is marked as unhealthy. This can help in quickly pinpointing problematic areas within your AWS services, enabling you to take corrective measures promptly.

```sql+postgres
select
  id,
  service_id,
  init_health_status
from
  aws_service_discovery_instance
where
  init_health_status = 'UNHEALTHY';
```

```sql+sqlite
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
Gain insights into the distribution of instances across various services, allowing you to understand the service utilization patterns. This can be particularly useful for load balancing and resource allocation.

```sql+postgres
select
  service_id,
  count(id)
from
  aws_service_discovery_instance
group by
  service_id;
```

```sql+sqlite
select
  service_id,
  count(id)
from
  aws_service_discovery_instance
group by
  service_id;
```

### Get service details of each instance
Gain insights into the specifics of each service associated with every instance. This is useful for understanding the relationships between services and instances, and for tracking service creation dates.

```sql+postgres
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

```sql+sqlite
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
Explore the specifics of each service discovery instance by examining the associated EC2 instance details. This allows for a comprehensive understanding of the service's operation and usage, providing valuable insights for optimization and management.

```sql+postgres
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

```sql+sqlite
select
  i.id,
  i.service_id,
  i.ec2_instance_id,
  ei.instance_type,
  ei.instance_state,
  ei.launch_time
from
  aws_service_discovery_instance as i
join
  aws_ec2_instance as ei
on
  ei.instance_id = i.ec2_instance_id
where
  i.ec2_instance_id is not null;
```

### Get the IP address configuration of service discovery instances 
Determine the IP address configuration of instances within a service discovery setup to better understand networking and connectivity details. This can be helpful in troubleshooting network-related issues or planning for network expansion.

```sql+postgres
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

```sql+sqlite
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