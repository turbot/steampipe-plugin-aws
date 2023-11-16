---
title: "Table: aws_ssm_document_permission - Query AWS SSM Document Permissions using SQL"
description: "Allows users to query AWS SSM Document Permissions, providing detailed information about the permissions associated with Systems Manager (SSM) documents."
---

# Table: aws_ssm_document_permission - Query AWS SSM Document Permissions using SQL

The `aws_ssm_document_permission` table in Steampipe provides information about the permissions associated with AWS Systems Manager (SSM) documents. This table allows DevOps engineers, security analysts, and system administrators to query document-specific permission details, including the type of permission, the account IDs that the permissions apply to, and the document version. Users can utilize this table to gather insights on document permissions, such as identifying the accounts that have access to specific SSM documents, verifying the type of access granted, and more. The schema outlines the various attributes of the SSM document permission, including the document name, permission type, account IDs, and the document version.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_document_permission` table, you can use the `.inspect aws_ssm_document_permission` command in Steampipe.

**Key columns**:

- `name`: The name of the Systems Manager document. This column is useful for joining with other tables that contain document-specific information.
- `permission_type`: The type of permission associated with the document. This column is important as it helps in identifying the level of access granted to the accounts.
- `account_ids`: The AWS account IDs that the document permissions apply to. This column is useful for correlating with other tables that contain account-specific information.

## Examples

### Basic info

```sql
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

```sql
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
