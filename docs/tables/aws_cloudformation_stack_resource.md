---
title: "Table: aws_cloudformation_stack_resource - Query AWS CloudFormation Stack Resources using SQL"
description: "Allows users to query AWS CloudFormation Stack Resources, providing details about each resource within the stack, including its status, type, and associated metadata. This table is useful for managing and analyzing AWS CloudFormation resources."
---

# Table: aws_cloudformation_stack_resource - Query AWS CloudFormation Stack Resources using SQL

The `aws_cloudformation_stack_resource` table in Steampipe provides information about Stack Resources within AWS CloudFormation. This table allows DevOps engineers to query resource-specific details, including the current status, resource type, and associated metadata. Users can utilize this table to gather insights on resources, such as resource status, the type of resources used in the stack, and more. The schema outlines the various attributes of the Stack Resource, including the stack name, resource status, logical resource id, and physical resource id.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudformation_stack_resource` table, you can use the `.inspect aws_cloudformation_stack_resource` command in Steampipe.

Key columns:

- `stack_name`: The name associated with the stack. This can be used to join with other tables that contain stack name information.
- `resource_status`: The status of the resource (e.g., CREATE_COMPLETE, DELETE_FAILED). This is useful for tracking the lifecycle of resources within a stack.
- `resource_type`: The type of resource provisioned by CloudFormation. This can be used to join with other tables that contain resource type information.

## Examples

### Basic info

```sql
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

```sql
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

### List resources having termination protection disabled

```sql
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

### List stack resources of type VPC

```sql
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

```sql
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
