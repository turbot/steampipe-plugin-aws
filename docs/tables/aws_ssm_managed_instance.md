# Table: aws_ssm_managed_instance

A managed instance is any machine configured for AWS Systems Manager. You can configure Amazon Elastic Compute Cloud (Amazon EC2) instances or on-premises machines in a hybrid environment as managed instances. Systems Manager supports various distributions of Linux, including Raspberry Pi devices, macOS, and Microsoft Windows Server.

## Examples

### Basic info

```sql
select
  instance_id,
  arn,
  resource_type,
  association_status,
  agent_version,
  platform_type
from
  aws_ssm_managed_instance;
```

### List managed instances with no associations

```sql
select
  instance_id,
  arn,
  resource_type,
  association_status
from
  aws_ssm_managed_instance
where
  association_status is null;
```


### List EC2 instances not managed by SSM

```sql
select
  i.instance_id,
  i.arn,
  m.instance_id is not null as ssm_managed
from
  aws_ec2_instance i
left join aws_ssm_managed_instance m on m.instance_id = i.instance_id
where 
  m.instance_id is null;
```
