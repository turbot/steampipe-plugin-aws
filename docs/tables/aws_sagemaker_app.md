---
title: "Steampipe Table: aws_sagemaker_app - Query AWS SageMaker App using SQL"
description: "Allows users to query AWS SageMaker App data, providing detailed insights into application configurations, user settings, and associated metadata."
folder: "SageMaker"
---

# Table: aws_sagemaker_app - Query AWS SageMaker App using SQL

The AWS SageMaker App is a component of Amazon SageMaker that provides a platform for developers and data scientists to build, train, and deploy machine learning models quickly. It offers a fully managed service that covers the entire machine learning workflow to label and prepare your data, choose an algorithm, train the model, tune and optimize it for deployment, make predictions, and take action. SageMaker App simplifies the process of building, training, and deploying machine learning models, allowing you to get your models to production faster with much less effort and at lower cost.

## Table Usage Guide

The `aws_sagemaker_app` table in Steampipe provides you with information about AWS SageMaker Apps. This table enables you, as a DevOps engineer, data scientist, or other technical professional, to query application-specific details, including resource names, ARNs, types, statuses, and creation times. You can utilize this table to gather insights on SageMaker Apps, such as their configurations, user settings, and more. The schema outlines the various attributes of the SageMaker App for you, including the app ARN, creation time, status, and associated tags.

## Examples

### Basic info
Explore the status and creation time of your AWS Sagemaker apps to understand their current operational state and longevity. This can aid in resource management and troubleshooting.

```sql+postgres
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_app;
```

```sql+sqlite
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_app;
```

### List apps that failed to create
Discover the segments that failed during the application creation process within AWS SageMaker. This is useful to understand and rectify issues that prevent successful application creation.

```sql+postgres
select
  name,
  arn,
  creation_time,
  status,
  failure_reason
from
  aws_sagemaker_app
where 
  status = 'Failed';
```

```sql+sqlite
select
  name,
  arn,
  creation_time,
  status,
  failure_reason
from
  aws_sagemaker_app
where 
  status = 'Failed';
```