---
title: "Table: aws_redshiftserverless_workgroup - Query AWS Redshift Serverless Workgroup using SQL"
description: "Allows users to query AWS Redshift Serverless Workgroup information, including workgroup details, query execution settings, and enforce workgroup configuration."
---

# Table: aws_redshiftserverless_workgroup - Query AWS Redshift Serverless Workgroup using SQL

The `aws_redshiftserverless_workgroup` table in Steampipe provides information about workgroups within AWS Redshift Serverless. This table allows DevOps engineers to query workgroup-specific details, including query execution settings, enforce workgroup configuration, and associated metadata. Users can utilize this table to gather insights on workgroups, such as workgroup settings, enforced configurations, query execution details, and more. The schema outlines the various attributes of the Redshift Serverless workgroup, including the workgroup name, state, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshiftserverless_workgroup` table, you can use the `.inspect aws_redshiftserverless_workgroup` command in Steampipe.

### Key columns:

- `name`: The name of the workgroup. This can be used to join this table with others that contain workgroup names.
- `state`: The state of the workgroup (ENABLED or DISABLED). This column can be useful in understanding the status of the workgroup.
- `tags`: The metadata tags assigned to the workgroup. These can be used to join this table with others that contain workgroup tags.

## Examples

### Basic info

```sql
select
  workgroup_name,
  workgroup_arn,
  workgroup_id,
  base_capacity,
  creation_date,
  region,
  status
from
  aws_redshiftserverless_workgroup;
```

### List unavailable workgroups

```sql
select
  workgroup_name,
  workgroup_arn,
  workgroup_id,
  base_capacity,
  creation_date,
  region,
  status
from
  aws_redshiftserverless_workgroup
where
  status <> 'AVAILABLE';
```

### List publicly accessible workgroups

```sql
select
  workgroup_name,
  workgroup_arn,
  workgroup_id,
  base_capacity,
  creation_date,
  region,
  status
from
  aws_redshiftserverless_workgroup
where
  publicly_accessible;
```

### Get total base capacity utilized by available workgroups

```sql
select
  sum(base_capacity) total_base_capacity
from
  aws_redshiftserverless_workgroup
where
  status = 'AVAILABLE';
```

### Get endpoint details of each workgroups

```sql
select
  workgroup_arn,
  endpoint ->> 'Address' as endpoint_address,
  endpoint ->> 'Port' as endpoint_port,
  endpoint -> 'VpcEndpoints' as endpoint_vpc_details
from
  aws_redshiftserverless_workgroup;
```

### List config parameters associated with each workgroup

```sql
select
  workgroup_arn,
  p ->> 'ParameterKey' as parameter_key,
  p ->> 'ParameterValue' as parameter_value
from
  aws_redshiftserverless_workgroup,
  jsonb_array_elements(config_parameters) p;
```
