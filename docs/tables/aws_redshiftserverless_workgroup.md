# Table: aws_redshiftserverless_workgroup

Amazon Redshift Serverless workgroups are collections of compute resources. The compute-related workgroup groups together compute resources like RPUs, VPC subnet groups, and security groups. Properties for the workgroup include network and security settings. Other resources that are grouped under workgroups include access and usage limits.

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
