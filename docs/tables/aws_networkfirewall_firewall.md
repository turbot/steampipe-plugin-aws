---
title: "Steampipe Table: aws_networkfirewall_firewall - Query AWS Network Firewall using SQL"
description: "Allows users to query AWS Network Firewalls for detailed information about each firewall's properties and settings."
folder: "Network Firewall"
---

# Table: aws_networkfirewall_firewall - Query AWS Network Firewall using SQL

The AWS Network Firewall is a managed service that makes it easy to deploy essential network protections for all of your Amazon Virtual Private Clouds (VPCs). The service can be set up, configured, and maintained through a simple console interface, APIs, or with the AWS CLI. It provides high availability, scalability, and you pay only for what you use.

## Table Usage Guide

The `aws_networkfirewall_firewall` table in Steampipe provides you with information about each firewall in AWS Network Firewall. This table enables you, as a network administrator, security analyst, or DevOps engineer, to query specific details about firewalls, including firewall policies, subnet mappings, and associated VPCs. You can utilize this table to gain insights into firewall configurations, such as firewall policy ARNs, VPC IDs, subnet IDs, and more. The schema outlines the various attributes of the firewall for you, including the firewall ARN, firewall name, firewall policy ARN, VPC ID, subnet mapping, delete protection status, and associated tags.

## Examples

### Basic info
Determine the areas in which your AWS Network Firewall is deployed to gain insights into the regions and associated VPCs. This can help you assess your firewall's coverage and ensure your resources are adequately protected.

```sql+postgres
select
  arn,
  name,
  description,
  vpc_id,
  policy_arn,
  region,
  tags
from
  aws_networkfirewall_firewall;
```

```sql+sqlite
select
  arn,
  name,
  description,
  vpc_id,
  policy_arn,
  region,
  tags
from
  aws_networkfirewall_firewall;
```

### List firewalls using default encryption
Determine the areas in which firewalls are utilizing the default encryption settings. This is useful for identifying potential security vulnerabilities and ensuring compliance with best practices.

```sql+postgres
select
  arn,
  name,
  description,
  encryption_configuration
from
  aws_networkfirewall_firewall
where
  encryption_configuration ->> 'Type' = 'AWS_OWNED_KMS_KEY';
```

```sql+sqlite
select
  arn,
  name,
  description,
  encryption_configuration
from
  aws_networkfirewall_firewall
where
  json_extract(encryption_configuration, '$.Type') = 'AWS_OWNED_KMS_KEY';
```

### List firewalls having deletion protection disabled
Discover the segments of your network that are potentially vulnerable due to firewalls with deletion protection disabled. This is beneficial in enhancing your security measures by identifying and rectifying areas of weakness within your network infrastructure.

```sql+postgres
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not delete_protection;
```

```sql+sqlite
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  delete_protection = 0;
```

### List firewalls having policy change protection disabled
Discover firewalls where policy change protection is disabled to identify potential security vulnerabilities in your network. This can help in prioritizing and addressing security loopholes to maintain a robust defense system.

```sql+postgres
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not policy_change_protection;
```

```sql+sqlite
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  policy_change_protection = 0;
```

### List firewalls having subnet change protection disabled
Explore which firewalls lack protection against subnet changes. This is beneficial in identifying potential security vulnerabilities within your network infrastructure.

```sql+postgres
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not subnet_change_protection;
```

```sql+sqlite
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  subnet_change_protection = 0;
```

### Get subnet details for each firewall
This query is useful to understand the relationship between your firewalls and subnets in your network. It helps you identify the specific locations where your firewalls are deployed, providing insights into your network's security infrastructure.

```sql+postgres
select
  f.arn,
  f.name,
  f.vpc_id,
  s ->> 'SubnetId' as subnet_id,
  cidr_block,
  availability_zone,
  default_for_az
from
  aws_networkfirewall_firewall f,
  jsonb_array_elements(subnet_mappings) s,
  aws_vpc_subnet vs
where
  vs.subnet_id = s ->> 'SubnetId';
```

```sql+sqlite
select
  f.arn,
  f.name,
  f.vpc_id,
  json_extract(s.value, '$.SubnetId') as subnet_id,
  cidr_block,
  availability_zone,
  default_for_az
from
  aws_networkfirewall_firewall f,
  json_each(f.subnet_mappings) as s,
  aws_vpc_subnet vs
where
  vs.subnet_id = json_extract(s.value, '$.SubnetId');
```

### Get KMS key details of firewalls encrypted with customer managed keys
Identify firewalls that are encrypted with customer-managed keys and gain insights into their key rotation status. This can be useful in ensuring that your organization's encryption practices are in line with its security policies.

```sql+postgres
select
  f.arn,
  f.name,
  f.vpc_id,
  k.arn as key_arn,
  key_rotation_enabled
from
  aws_networkfirewall_firewall f,
  aws_kms_key k
where
  k.id = encryption_configuration ->> 'KeyId'
  and not f.encryption_configuration ->> 'Type' = 'AWS_OWNED_KMS_KEY';
```

```sql+sqlite
select
  f.arn,
  f.name,
  f.vpc_id,
  k.arn as key_arn,
  key_rotation_enabled
from
  aws_networkfirewall_firewall f,
  aws_kms_key k
where
  k.id = json_extract(f.encryption_configuration, '$.KeyId')
  and not json_extract(f.encryption_configuration, '$.Type') = 'AWS_OWNED_KMS_KEY';
```

### Get logging configuration details of firewall
The detailed insight into log types and destinations aids in ensuring that the network firewall configurations comply with organizational policies and regulatory standards. This is essential for audits, where evidence of proper log management practices needs to be presented.

```sql+postgres
select
  name,
  arn,
  l -> 'LogDestination' as log_destination,
  l ->> 'LogDestinationType' as log_destination_type,
  l ->> 'LogType' as log_type
from
  aws_networkfirewall_firewall,
  jsonb_array_elements(logging_configuration) as l;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(l.value, '$.LogDestination') as log_destination,
  json_extract(l.value, '$.LogDestinationType') as log_destination_type,
  json_extract(l.value, '$.LogType') as log_type
from
  aws_networkfirewall_firewall,
  json_each(aws_networkfirewall_firewall.logging_configuration) as l;
```