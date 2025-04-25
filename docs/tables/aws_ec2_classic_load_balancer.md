---
title: "Steampipe Table: aws_ec2_classic_load_balancer - Query AWS EC2 Classic Load Balancer using SQL"
description: "Allows users to query Classic Load Balancers within Amazon EC2."
folder: "EC2"
---

# Table: aws_ec2_classic_load_balancer - Query AWS EC2 Classic Load Balancer using SQL

The AWS EC2 Classic Load Balancer automatically distributes incoming application traffic across multiple Amazon EC2 instances in the cloud. It enables you to achieve greater levels of fault tolerance in your applications, seamlessly providing the required amount of load balancing capacity needed to distribute application traffic. This service offers a highly available, scalable, and predictable performance to distribute the workload evenly to the backend servers.

## Table Usage Guide

The `aws_ec2_classic_load_balancer` table in Steampipe provides you with information about Classic Load Balancers within Amazon Elastic Compute Cloud (EC2). This table allows you, as a cloud engineer, developer, or administrator, to query load balancer-specific details, including its availability zones, security groups, backend server descriptions, and listener descriptions. You can utilize this table to gather insights on load balancers, such as their configurations, attached instances, health checks, and more. The schema outlines the various attributes of the Classic Load Balancer for you, including the load balancer name, DNS name, created time, and associated tags.

## Examples

### Instances associated with classic load balancers
Identify the instances that are linked with classic load balancers to effectively manage and balance network traffic.

```sql+postgres
select
  name,
  instances
from
  aws_ec2_classic_load_balancer;
```

```sql+sqlite
select
  name,
  instances
from
  aws_ec2_classic_load_balancer;
```

### List of classic load balancers whose logging is not enabled
Determine the areas in which classic load balancers are operating without logging enabled. This is useful for identifying potential security gaps, as logging provides a record of all requests handled by the load balancer.

```sql+postgres
select
  name,
  access_log_enabled
from
  aws_ec2_classic_load_balancer
where
  access_log_enabled = 'false';
```

```sql+sqlite
select
  name,
  access_log_enabled
from
  aws_ec2_classic_load_balancer
where
  access_log_enabled = 'false';
```


### Security groups attached to each classic load balancer
Identify the security groups associated with each classic load balancer to ensure proper access control and minimize potential security risks.

```sql+postgres
select
  name,
  jsonb_array_elements_text(security_groups) as sg
from
  aws_ec2_classic_load_balancer;
```

```sql+sqlite
select
  name,
  json_extract(json_each.value, '$') as sg
from
  aws_ec2_classic_load_balancer,
  json_each(security_groups);
```

### Classic load balancers listener info
Uncover the details of your classic load balancer's listeners to understand how each instance is configured, including the protocols used, port numbers, SSL certificates, and any associated policy names. This information can help you manage and optimize your load balancing strategy.

```sql+postgres
select
  name,
  listener_description -> 'Listener' ->> 'InstancePort' as instance_port,
  listener_description -> 'Listener' ->> 'InstanceProtocol' as instance_protocol,
  listener_description -> 'Listener' ->> 'LoadBalancerPort' as load_balancer_port,
  listener_description -> 'Listener' ->> 'Protocol' as load_balancer_protocol,
  listener_description -> 'SSLCertificateId' ->> 'SSLCertificateId' as ssl_certificate,
  listener_description -> 'Listener' ->> 'PolicyNames' as policy_names
from
  aws_ec2_classic_load_balancer
  cross join jsonb_array_elements(listener_descriptions) as listener_description;
```

```sql+sqlite
select
  name,
  json_extract(listener_description.value, '$.Listener.InstancePort') as instance_port,
  json_extract(listener_description.value, '$.Listener.InstanceProtocol') as instance_protocol,
  json_extract(listener_description.value, '$.Listener.LoadBalancerPort') as load_balancer_port,
  json_extract(listener_description.value, '$.Listener.Protocol') as load_balancer_protocol,
  json_extract(listener_description.value, '$.SSLCertificateId.SSLCertificateId') as ssl_certificate,
  json_extract(listener_description.value, '$.Listener.PolicyNames') as policy_names
from
  aws_ec2_classic_load_balancer,
  json_each(listener_descriptions) as listener_description;
```

### Health check info
Explore the health status of your classic load balancers in AWS EC2 by analyzing parameters such as threshold values, check intervals, and timeouts. This information can be crucial for maintaining optimal server performance and minimizing downtime.

```sql+postgres
select
  name,
  healthy_threshold,
  health_check_interval,
  health_check_target,
  health_check_timeout,
  unhealthy_threshold
from
  aws_ec2_classic_load_balancer;
```

```sql+sqlite
select
  name,
  healthy_threshold,
  health_check_interval,
  health_check_target,
  health_check_timeout,
  unhealthy_threshold
from
  aws_ec2_classic_load_balancer;
```