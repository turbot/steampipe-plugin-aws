---
title: "Steampipe Table: aws_config_conformance_pack - Query AWS Config Conformance Packs using SQL"
description: "Allows users to query AWS Config Conformance Packs to fetch information about the AWS Config conformance packs deployed on an AWS account."
folder: "Config"
---

# Table: aws_config_conformance_pack - Query AWS Config Conformance Packs using SQL

The AWS Config Conformance Pack is a collection of AWS Config rules and remediation actions that can be easily deployed as a single entity in an account and a region. These packs can be used to create a common baseline of security, compliance, or operational best practices across multiple accounts in your organization. AWS Config continuously monitors and records your AWS resource configurations and allows you to automate the evaluation of recorded configurations against desired configurations.

## Table Usage Guide

The `aws_config_conformance_pack` table in Steampipe provides you with information about AWS Config conformance packs within the AWS Config service. This table allows you, as a DevOps engineer, to query conformance pack-specific details, including pack names, delivery S3 bucket, and associated metadata. You can utilize this table to gather insights on conformance packs, such as pack ARN, creation time, last update requested time, input parameters, and more. The schema outlines the various attributes of the conformance pack for you, including the pack ARN, delivery S3 bucket, input parameters, and associated tags.

## Examples

### Basic info
Explore the general information about AWS Config Conformance Packs, such as who created them and when they were last updated. This can help understand the management and status of these resources in your AWS environment.

```sql+postgres
select
  name,
  conformance_pack_id,
  created_by,
  last_update_requested_time,
  title,
  akas
from
  aws_config_conformance_pack;
```

```sql+sqlite
select
  name,
  conformance_pack_id,
  created_by,
  last_update_requested_time,
  title,
  akas
from
  aws_config_conformance_pack;
```


### Get S3 bucket info for each conformance pack
Explore which conformance packs are associated with each S3 bucket. This can help streamline and improve the management of AWS configurations.

```sql+postgres
select
  name,
  conformance_pack_id,
  delivery_s3_bucket,
  delivery_s3_key_prefix
from
  aws_config_conformance_pack;
```

```sql+sqlite
select
  name,
  conformance_pack_id,
  delivery_s3_bucket,
  delivery_s3_key_prefix
from
  aws_config_conformance_pack;
```


### Get input parameter details of each conformance pack
Determine the settings of each conformance pack in your AWS Config service. This helps in understanding how each pack is configured and can assist in identifying any discrepancies or areas for optimization.

```sql+postgres
select
  name,
  inp ->> 'ParameterName' as parameter_name,
  inp ->> 'ParameterValue' as parameter_value,
  title,
  akas
from
  aws_config_conformance_pack,
  jsonb_array_elements(input_parameters) as inp;
```

```sql+sqlite
select
  aws_config_conformance_pack.name,
  json_extract(inp.value, '$.ParameterName') as parameter_name,
  json_extract(inp.value, '$.ParameterValue') as parameter_value,
  title,
  akas
from
  aws_config_conformance_pack,
  json_each(input_parameters) as inp;
```