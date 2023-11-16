---
title: "Table: aws_ec2_classic_load_balancer - Query AWS EC2 Classic Load Balancer using SQL"
description: "Allows users to query Classic Load Balancers within Amazon EC2."
---

# Table: aws_ec2_classic_load_balancer - Query AWS EC2 Classic Load Balancer using SQL

The `aws_ec2_classic_load_balancer` table in Steampipe provides information about Classic Load Balancers within Amazon Elastic Compute Cloud (EC2). This table allows cloud engineers, developers, and administrators to query load balancer-specific details, including its availability zones, security groups, backend server descriptions, and listener descriptions. Users can utilize this table to gather insights on load balancers, such as their configurations, attached instances, health checks, and more. The schema outlines the various attributes of the Classic Load Balancer, including the load balancer name, DNS name, created time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_classic_load_balancer` table, you can use the `.inspect aws_ec2_classic_load_balancer` command in Steampipe.

### Key columns:

- `load_balancer_name`: The name of the load balancer. This can be used to join with other tables where the load balancer name is required.
- `availability_zones`: The availability zones for the load balancer. This is useful for joining with tables that need to understand the geographical distribution of resources.
- `security_groups`: The security groups assigned to the load balancer's network interface. It is useful for joining with security group tables to understand the security posture of the load balancer.

## Examples

### Instances associated with classic load balancers

```sql
select
  name,
  instances
from
  aws_ec2_classic_load_balancer;
```


### List of classic load balancers whose logging is not enabled

```sql
select
  name,
  access_log_enabled
from
  aws_ec2_classic_load_balancer
where
  access_log_enabled = 'false';
```


### Security groups attached to each classic load balancer

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as sg
from
  aws_ec2_classic_load_balancer;
```


### Classic load balancers listener info

```sql
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


### Health check info

```sql
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