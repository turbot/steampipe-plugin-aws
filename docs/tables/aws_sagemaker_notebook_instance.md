---
title: "Steampipe Table: aws_sagemaker_notebook_instance - Query AWS SageMaker Notebook Instances using SQL"
description: "Allows users to query AWS SageMaker Notebook Instances to gather information about their configuration, status, and other related details."
folder: "SageMaker"
---

# Table: aws_sagemaker_notebook_instance - Query AWS SageMaker Notebook Instances using SQL

The AWS SageMaker Notebook Instances are a fully managed service that provides Jupyter notebooks for data exploration, cleaning, and preprocessing. They also provide a development environment to create machine learning models and experiments. These instances allow you to seamlessly connect to your data stored in AWS S3, AWS DynamoDB, AWS Redshift, and more, facilitating easier data manipulation and analysis.

## Table Usage Guide

The `aws_sagemaker_notebook_instance` table in Steampipe provides you with information about Notebook Instances within AWS SageMaker. This table allows you, as a DevOps engineer, data scientist, or other AWS user, to query Notebook Instance-specific details, including instance status, instance type, associated roles, and other metadata. You can utilize this table to gather insights on instances, such as instances with certain roles, instance statuses, and more. The schema outlines the various attributes of the SageMaker Notebook Instance for you, including the instance name, instance type, role ARN, creation time, and associated tags.

## Examples

### Basic info
Determine the areas in which AWS SageMaker notebook instances are being used, by examining their creation times, instance types, and current statuses. This allows for better resource management and operational oversight.

```sql+postgres
select
  name,
  arn,
  creation_time,
  instance_type,
  notebook_instance_status
from
  aws_sagemaker_notebook_instance;
```

```sql+sqlite
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
Identify instances where AWS SageMaker notebook instances lack encryption at rest, a crucial security feature. This can help in enhancing data security by pinpointing areas that need attention.

```sql+postgres
select
  name,
  kms_key_id
from
  aws_sagemaker_notebook_instance
where
  kms_key_id is null;
```

```sql+sqlite
select
  name,
  kms_key_id
from
  aws_sagemaker_notebook_instance
where
  kms_key_id is null;
```


### List publicly available notebook instances
Uncover the details of SageMaker notebook instances that have disabled direct internet access, allowing you to assess security measures and ensure data protection.

```sql+postgres
select
  name,
  direct_internet_access
from
  aws_sagemaker_notebook_instance
where
  direct_internet_access = 'Disabled';
```

```sql+sqlite
select
  name,
  direct_internet_access
from
  aws_sagemaker_notebook_instance
where
  direct_internet_access = 'Disabled';
```


### List notebook instances that allow root access
Identify instances where root access is enabled in your AWS Sagemaker notebook instances, which could potentially pose security risks. This is useful for maintaining and improving security measures within your system.

```sql+postgres
select
  name,
  root_access
from
  aws_sagemaker_notebook_instance
where
  root_access = 'Enabled';
```

```sql+sqlite
select
  name,
  root_access
from
  aws_sagemaker_notebook_instance
where
  root_access = 'Enabled';
```