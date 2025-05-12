---
title: "Steampipe Table: aws_rds_db_instance - Query AWS RDS DB Instances using SQL"
description: "Allows users to query AWS RDS DB Instances for detailed information about the configuration, status, and other metadata associated with each database instance."
folder: "RDS"
---

# Table: aws_rds_db_instance - Query AWS RDS DB Instances using SQL

Title: AWS RDS DB Instances

Description: AWS RDS DB Instances are a part of Amazon Relational Database Service (RDS), which makes it easier to set up, operate, and scale a relational database in the cloud. They provide cost-efficient and resizable capacity while automating time-consuming administration tasks such as hardware provisioning, database setup, patching, and backups. It frees you to focus on your applications so you can give them the fast performance, high availability, security and compatibility they need.

## Table Usage Guide

The `aws_rds_db_instance` table in Steampipe provides you with comprehensive information about each database instance within Amazon Relational Database Service (RDS). This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query detailed information about each DB instance, including its configuration, status, performance metrics, and other associated metadata. You can leverage this table to gather insights about DB instances, such as instance specifications, security configurations, backup policies, and more. The schema outlines the various attributes of the DB instance for you, including the DB instance identifier, instance class, engine version, storage type, and associated tags.

## Examples

### Basic info
Explore which Amazon RDS database instances are publicly accessible, along with their identifiers, classes, engines, and versions. This information can help determine areas where security might be improved by limiting public access.

```sql+postgres
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  publicly_accessible
from
  aws_rds_db_instance
```

```sql+sqlite
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  publicly_accessible
from
  aws_rds_db_instance
```

### List DB instances which are publicly accessible
Determine the areas in which your database instances are publicly accessible. This is useful for identifying potential security risks and ensuring that your data is protected.

```sql+postgres
select
  db_instance_identifier,
  publicly_accessible
from
  aws_rds_db_instance
where
  publicly_accessible;
```

```sql+sqlite
select
  db_instance_identifier,
  publicly_accessible
from
  aws_rds_db_instance
where
  publicly_accessible = 1;
```

### List DB instances which are not authenticated through IAM users and roles
Identify instances where database instances are not utilizing IAM users and roles for authentication. This query is useful in highlighting potential security risks and enforcing best practices for access management.

```sql+postgres
select
  db_instance_identifier,
  iam_database_authentication_enabled
from
  aws_rds_db_instance
where
  not iam_database_authentication_enabled;
```

```sql+sqlite
select
  db_instance_identifier,
  iam_database_authentication_enabled
from
  aws_rds_db_instance
where
  iam_database_authentication_enabled = 0;
```

### Get VPC and subnet info for each DB instance
Determine the areas in which your database instances are located and the security measures in place. This helps ensure your databases are in the desired regions and properly secured, enhancing your data management strategy.

```sql+postgres
select
  db_instance_identifier as attached_vpc,
  vsg ->> 'VpcSecurityGroupId' as vpc_security_group_id,
  vsg ->> 'Status' as status,
  sub -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  sub ->> 'SubnetIdentifier' as subnet_identifier,
  sub -> 'SubnetOutpost' ->> 'Arn' as subnet_outpost,
  sub ->> 'SubnetStatus' as subnet_status
from
  aws_rds_db_instance
  cross join jsonb_array_elements(vpc_security_groups) as vsg
  cross join jsonb_array_elements(subnets) as sub;
```

```sql+sqlite
select
  db_instance_identifier as attached_vpc,
  json_extract(vsg.value, '$.VpcSecurityGroupId') as vpc_security_group_id,
  json_extract(vsg.value, '$.Status') as status,
  json_extract(json_extract(sub.value, '$.SubnetAvailabilityZone'), '$.Name') as subnet_availability_zone,
  json_extract(sub.value, '$.SubnetIdentifier') as subnet_identifier,
  json_extract(json_extract(sub.value, '$.SubnetOutpost'), '$.Arn') as subnet_outpost,
  json_extract(sub.value, '$.SubnetStatus') as subnet_status
from
  aws_rds_db_instance,
  json_each(vpc_security_groups) as vsg,
  json_each(subnets) as sub;
```

### List DB instances with deletion protection disabled
Discover the segments that consist of DB instances where deletion protection is not enabled. This is useful for identifying potential security vulnerabilities within your AWS RDS DB instances.

```sql+postgres
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  deletion_protection
from
  aws_rds_db_instance
where
  not deletion_protection;
```

```sql+sqlite
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  deletion_protection
from
  aws_rds_db_instance
where
  deletion_protection = 0;
```

### List DB instances with unecrypted storage
Discover the segments that are using unencrypted storage within your database instances to assess potential vulnerabilities and improve your system's security. This query is particularly useful in identifying areas where sensitive data may be at risk due to lack of encryption.

```sql+postgres
select
  db_instance_identifier,
  class,
  allocated_storage,
  deletion_protection
from
  aws_rds_db_instance
where
  not storage_encrypted;
```

```sql+sqlite
select
  db_instance_identifier,
  class,
  allocated_storage,
  deletion_protection
from
  aws_rds_db_instance
where
  storage_encrypted = 0;
```

### Get endpoint info for each DB instance
Explore which database instances are linked to specific endpoints. This information can help assess the network connections and their respective ports, aiding in efficient database management and troubleshooting.

```sql+postgres
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_rds_db_instance;
```

```sql+sqlite
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_rds_db_instance;
```

### List SQL Server DB instances with SSL disabled in assigned parameter group
Identify instances where SQL Server database instances have SSL disabled in their assigned parameter groups. This is useful for uncovering potential security vulnerabilities related to unencrypted data transmission.

