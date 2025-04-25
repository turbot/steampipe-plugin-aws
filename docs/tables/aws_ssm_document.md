---
title: "Steampipe Table: aws_ssm_document - Query AWS SSM Documents using SQL"
description: "Allows users to query AWS SSM Documents and retrieve detailed information about each document, including its name, version, owner, status, and permissions, among others."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_document - Query AWS SSM Documents using SQL

The AWS Systems Manager Document (SSM Document) is a resource that defines the actions that Systems Manager performs on your managed instances. These documents can be used to automate tasks and ensure they are done consistently across multiple instances. SSM Documents support multiple types of actions, including running scripts, applying patches, and more, enabling you to manage your AWS resources effectively.

## Table Usage Guide

The `aws_ssm_document` table in Steampipe provides you with information about SSM documents within AWS Systems Manager (SSM). This table enables you, as a DevOps engineer, to query document-specific details, including the document name, version, owner, status, and permissions, among others. You can utilize this table to gather insights on SSM documents, such as their current status, the document format, and the permissions associated with each document. The schema outlines for you the various attributes of the SSM document, including the document name, document version, owner, permissions, and associated tags.

## Examples

### Basic info
This query allows you to explore and understand the status, owner, and platform details of documents within your AWS Simple Systems Manager. It helps in identifying where potential changes or updates may be necessary, providing insights for better system management.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which your AWS account owns documents. This can help you understand your resource ownership and manage your AWS resources more effectively.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which documents within the AWS SSM service are not owned by Amazon. This can be useful to identify potential security risks or to audit ownership of documents.

```sql+postgres
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

```sql+sqlite
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
  owner_type <> 'Amazon';
```

### List documents that are shared publicly
Discover the segments that consist of documents which are shared publicly. This query is handy in identifying potential security risks by pinpointing documents that are open to all, thus allowing for appropriate action to be taken.

```sql+postgres
with ssm_documents as (
  select
    name,
    owner,
    region,
    account_id
  from
    aws_ssm_document
  where
    owner_type = 'Self'
  order by
    name
)
select
  d.name,
  d.owner,
  p.account_ids
from
  ssm_documents as d
  left join aws_ssm_document_permission as p on p.document_name = d.name and p.region = d.region and p.account_id = d.account_id
where
  p.account_ids :: jsonb ? 'all';
```

```sql+sqlite
with ssm_documents as (
  select
    name,
    owner,
    region,
    account_id
  from
    aws_ssm_document
  where
    owner_type = 'Self'
  order by
    name
)
select
  d.name,
  d.owner,
  p.account_ids
from
  ssm_documents as d
  left join aws_ssm_document_permission as p on p.document_name = d.name 
    and p.region = d.region 
    and p.account_id = d.account_id
where
  json_extract(account_ids, '$.all') is not null;
```

### Get a specific document
This query allows users to pinpoint the specific details of a document within the AWS Simple Systems Manager (SSM), particularly useful for those needing to assess a document's approved version or creation date. It's particularly beneficial when managing or auditing AWS resources.

```sql+postgres
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

```sql+sqlite
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