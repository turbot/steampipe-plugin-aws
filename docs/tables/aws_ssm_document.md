---
title: "Table: aws_ssm_document - Query AWS SSM Documents using SQL"
description: "Allows users to query AWS SSM Documents and retrieve detailed information about each document, including its name, version, owner, status, and permissions, among others."
---

# Table: aws_ssm_document - Query AWS SSM Documents using SQL

The `aws_ssm_document` table in Steampipe provides information about SSM documents within AWS Systems Manager (SSM). This table allows DevOps engineers to query document-specific details, including the document name, version, owner, status, and permissions, among others. Users can utilize this table to gather insights on SSM documents, such as their current status, the document format, and the permissions associated with each document. The schema outlines the various attributes of the SSM document, including the document name, document version, owner, permissions, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_document` table, you can use the `.inspect aws_ssm_document` command in Steampipe.

### Key columns:

- `name`: The name of the SSM document. This can be used to join with other tables that contain information about SSM documents.
- `document_version`: The version of the SSM document. This is useful for tracking the various versions of a document.
- `owner`: The AWS user who owns the SSM document. This can be used to join with other tables that contain information about AWS users.

## Examples

### Basic info

```sql
select
  name,
  document_version,
  status,
  owner,
  document_format,
  document_type,
  platform_types,
  region
from
  aws_ssm_document;
```

### List documents owned by the AWS account

```sql
select
  name,
  owner,
  document_version,
  status,
  document_format,
  document_type
from
  aws_ssm_document
where
  owner_type = 'Self';
```

### List documents not owned by Amazon

```sql
select
  name,
  owner,
  document_version,
  status,
  document_format,
  document_type
from
  aws_ssm_document
where
  owner_type != 'Amazon';
```

### List documents that are shared publicly

```sql
select
  name,
  owner,
  account_ids
from
  aws_ssm_document
where
  owner_type = 'Self'
  and account_ids :: jsonb ? 'all';
```

### Get a specific document

```sql
select
  name,
  arn,
  approved_version,
  created_date,
  document_type
from
  aws_ssm_document
where
  arn = 'arn:aws:ssm:ap-south-1:112233445566:document/AWS-ASGEnterStandby';
```
