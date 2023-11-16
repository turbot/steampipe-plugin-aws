---
title: "Table: aws_rds_db_instance - Query AWS RDS DB Instances using SQL"
description: "Allows users to query AWS RDS DB Instances for detailed information about the configuration, status, and other metadata associated with each database instance."
---

# Table: aws_rds_db_instance - Query AWS RDS DB Instances using SQL

The `aws_rds_db_instance` table in Steampipe provides comprehensive information about each database instance within Amazon Relational Database Service (RDS). This table enables DevOps engineers, database administrators, and other technical professionals to query detailed information about each DB instance, including its configuration, status, performance metrics, and other associated metadata. Users can leverage this table to gather insights about DB instances, such as instance specifications, security configurations, backup policies, and more. The schema outlines the various attributes of the DB instance, including the DB instance identifier, instance class, engine version, storage type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance` table, you can use the `.inspect aws_rds_db_instance` command in Steampipe.

### Key columns:

- `db_instance_identifier`: The identifier for the DB instance. This identifier is unique across all DB instances within an AWS account and can be used to join this table with other tables that contain DB instance identifiers.
- `arn`: The Amazon Resource Name (ARN) of the DB instance. The ARN is a unique identifier for the DB instance and can be used to join this table with other tables that contain DB instance ARNs.
- `vpc_id`: The identifier of the VPC that the DB instance is in. This identifier can be used to join this table with other tables that contain VPC identifiers, providing insights into the network configuration of the DB instance.

## Examples

### Basic info

```sql
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

```sql
select
  db_instance_identifier,
  publicly_accessible
from
  aws_rds_db_instance
where
  publicly_accessible;
```

### List DB instances which are not authenticated through IAM users and roles

```sql
select
  db_instance_identifier,
  iam_database_authentication_enabled
from
  aws_rds_db_instance
where
  not iam_database_authentication_enabled;
```

### Get VPC and subnet info for each DB instance

```sql
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

### List DB instances with deletion protection disabled

```sql
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

### List DB instances with unecrypted storage

```sql
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

### Get endpoint info for each DB instance

```sql
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_rds_db_instance;
```

### List SQL Server DB instances with SSL disabled in assigned parameter group

```sql
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

### List DB instance pending maintenance actions

```sql
select
  actions ->> 'ResourceIdentifier' as db_instance_identifier,
  details ->> 'Action' as action,
  details ->> 'OptInStatus' as opt_in_status,
  details ->> 'ForcedApplyDate' as forced_apply_date,
  details ->> 'CurrentApplyDate' as current_apply_date,
  details ->> 'AutoAppliedAfterDate' as auto_applied_after_date
from
  aws_rds_db_instance,
  jsonb_array_elements(pending_maintenance_actions) as actions,
  jsonb_array_elements(actions -> 'PendingMaintenanceActionDetails') as details;
```

### List certificate details associated to the instance

```sql
select
  arn,
  certificate ->> 'CertificateArn' as certificate_arn,
  certificate ->> 'CertificateType' as certificate_type,
  certificate ->> 'ValidFrom' as valid_from,
  certificate ->> 'ValidTill' as valid_till
from
  aws_rds_db_instance;
```

### List certificates valid for less than 90 days

```sql
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