---
title: "Steampipe Table: aws_redshiftserverless_workgroup - Query AWS Redshift Serverless Workgroup using SQL"
description: "Allows users to query AWS Redshift Serverless Workgroup information, including workgroup details, query execution settings, and enforce workgroup configuration."
folder: "Redshift"
---

# Table: aws_redshiftserverless_workgroup - Query AWS Redshift Serverless Workgroup using SQL

The AWS Redshift Serverless Workgroup is a feature of Amazon Redshift that enables you to separate query processing among different sets of users. It allows you to manage query concurrency, memory allocation, and user access for better performance and security. Through SQL, you can query and analyze the workgroup's configurations and usage statistics.

## Table Usage Guide

The `aws_redshiftserverless_workgroup` table in Steampipe provides you with information about workgroups within AWS Redshift Serverless. This table allows you as a DevOps engineer to query workgroup-specific details, including query execution settings, enforce workgroup configuration, and associated metadata. You can utilize this table to gather insights on workgroups, such as workgroup settings, enforced configurations, query execution details, and more. The schema outlines the various attributes of the Redshift Serverless workgroup for you, including the workgroup name, state, creation time, and associated tags.

## Examples

### Basic info
Explore which AWS Redshift Serverless Workgroups are available and assess their status. This can help in efficiently managing resources and understanding server capacity.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which workgroups are not currently available. This is useful for identifying potential issues or bottlenecks within your AWS Redshift serverless environment.

```sql+postgres
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

```sql+sqlite
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
  status != 'AVAILABLE';
```

### List publicly accessible workgroups
Discover the segments that are publicly accessible within your workgroups. This is useful for understanding potential security risks and which areas of your system may be exposed to the public.

```sql+postgres
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

```sql+sqlite
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
  publicly_accessible = 1;
```

### Get total base capacity utilized by available workgroups
Determine the total capacity utilized by all available workgroups in your AWS Redshift Serverless environment. This is useful for understanding the extent of your resource usage and planning for future capacity needs.

```sql+postgres
select
  sum(base_capacity) total_base_capacity
from
  aws_redshiftserverless_workgroup
where
  status = 'AVAILABLE';
```

```sql+sqlite
select
  sum(base_capacity) total_base_capacity
from
  aws_redshiftserverless_workgroup
where
  status = 'AVAILABLE';
```

### Get endpoint details of each workgroups
Discover the segments that provide details about the endpoints of each workgroup in your AWS Redshift serverless environment. This can be beneficial when you need to understand the connectivity details of your serverless workgroups for auditing or troubleshooting purposes.

```sql+postgres
select
  workgroup_arn,
  endpoint ->> 'Address' as endpoint_address,
  endpoint ->> 'Port' as endpoint_port,
  endpoint -> 'VpcEndpoints' as endpoint_vpc_details
from
  aws_redshiftserverless_workgroup;
```

```sql+sqlite
select
  workgroup_arn,
  json_extract(endpoint, '$.Address') as endpoint_address,
  json_extract(endpoint, '$.Port') as endpoint_port,
  json_extract(endpoint, '$.VpcEndpoints') as endpoint_vpc_details
from
  aws_redshiftserverless_workgroup;
```

### List config parameters associated with each workgroup
Discover the segments that contain specific configurations in each workgroup to understand how they are set up and operate. This can be particularly useful for auditing, debugging, or optimization purposes.

```sql+postgres
select
  workgroup_arn,
  p ->> 'ParameterKey' as parameter_key,
  p ->> 'ParameterValue' as parameter_value
from
  aws_redshiftserverless_workgroup,
  jsonb_array_elements(config_parameters) p;
```

```sql+sqlite
select
  workgroup_arn,
  json_extract(p.value, '$.ParameterKey') as parameter_key,
  json_extract(p.value, '$.ParameterValue') as parameter_value
from
  aws_redshiftserverless_workgroup,
  json_each(config_parameters) as p;
```