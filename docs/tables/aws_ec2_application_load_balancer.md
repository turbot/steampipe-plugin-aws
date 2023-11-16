---
title: "Table: aws_ec2_application_load_balancer - Query AWS EC2 Application Load Balancer using SQL"
description: "Allows users to query AWS EC2 Application Load Balancer, providing detailed information about each load balancer within an AWS account. This includes its current state, availability zones, security groups, and other important attributes."
---

# Table: aws_ec2_application_load_balancer - Query AWS EC2 Application Load Balancer using SQL

The `aws_ec2_application_load_balancer` table in Steampipe allows users to gain insights into the Application Load Balancers within their AWS EC2 service. The table provides detailed information about each Application Load Balancer, including its current state, associated security groups, availability zones, type, scheme, and other important attributes. This table can be used by DevOps engineers to query load balancer-specific details, monitor the health of the load balancers, assess load balancing configurations, and much more. The schema outlines various attributes of the Application Load Balancer, such as the ARN, DNS name, canonical hosted zone ID, and creation date, among others.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_application_load_balancer` table, you can use the `.inspect aws_ec2_application_load_balancer` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the load balancer. This provides a unique identifier for the load balancer, and can be used to join with other tables that also contain load balancer ARNs.
- `availability_zones`: Information about the Availability Zones. This is useful for understanding the geographical distribution and redundancy of your load balancers.
- `security_groups`: The IDs of the security groups for the load balancer. This is useful for querying and managing the security configurations associated with your load balancers.

## Examples

### Security group attached to the application load balancers

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as attached_security_group
from
  aws_ec2_application_load_balancer;
```


### Availability zone information

```sql
select
  name,
  az ->> 'LoadBalancerAddresses' as load_balancer_addresses,
  az ->> 'OutpostId' as outpost_id,
  az ->> 'SubnetId' as subnet_id,
  az ->> 'ZoneName' as zone_name
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(availability_zones) as az;
```


### List of application load balancers whose availability zone count is less than 1

```sql
select
  name,
  count(az ->> 'ZoneName') < 2 as zone_count_1
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```


### List of application load balancers whose logging is not enabled

```sql
select
  name,
  lb ->> 'Key' as logging_key,
  lb ->> 'Value' as logging_value
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'access_logs.s3.enabled'
  and lb ->> 'Value' = 'false';
```


### List of application load balancers whose deletion protection is not enabled

```sql
select
  name,
  lb ->> 'Key' as deletion_protection_key,
  lb ->> 'Value' as deletion_protection_value
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'deletion_protection.enabled'
  and lb ->> 'Value' = 'false';
```
