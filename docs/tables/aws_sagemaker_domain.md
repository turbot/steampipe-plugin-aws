---
title: "Table: aws_sagemaker_domain - Query AWS SageMaker Domains using SQL"
description: "Allows users to query AWS SageMaker Domains to retrieve data about AWS SageMaker Studio domains, including domain details, status, and associated metadata."
---

# Table: aws_sagemaker_domain - Query AWS SageMaker Domains using SQL

The `aws_sagemaker_domain` table in Steampipe provides information about domains within AWS SageMaker Studio. This table allows data scientists, machine learning engineers, and DevOps engineers to query domain-specific details, including the domain status, creation time, and associated metadata. Users can utilize this table to gather insights on domains, such as the status of a domain, the creation time, the associated app network access type, and more. The schema outlines the various attributes of the SageMaker domain, including the domain ID, domain ARN, domain name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sagemaker_domain` table, you can use the `.inspect aws_sagemaker_domain` command in Steampipe.

**Key columns**:

- `domain_id`: The ID of the domain. This is a unique identifier and can be used to join this table with other tables that contain domain-specific information.
- `domain_arn`: The Amazon Resource Name (ARN) of the domain. This can be used to join this table with other tables that require the domain's ARN.
- `status`: The status of the domain. This column can be used to track the status of a domain and can be useful for monitoring and alerting purposes.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_domain;
```

### List sagemaker domains where EFS volume is unencrypted

```sql
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

```sql
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