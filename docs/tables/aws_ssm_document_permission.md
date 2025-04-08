---
title: "Steampipe Table: aws_ssm_document_permission - Query AWS SSM Document Permissions using SQL"
description: "Allows users to query AWS SSM Document Permissions, providing detailed information about the permissions associated with Systems Manager (SSM) documents."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_document_permission - Query AWS SSM Document Permissions using SQL

The AWS SSM Document Permission is a component of AWS Systems Manager that allows you to manage access permissions to your SSM documents. These documents define the actions that Systems Manager performs on your managed instances. With document permissions, you can specify which AWS Identity and Access Management (IAM) users and roles can use these documents, providing a secure way to distribute commands or configurations to your servers.

## Table Usage Guide

The `aws_ssm_document_permission` table in Steampipe provides you with information about the permissions associated with AWS Systems Manager (SSM) documents. This table allows you, as a DevOps engineer, security analyst, or system administrator, to query document-specific permission details, including the type of permission, the account IDs that the permissions apply to, and the document version. You can utilize this table to gather insights on document permissions, such as identifying the accounts that have access to specific SSM documents, verifying the type of access granted, and more. The schema outlines the various attributes of the SSM document permission for you, including the document name, permission type, account IDs, and the document version.

**Important Notes**
- You must specify the `document_name` column in the `where` clause to query the table.

## Examples

### Basic info
Explore which AWS accounts have permission to the 'ConfigureS3BucketLogging' document. This can be useful to ensure only the intended accounts have access, enhancing security and compliance.

```sql+postgres
select
  document_name,
  shared_account_id,
  shared_document_version,
  account_ids,
  title
from
  aws_ssm_document_permission
where
  document_name = 'ConfigureS3BucketLogging';
```

```sql+sqlite
select
  document_name,
  shared_account_id,
  shared_document_version,
  account_ids,
  title
from
  aws_ssm_document_permission
where
  document_name = 'ConfigureS3BucketLogging';
```

### Get document details for the permissions
This query is useful for exploring the permissions and versions of a specific document in the AWS SSM service, in this case 'ConfigureS3BucketLogging'. It helps you understand who has access to the document and what versions of the document are approved, providing insights into document management and control.

```sql+postgres
select
  p.document_name,
  p.shared_account_id,
  p.shared_document_version,
  d.approved_version,
  d.attachments_information,
  d.created_date,
  d.default_version
from
  aws_ssm_document_permission as p,
  aws_ssm_document as d
where
  p.document_name = 'ConfigureS3BucketLogging';
```

```sql+sqlite
select
  p.document_name,
  p.shared_account_id,
  p.shared_document_version,
  d.approved_version,
  d.attachments_information,
  d.created_date,
  d.default_version
from
  aws_ssm_document_permission as p,
  aws_ssm_document as d
where
  p.document_name = 'ConfigureS3BucketLogging';
```