---
title: "Steampipe Table: aws_emr_studio - Query AWS EMR Studio using SQL"
description: "Allows users to query AWS EMR Studio for detailed information about the configuration, security settings, and other metadata of each studio."
folder: "EMR"
---

# Table: aws_emr_studio - Query AWS EMR Studio using SQL

AWS EMR Studio is an integrated development environment (IDE) that makes it easy for data scientists and data engineers to develop, visualize, and debug data engineering and data science applications written in R, Python, Scala, and SQL. EMR Studio provides fully managed Jupyter notebooks and tools to help you develop and debug applications.

## Table Usage Guide

The `aws_emr_studio` table in Steampipe provides you with information about EMR Studios within AWS Elastic MapReduce (EMR). This table allows you, as a DevOps engineer or data engineer, to query studio-specific details, including authentication settings, security configurations, and associated metadata. You can utilize this table to gather insights on studios, such as their authentication mode, security group configurations, and network settings. The schema outlines the various attributes of the EMR Studio for you, including the studio ID, name, ARN, and associated security and network configurations.

## Examples

### Basic info
Explore the basic information about your EMR Studios, including their names, IDs, and authentication modes. This can help you understand the configuration of your studios at a glance.

```sql+postgres
select
  name,
  studio_id,
  auth_mode,
  url,
  creation_time
from
  aws_emr_studio;
```

```sql+sqlite
select
  name,
  studio_id,
  auth_mode,
  url,
  creation_time
from
  aws_emr_studio;
```

### List studios with IAM authentication
Identify studios that use IAM for authentication. This can be useful for understanding your authentication setup and ensuring proper security configurations.

```sql+postgres
select
  name,
  studio_id,
  auth_mode,
  service_role,
  user_role
from
  aws_emr_studio
where
  auth_mode = 'IAM';
```

```sql+sqlite
select
  name,
  studio_id,
  auth_mode,
  service_role,
  user_role
from
  aws_emr_studio
where
  auth_mode = 'IAM';
```

### Get security group details for studios
Explore the security group configurations for your EMR Studios to ensure proper network security settings are in place.

```sql+postgres
select
  name,
  studio_id,
  engine_security_group_id,
  workspace_security_group_id,
  vpc_id
from
  aws_emr_studio;
```

```sql+sqlite
select
  name,
  studio_id,
  engine_security_group_id,
  workspace_security_group_id,
  vpc_id
from
  aws_emr_studio;
```

### List studios with specific VPC
Find studios that are associated with a particular VPC to understand your network configuration and resource organization.

```sql+postgres
select
  name,
  studio_id,
  vpc_id,
  subnet_ids
from
  aws_emr_studio
where
  vpc_id = 'vpc-1234567890abcdef0';
```

```sql+sqlite
select
  name,
  studio_id,
  vpc_id,
  subnet_ids
from
  aws_emr_studio
where
  vpc_id = 'vpc-1234567890abcdef0';
```

### Get identity provider configuration details
Examine the identity provider settings for studios using IAM Identity Center for authentication.

```sql+postgres
select
  name,
  studio_id,
  auth_mode,
  idp_auth_url,
  idp_relay_state_parameter_name
from
  aws_emr_studio
where
  auth_mode = 'SSO';
```

```sql+sqlite
select
  name,
  studio_id,
  auth_mode,
  idp_auth_url,
  idp_relay_state_parameter_name
from
  aws_emr_studio
where
  auth_mode = 'SSO';
```

### List studios with default S3 location configuration
Find studios that have been configured with a default S3 location for workspace backups and notebook files.

```sql+postgres
select
  name,
  studio_id,
  default_s3_location,
  creation_time
from
  aws_emr_studio
where
  default_s3_location is not null;
```

```sql+sqlite
select
  name,
  studio_id,
  default_s3_location,
  creation_time
from
  aws_emr_studio
where
  default_s3_location is not null;
```

### Find studios with missing security configurations
Identify studios that might be missing critical security configurations, such as security groups or VPC settings.

```sql+postgres
select
  name,
  studio_id,
  engine_security_group_id,
  workspace_security_group_id,
  vpc_id,
  subnet_ids
from
  aws_emr_studio
where
  engine_security_group_id is null
  or workspace_security_group_id is null
  or vpc_id is null
  or subnet_ids is null;
```

```sql+sqlite
select
  name,
  studio_id,
  engine_security_group_id,
  workspace_security_group_id,
  vpc_id,
  subnet_ids
from
  aws_emr_studio
where
  engine_security_group_id is null
  or workspace_security_group_id is null
  or vpc_id is null
  or json_array_length(subnet_ids) = 0;
```

### List studios by creation date
Find studios created within a specific time period to track recent deployments or for audit purposes.

```sql+postgres
select
  name,
  studio_id,
  creation_time,
  auth_mode,
  vpc_id
from
  aws_emr_studio
where
  creation_time >= now() - interval '30 days'
order by
  creation_time desc;
```

```sql+sqlite
select
  name,
  studio_id,
  creation_time,
  auth_mode,
  vpc_id
from
  aws_emr_studio
where
  creation_time >= datetime('now', '-30 days')
order by
  creation_time desc;
```

### Find studios with specific subnet configurations
Identify studios that are using specific subnets, which can be useful for network planning and security analysis.

```sql+postgres
select
  name,
  studio_id,
  vpc_id,
  subnet_ids
from
  aws_emr_studio
where
  subnet_ids ? 'subnet-1234567890abcdef0';
```

```sql+sqlite
select
  name,
  studio_id,
  vpc_id,
  subnet_ids
from
  aws_emr_studio
where
  json_extract(subnet_ids, '$[*]') like '%subnet-1234567890abcdef0%';
```

### List studios with service role configurations
Find studios that have specific service role configurations, which is important for understanding IAM permissions.

```sql+postgres
select
  name,
  studio_id,
  service_role,
  user_role,
  auth_mode
from
  aws_emr_studio
where
  service_role like '%EMRStudioServiceRole%';
```

```sql+sqlite
select
  name,
  studio_id,
  service_role,
  user_role,
  auth_mode
from
  aws_emr_studio
where
  service_role like '%EMRStudioServiceRole%';
```
