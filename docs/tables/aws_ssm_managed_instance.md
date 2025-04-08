---
title: "Steampipe Table: aws_ssm_managed_instance - Query AWS SSM Managed Instances using SQL"
description: "Allows users to query AWS SSM Managed Instances to retrieve their configuration and status information."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_managed_instance - Query AWS SSM Managed Instances using SQL

The AWS Systems Manager Managed Instance is a compute resource in your environment that is configured for use with AWS Systems Manager. These can be Amazon EC2 instances or servers and virtual machines (VMs) in your on-premises environment. They provide secure and scalable configuration and automation management, enabling you to automate administrative tasks, apply compliance policies, and manage resources across your environment.

## Table Usage Guide

The `aws_ssm_managed_instance` table in Steampipe provides you with information about managed instances within AWS Systems Manager (SSM). This table allows you, as a DevOps engineer, to query managed instance-specific details, including instance ID, name, platform type, platform version, and associated metadata. You can utilize this table to gather insights on instances, such as their operational status, last ping time, agent version, and more. The schema outlines the various attributes of the managed instance for you, including the instance ARN, registration date, resource type, and associated tags.

## Examples

### Basic info
Gain insights into the status and characteristics of managed instances in AWS Simple Systems Manager (SSM). This can help in monitoring and managing resources effectively, identifying any issues with association status or outdated agent versions, and understanding the distribution of resources across different platform types.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which managed instances lack associations. This could be useful in identifying potential gaps in your resource management, allowing for more efficient allocation and utilization of resources.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which EC2 instances are not managed by the Systems Manager (SSM) to identify potential gaps in your management strategy. This query is useful for ensuring all instances are appropriately managed and can highlight areas needing attention.

```sql+postgres
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

```sql+sqlite
select
  i.instance_id,
  i.arn,
  case when m.instance_id is not null then 1 else 0 end as ssm_managed
from
  aws_ec2_instance i
left join aws_ssm_managed_instance m on m.instance_id = i.instance_id
where 
  m.instance_id is null;
```