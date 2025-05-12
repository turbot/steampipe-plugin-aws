---
title: "Steampipe Table: aws_cloudformation_stack_resource - Query AWS CloudFormation Stack Resources using SQL"
description: "Allows users to query AWS CloudFormation Stack Resources, providing details about each resource within the stack, including its status, type, and associated metadata. This table is useful for managing and analyzing AWS CloudFormation resources."
folder: "CloudFormation"
---

# Table: aws_cloudformation_stack_resource - Query AWS CloudFormation Stack Resources using SQL

The AWS CloudFormation Stack Resources are the AWS resources that are part of a stack. AWS CloudFormation simplifies the process of managing your AWS resources by treating all the resources as a single unit, called a stack. These resources can be created, updated, or deleted in a single operation, making it easier to manage and configure all the resources collectively.

## Table Usage Guide

The `aws_cloudformation_stack_resource` table in Steampipe provides you with information about Stack Resources within AWS CloudFormation. This table allows you, as a DevOps engineer, to query resource-specific details, including the current status, resource type, and associated metadata. You can utilize this table to gather insights on resources, such as resource status, the type of resources used in the stack, and more. The schema outlines the various attributes of the Stack Resource for you, including the stack name, resource status, logical resource id, and physical resource id.

## Examples

### Basic info
Explore the status and type of resources within your AWS CloudFormation stack to better understand your stack's configuration and resource allocation. This allows for effective resource management and helps identify potential issues in your stack's setup.

```sql+postgres
select
  stack_name,
  stack_id,
  logical_resource_id,
  resource_type,
  resource_status
from
  aws_cloudformation_stack_resource;
```

```sql+sqlite
select
  stack_name,
  stack_id,
  logical_resource_id,
  resource_type,
  resource_status
from
  aws_cloudformation_stack_resource;
```

### List cloudformation stack resources having rollback disabled
Determine the areas in your AWS CloudFormation setup where rollback is disabled, allowing you to understand potential risk points in your infrastructure. This can be useful in identifying instances where a failure in stack creation or update could lead to resource inconsistencies.

```sql+postgres
select
  s.name,
  s.disable_rollback,
  r.logical_resource_id,
  r.resource_status
from
  aws_cloudformation_stack_resource as r,
  aws_cloudformation_stack as s
where
  r.stack_id = s.id
  and s.disable_rollback;
```

```sql+sqlite
select
  s.name,
  s.disable_rollback,
  r.logical_resource_id,
  r.resource_status
from
  aws_cloudformation_stack_resource as r
join
  aws_cloudformation_stack as s
on
  r.stack_id = s.id
where
  s.disable_rollback = 1;
```

### List resources having termination protection disabled
Determine the areas in which resources could be at risk due to disabled termination protection. This is useful for identifying potential vulnerabilities within your CloudFormation stacks.

```sql+postgres
select
  s.name,
  s.enable_termination_protection,
  s.disable_rollback,
  r.logical_resource_id,
  r.resource_status
from
  aws_cloudformation_stack_resource as r,
  aws_cloudformation_stack as s
where
  r.stack_id = s.id
  and not enable_termination_protection;
```

```sql+sqlite
select
  s.name,
  s.enable_termination_protection,
  s.disable_rollback,
  r.logical_resource_id,
  r.resource_status
from
  aws_cloudformation_stack_resource as r
join
  aws_cloudformation_stack as s
on
  r.stack_id = s.id
where
  not s.enable_termination_protection;
```

### List stack resources of type VPC
Discover the segments that are utilizing Virtual Private Cloud (VPC) resources within your AWS CloudFormation stacks. This is useful for understanding your resource allocation and identifying any potential areas of optimization.

```sql+postgres
select
  stack_name,
  stack_id,
  logical_resource_id,
  resource_status,
  resource_type
from
  aws_cloudformation_stack_resource
where
  resource_type = 'AWS::EC2::VPC';
```

```sql+sqlite
select
  stack_name,
  stack_id,
  logical_resource_id,
  resource_status,
  resource_type
from
  aws_cloudformation_stack_resource
where
  resource_type = 'AWS::EC2::VPC';
```

### List resources that failed to update
Identify instances where updates to cloud resources failed. This can help in troubleshooting and rectifying issues to ensure smooth operation of your cloud infrastructure.

```sql+postgres
select
  stack_name,
  logical_resource_id,
  resource_status,
  resource_type
from
  aws_cloudformation_stack_resource
where
  resource_status = 'UPDATE_FAILED';
```

```sql+sqlite
select
  stack_name,
  logical_resource_id,
  resource_status,
  resource_type
from
  aws_cloudformation_stack_resource
where
  resource_status = 'UPDATE_FAILED';
```