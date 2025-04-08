---
title: "Steampipe Table: aws_sagemaker_domain - Query AWS SageMaker Domains using SQL"
description: "Allows users to query AWS SageMaker Domains to retrieve data about AWS SageMaker Studio domains, including domain details, status, and associated metadata."
folder: "SageMaker"
---

# Table: aws_sagemaker_domain - Query AWS SageMaker Domains using SQL

The AWS SageMaker Domain is a fully managed service that provides every developer and data scientist with the ability to build, train, and deploy machine learning (ML) models quickly. SageMaker removes the heavy lifting from each step of the machine learning process to make it easier to develop high-quality models. It offers a set of tools for developers and data scientists to iteratively develop, tune, and deploy machine learning models.

## Table Usage Guide

The `aws_sagemaker_domain` table in Steampipe provides you with information about domains within AWS SageMaker Studio. This table allows you, as a data scientist, machine learning engineer, or DevOps engineer, to query domain-specific details, including the domain status, creation time, and associated metadata. You can utilize this table to gather insights on domains, such as the status of a domain, the creation time, the associated app network access type, and more. The schema outlines the various attributes of the SageMaker domain, including the domain ID, domain ARN, domain name, and associated tags for you.

## Examples

### Basic info
Explore which AWS Sagemaker domains are active or inactive and their respective creation times. This can be useful in managing and monitoring the lifecycle of your machine learning environments.

```sql+postgres
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_domain;
```

```sql+sqlite
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_domain;
```

### List sagemaker domains where EFS volume is unencrypted
Discover the segments that have unencrypted EFS volumes in SageMaker domains. This is useful to identify potential security risks and take necessary corrective actions.

```sql+postgres
select
  name,
  creation_time,
  home_efs_file_system_id,
  kms_key_id
from
  aws_sagemaker_domain
where 
  kms_key_id is null;
```

```sql+sqlite
select
  name,
  creation_time,
  home_efs_file_system_id,
  kms_key_id
from
  aws_sagemaker_domain
where 
  kms_key_id is null;
```

### List publicly accessible sagemaker domains
Discover the segments that have 'PublicInternetOnly' as their network access type to identify publicly accessible domains. This is particularly useful in assessing the security and accessibility of your network resources.

```sql+postgres
select
  name,
  arn,
  creation_time,
  app_network_access_type
from
  aws_sagemaker_domain
where 
  app_network_access_type = 'PublicInternetOnly';
```

```sql+sqlite
select
  name,
  arn,
  creation_time,
  app_network_access_type
from
  aws_sagemaker_domain
where 
  app_network_access_type = 'PublicInternetOnly';
```