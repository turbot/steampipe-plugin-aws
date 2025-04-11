---
title: "Steampipe Table: aws_ec2_target_group - Query AWS EC2 Target Groups using SQL"
description: "Allows users to query AWS EC2 Target Groups and provides information about each Target Group within an AWS account."
folder: "EC2"
---

# Table: aws_ec2_target_group - Query AWS EC2 Target Groups using SQL

An AWS EC2 Target Group is a component of the Elastic Load Balancing service. It is used to route requests to one or more registered targets, such as EC2 instances, as part of a load balancing configuration. This allows the distribution of network traffic to multiple resources, improving availability and fault tolerance in your applications.

## Table Usage Guide

The `aws_ec2_target_group` table in Steampipe provides you with information about each Target Group within your AWS account. This table allows you, as a DevOps engineer, security auditor, or other technical professional, to query Target Group-specific details, including the associated load balancer, health check configuration, and attributes. You can utilize this table to gather insights on Target Groups, such as their configurations, associated resources, and more. The schema outlines the various attributes of the Target Group for you, including the ARN, Health Check parameters, and associated tags.

## Examples

### Basic target group info
Explore the different target groups within your AWS EC2 instances to understand their associated load balancer resources and the virtual private cloud (VPC) they belong to. This can help in managing and optimizing your cloud resources effectively.

```sql+postgres
select
  target_group_name,
  target_type,
  load_balancer_arns,
  vpc_id
from
  aws_ec2_target_group;
```

```sql+sqlite
select
  target_group_name,
  target_type,
  load_balancer_arns,
  vpc_id
from
  aws_ec2_target_group;
```


### Health check info of target groups
This query is used to gain insights into the health check configurations of target groups within an AWS EC2 environment. Its practical application lies in its ability to help identify potential issues or vulnerabilities in the system, ensuring optimal performance and security.

```sql+postgres
select
  health_check_enabled,
  protocol,
  matcher_http_code,
  healthy_threshold_count,
  unhealthy_threshold_count,
  health_check_enabled,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  health_check_timeout_seconds
from
  aws_ec2_target_group;
```

```sql+sqlite
select
  health_check_enabled,
  protocol,
  matcher_http_code,
  healthy_threshold_count,
  unhealthy_threshold_count,
  health_check_enabled,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  health_check_timeout_seconds
from
  aws_ec2_target_group;
```


### Registered target for each target group
Determine the areas in which each registered target is located for a specific target group. This can be useful for identifying potential issues with load balancing or for optimizing resource allocation across different availability zones.

```sql+postgres
select
  target_group_name,
  target_type,
  target -> 'Target' ->> 'AvailabilityZone' as availability_zone,
  target -> 'Target' ->> 'Id' as id,
  target -> 'Target' ->> 'Port' as port
from
  aws_ec2_target_group
  cross join jsonb_array_elements(target_health_descriptions) as target;
```

```sql+sqlite
select
  target_group_name,
  target_type,
  json_extract(target.value, '$.Target.AvailabilityZone') as availability_zone,
  json_extract(target.value, '$.Target.Id') as id,
  json_extract(target.value, '$.Target.Port') as port
from
  aws_ec2_target_group,
  json_each(target_health_descriptions) as target;
```

### Health status of registered targets
Identify instances where the health status of registered targets in EC2 instances can be assessed. This allows for proactive management of resources by pinpointing potential issues or disruptions in the target groups.

```sql+postgres
select
  target_group_name,
  target_type,
  target -> 'TargetHealth' ->> 'Description' as description,
  target -> 'TargetHealth' ->> 'Reason' reason,
  target -> 'TargetHealth' ->> 'State' as state
from
  aws_ec2_target_group
  cross join jsonb_array_elements(target_health_descriptions) as target;
```

```sql+sqlite
select
  target_group_name,
  target_type,
  json_extract(target.value, '$.TargetHealth.Description') as description,
  json_extract(target.value, '$.TargetHealth.Reason') as reason,
  json_extract(target.value, '$.TargetHealth.State') as state
from
  aws_ec2_target_group,
  json_each(target_health_descriptions) as target;
```