```sql+postgres
with db_parameter_group as (
  select
    name as db_parameter_group_name,
    pg ->> 'ParameterName' as parameter_name,
    pg ->> 'ParameterValue' as parameter_value
  from
    aws_rds_db_parameter_group,
    jsonb_array_elements(parameters) as pg
  where
    -- The example is limited to SQL Server, this may change based on DB engine
    pg ->> 'ParameterName' like 'rds.force_ssl'
    and name not like 'default.%'
),
 rds_associated_parameter_group as (
  select
    db_instance_identifier as db_instance_identifier,
    arn,
    pg ->> 'DBParameterGroupName' as DBParameterGroupName
  from
    aws_rds_db_instance,
    jsonb_array_elements(db_parameter_groups) as pg
  where
    engine like 'sqlserve%'
)
select
  rds.db_instance_identifier as name,
  rds.DBParameterGroupName,
  parameter_name,
  parameter_value
from
  rds_associated_parameter_group as rds
  left join db_parameter_group d on rds.DBParameterGroupName = d.db_parameter_group_name
where
  parameter_value = '0'
```

```sql+sqlite
with db_parameter_group as (
  select
    name as db_parameter_group_name,
    json_extract(pg.value, '$.ParameterName') as parameter_name,
    json_extract(pg.value, '$.ParameterValue') as parameter_value
  from
    aws_rds_db_parameter_group,
    json_each(parameters) as pg
  where
    json_extract(pg.value, '$.ParameterName') like 'rds.force_ssl'
    and name not like 'default.%'
),
 rds_associated_parameter_group as (
  select
    db_instance_identifier as db_instance_identifier,
    arn,
    json_extract(pg.value, '$.DBParameterGroupName') as DBParameterGroupName
  from
    aws_rds_db_instance,
    json_each(db_parameter_groups) as pg
  where
    engine like 'sqlserve%'
)
select
  rds.db_instance_identifier as name,
  rds.DBParameterGroupName,
  parameter_name,
  parameter_value
from
  rds_associated_parameter_group as rds
  left join db_parameter_group d on rds.DBParameterGroupName = d.db_parameter_group_name
where
  parameter_value = '0';
```

### List certificate details associated to the instance
Discover the segments that highlight the validity and type of certificates associated with a specific instance. This can be especially useful in managing and tracking certificate expiration dates, ensuring the security and reliability of your database instances.

```sql+postgres
select
  arn,
  certificate ->> 'CertificateArn' as certificate_arn,
  certificate ->> 'CertificateType' as certificate_type,
  certificate ->> 'ValidFrom' as valid_from,
  certificate ->> 'ValidTill' as valid_till
from
  aws_rds_db_instance;
```

```sql+sqlite
select
  arn,
  json_extract(certificate, '$.CertificateArn') as certificate_arn,
  json_extract(certificate, '$.CertificateType') as certificate_type,
  json_extract(certificate, '$.ValidFrom') as valid_from,
  json_extract(certificate, '$.ValidTill') as valid_till
from
  aws_rds_db_instance;
```

### List certificates valid for less than 90 days
Discover the segments that have certificates valid for less than 90 days in your AWS RDS database instances. This can help in proactive renewal of certificates, thereby avoiding any service disruption due to certificate expiry.

```sql+postgres
select
  arn,
  certificate ->> 'CertificateArn' as certificate_arn,
  certificate ->> 'CertificateType' as certificate_type,
  certificate ->> 'ValidFrom' as valid_from,
  certificate ->> 'ValidTill' as valid_till
from
  aws_rds_db_instance
where
  (certificate ->> 'ValidTill')::timestamp <= (current_date - interval '90' day);
```

```sql+sqlite
select
  arn,
  json_extract(certificate, '$.CertificateArn') as certificate_arn,
  json_extract(certificate, '$.CertificateType') as certificate_type,
  json_extract(certificate, '$.ValidFrom') as valid_from,
  json_extract(certificate, '$.ValidTill') as valid_till
from
  aws_rds_db_instance
where
  julianday('now') - julianday(json_extract(certificate, '$.ValidTill')) >= 90;
```

### Listing RDS DB Instances with Existing Processor Features
Supports Infrastructure as Code (IaC) and Automation For organizations using IaC practices or automation in their cloud environments, such queries can help in generating reports, monitoring configurations, or triggering workflows based on the state of RDS instances.

```sql+postgres
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  kms_key_id,
  processor_features
from
  aws_rds_db_instance
where
  processor_features is not null;
```

```sql+sqlite
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  kms_key_id,
  processor_features
from
  aws_rds_db_instance
where
  processor_features is not null;
```

### Get RDS DB instances pending maintenance actions
Get DB instances pending maintenance actions to plan and prioritize maintenance schedules effectively.

```sql+postgres
select
  a.db_instance_identifier,
  b.action,
  a.status,
  b.opt_in_status,
  b.forced_apply_date,
  b.current_apply_date,
  b.auto_applied_after_date
from 
  aws_rds_db_instance as a
  join aws_rds_pending_maintenance_action as b on b.resource_identifier = a.arn;
```

```sql+sqlite
select
  a.db_instance_identifier,
  b.action,
  a.status,
  b.opt_in_status,
  b.forced_apply_date,
  b.current_apply_date,
  b.auto_applied_after_date
from 
  aws_rds_db_instance as a
  join aws_rds_pending_maintenance_action as b on b.resource_identifier = a.arn;
```