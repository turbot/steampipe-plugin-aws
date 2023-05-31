# Table: aws_service_discovery_instance

In AWS, Service Discovery Instance refers to a resource that represents an individual instance of a service within a service discovery namespace.

AWS Service Discovery is a managed service that makes it easy to discover and locate services within your infrastructure. It enables automatic registration and DNS-based service discovery for microservices running on AWS.

By registering service instances with AWS Service Discovery, you can easily discover and access those instances using the DNS name associated with the service. This simplifies the process of locating and communicating with services within a distributed architecture, particularly in a dynamic and scalable environment.

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
