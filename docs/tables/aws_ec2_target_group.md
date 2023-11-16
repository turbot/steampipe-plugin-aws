---
title: "Table: aws_ec2_target_group - Query AWS EC2 Target Groups using SQL"
description: "Allows users to query AWS EC2 Target Groups and provides information about each Target Group within an AWS account."
---

# Table: aws_ec2_target_group - Query AWS EC2 Target Groups using SQL

The `aws_ec2_target_group` table in Steampipe provides information about each Target Group within an AWS account. This table allows DevOps engineers, security auditors, and other technical professionals to query Target Group-specific details, including the associated load balancer, health check configuration, and attributes. Users can utilize this table to gather insights on Target Groups, such as their configurations, associated resources, and more. The schema outlines the various attributes of the Target Group, including the ARN, Health Check parameters, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_target_group` table, you can use the `.inspect aws_ec2_target_group` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the target group. This can be used to join with other tables where a target group ARN is referenced.
- `load_balancer_arns`: The ARNs of the load balancers that route traffic to this target group. This can be used to correlate with other tables that reference Load Balancer ARNs.
- `target_type`: The type of targets registered with the target group. This can be used to filter target groups based on their target type (instance, ip, or lambda).

## Examples

### Basic target group info

```sql
select
  target_group_name,
  target_type,
  load_balancer_arns,
  vpc_id
from
  aws_ec2_target_group;
```


### Health check info of target groups

```sql
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

```sql
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


### Health status of registered targets

```sql
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
