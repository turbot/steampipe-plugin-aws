---
title: "Steampipe Table: aws_ssm_parameter - Query AWS Systems Manager Parameter Store using SQL"
description: "Allows users to query AWS Systems Manager Parameter Store to retrieve information about parameters, their types, values, and associated metadata."
folder: "Resource Access Manager"
---

# Table: aws_ssm_parameter - Query AWS Systems Manager Parameter Store using SQL

The AWS Systems Manager Parameter Store provides secure, hierarchical storage for configuration data management and secrets management. It allows you to centrally manage your configuration data, whether plain-text data such as database strings or secrets like passwords, thus improving the security of your data by using AWS Key Management Service (KMS). Parameter Store is designed to use with other AWS services to pull configuration data and keep your applications secure and scalable.

## Table Usage Guide

The `aws_ssm_parameter` table in Steampipe provides you with information about parameters within the AWS Systems Manager Parameter Store. This table allows you, as a DevOps engineer, to query parameter-specific details, such as parameter names, types, values, and associated metadata. You can utilize this table to gather insights on parameters, such as parameter descriptions, last modification dates, and the user who last modified the parameter. The schema outlines the various attributes of the parameter for you, including the parameter ARN, type, value, and associated tags.

## Examples

### SSM parameter basic info
Explore the basic information of AWS SSM parameters to understand their types, data types, tiers, and the regions they are located in. This can help in managing and organizing these parameters efficiently.

```sql+postgres
select
  name,
  type,
  data_type,
  tier,
  region
from
  aws_ssm_parameter;
```

```sql+sqlite
select
  name,
  type,
  data_type,
  tier,
  region
from
  aws_ssm_parameter;
```


### Policy details of advanced tier ssm parameter
Explore the policy details of advanced tier parameters within AWS's Simple Systems Manager (SSM). This query can be used to understand the policy type, status, and text, providing valuable insights into the configuration and usage of these parameters.

```sql+postgres
select
  name,
  tier,
  p ->> 'PolicyType' as policy_type,
  p ->> 'PolicyStatus' as Policy_status,
  p ->> 'PolicyText' as policy_text
from
  aws_ssm_parameter,
  jsonb_array_elements(policies) as p;
```

```sql+sqlite
select
  name,
  tier,
  json_extract(p.value, '$.PolicyType') as policy_type,
  json_extract(p.value, '$.PolicyStatus') as policy_status,
  json_extract(p.value, '$.PolicyText') as policy_text
from
  aws_ssm_parameter,
  json_each(policies) as p;
```


### List of SSM parameters which do not have owner or app_id tag key
Determine the areas in which AWS SSM parameters are missing essential tags such as 'owner' or 'app_id'. This is useful in identifying potential gaps in your tagging strategy, which could impact resource management and cost allocation.

```sql+postgres
select
  name
from
  aws_ssm_parameter
where
  tags -> 'owner' is null
  or tags -> 'app_id' is null;
```

```sql+sqlite
select
  name
from
  aws_ssm_parameter
where
  json_extract(tags, '$.owner') is null
  or json_extract(tags, '$.app_id') is null;
```