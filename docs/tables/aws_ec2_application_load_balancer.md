---
title: "Steampipe Table: aws_ec2_application_load_balancer - Query AWS EC2 Application Load Balancer using SQL"
description: "Allows users to query AWS EC2 Application Load Balancer, providing detailed information about each load balancer within an AWS account. This includes its current state, availability zones, security groups, and other important attributes."
folder: "EC2"
---

# Table: aws_ec2_application_load_balancer - Query AWS EC2 Application Load Balancer using SQL

The AWS EC2 Application Load Balancer is a resource within Amazon's Elastic Compute Cloud (EC2) service that automatically distributes incoming application traffic across multiple targets, such as EC2 instances, in multiple Availability Zones. This enhances the fault tolerance of your applications. The load balancer serves as a single point of contact for clients, which increases the availability of your application.

## Table Usage Guide

The `aws_ec2_application_load_balancer` table in Steampipe allows you to gain insights into the Application Load Balancers within your AWS EC2 service. The table provides detailed information about each Application Load Balancer, including its current state, associated security groups, availability zones, type, scheme, and other important attributes. You can use this table to query load balancer-specific details, monitor the health of the load balancers, assess load balancing configurations, and much more. The schema outlines various attributes of the Application Load Balancer, such as the ARN, DNS name, canonical hosted zone ID, and creation date, among others.

## Examples

### Security group attached to the application load balancers
Explore which security groups are linked to your application load balancers, enabling you to assess potential vulnerabilities and ensure optimal security configurations. This can be particularly useful for identifying security loopholes and reinforcing your system's defenses.

```sql+postgres
select
  name,
  jsonb_array_elements_text(security_groups) as attached_security_group
from
  aws_ec2_application_load_balancer;
```

```sql+sqlite
select
  name,
  json_extract(json_each.value, '$') as attached_security_group
from
  aws_ec2_application_load_balancer,
  json_each(security_groups);
```

### Availability zone information
Discover the segments that provide insights into the availability zones of your AWS EC2 application load balancer. This can be particularly useful for understanding your load balancer's distribution and identifying potential areas for improvement or troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(az.value, '$.LoadBalancerAddresses') as load_balancer_addresses,
  json_extract(az.value, '$.OutpostId') as outpost_id,
  json_extract(az.value, '$.SubnetId') as subnet_id,
  json_extract(az.value, '$.ZoneName') as zone_name
from
  aws_ec2_application_load_balancer,
  json_each(availability_zones) as az;
```

### List of application load balancers whose availability zone count is less than 1
Explore which application load balancers are potentially at risk due to being located in less than two availability zones. This is useful for identifying weak points in your infrastructure and improving system resilience.

```sql+postgres
select
  name,
  count(az ->> 'ZoneName') < 2 as zone_count_1
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```

```sql+sqlite
select
  name,
  count(json_extract(az.value, '$.ZoneName')) < 2 as zone_count_1
from
  aws_ec2_application_load_balancer,
  json_each(availability_zones) as az
group by
  name;
```


### List of application load balancers whose logging is not enabled
Identify instances where application load balancers in your AWS EC2 environment have their logging feature disabled. This is useful for maintaining security and compliance by ensuring all load balancers are properly recording activity.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as logging_key,
  json_extract(lb.value, '$.Value') as logging_value
from
  aws_ec2_application_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'access_logs.s3.enabled'
  and json_extract(lb.value, '$.Value') = 'false';
```


### List of application load balancers whose deletion protection is not enabled
Identify instances where application load balancers are not safeguarded against unintended deletion. This information can be useful in ensuring system resilience and minimizing service disruptions.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as deletion_protection_key,
  json_extract(lb.value, '$.Value') as deletion_protection_value
from
  aws_ec2_application_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'deletion_protection.enabled'
  and json_extract(lb.value, '$.Value') = 'false';
```