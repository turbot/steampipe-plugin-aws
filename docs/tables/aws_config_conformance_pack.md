---
title: "Table: aws_config_conformance_pack - Query AWS Config Conformance Packs using SQL"
description: "Allows users to query AWS Config Conformance Packs to fetch information about the AWS Config conformance packs deployed on an AWS account."
---

# Table: aws_config_conformance_pack - Query AWS Config Conformance Packs using SQL

The `aws_config_conformance_pack` table in Steampipe provides information about AWS Config conformance packs within AWS Config service. This table allows DevOps engineers to query conformance pack-specific details, including pack names, delivery S3 bucket, and associated metadata. Users can utilize this table to gather insights on conformance packs, such as pack ARN, creation time, last update requested time, input parameters, and more. The schema outlines the various attributes of the conformance pack, including the pack ARN, delivery S3 bucket, input parameters, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_config_conformance_pack` table, you can use the `.inspect aws_config_conformance_pack` command in Steampipe.

### Key columns:

- `name`: The name of the conformance pack. This column can be used to join with other tables to fetch more detailed information about the conformance pack.
- `conformance_pack_id`: The ID of the conformance pack. This is a unique identifier for the conformance pack and can be used to join with other tables.
- `delivery_s3_bucket`: The S3 bucket where AWS Config delivers the conformance pack template. This can be used to join with S3 related tables for more insights.

## Examples

### Basic info

```sql
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

```sql
select
  name,
  conformance_pack_id,
  delivery_s3_bucket,
  delivery_s3_key_prefix
from
  aws_config_conformance_pack;
```


### Get input parameter details of each conformance pack

```sql
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

