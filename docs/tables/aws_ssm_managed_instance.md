---
title: "Table: aws_ssm_managed_instance - Query AWS SSM Managed Instances using SQL"
description: "Allows users to query AWS SSM Managed Instances to retrieve their configuration and status information."
---

# Table: aws_ssm_managed_instance - Query AWS SSM Managed Instances using SQL

The `aws_ssm_managed_instance` table in Steampipe provides information about managed instances within AWS Systems Manager (SSM). This table allows DevOps engineers to query managed instance-specific details, including instance ID, name, platform type, platform version, and associated metadata. Users can utilize this table to gather insights on instances, such as their operational status, last ping time, agent version, and more. The schema outlines the various attributes of the managed instance, including the instance ARN, registration date, resource type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_managed_instance` table, you can use the `.inspect aws_ssm_managed_instance` command in Steampipe.

**Key columns**:

- `instance_id`: The ID of the managed instance. This is a key column for joining with other tables because it uniquely identifies each managed instance.
- `name`: The name of the managed instance. This column is useful for joining with other tables when the instance ID is not known.
- `platform_type`: The type of operating system running on the managed instance. This column can be used to filter instances based on their operating system type.

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
