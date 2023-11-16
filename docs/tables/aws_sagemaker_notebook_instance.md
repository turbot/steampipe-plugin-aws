---
title: "Table: aws_sagemaker_notebook_instance - Query AWS SageMaker Notebook Instances using SQL"
description: "Allows users to query AWS SageMaker Notebook Instances to gather information about their configuration, status, and other related details."
---

# Table: aws_sagemaker_notebook_instance - Query AWS SageMaker Notebook Instances using SQL

The `aws_sagemaker_notebook_instance` table in Steampipe provides information about Notebook Instances within AWS SageMaker. This table allows DevOps engineers, data scientists, and other AWS users to query Notebook Instance-specific details, including instance status, instance type, associated roles, and other metadata. Users can utilize this table to gather insights on instances, such as instances with certain roles, instance statuses, and more. The schema outlines the various attributes of the SageMaker Notebook Instance, including the instance name, instance type, role ARN, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sagemaker_notebook_instance` table, you can use the `.inspect aws_sagemaker_notebook_instance` command in Steampipe.

**Key columns**:

- `name`: The name of the notebook instance. It can be used to join this table with other tables that contain notebook instance names.
- `instance_type`: The type of the notebook instance. This column can be used to filter instances based on their types.
- `role_arn`: The Amazon Resource Name (ARN) of the IAM role associated with the notebook instance. This column can be used to join this table with other IAM role-related tables.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  instance_type,
  notebook_instance_status
from
  aws_sagemaker_notebook_instance;
```


### List notebook instances that do not have encryption at rest enabled

```sql
select
  name,
  kms_key_id
from
  aws_sagemaker_notebook_instance
where
  kms_key_id is null;
```


### List publicly available notebook instances

```sql
select
  name,
  direct_internet_access
from
  aws_sagemaker_notebook_instance
where
  direct_internet_access = 'Disabled';
```


### List notebook instances that allow root access

```sql
select
  name,
  root_access
from
  aws_sagemaker_notebook_instance
where
  root_access = 'Enabled';
```
