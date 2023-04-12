# Table: aws_ssm_document_permission

Describes the permissions for a AWS Systems Manager document (SSM document). If you created the document, you are the owner. If a document is shared, it can either be shared privately (by specifying a user's AWS account ID) or publicly (All).

**Note**: `document_name` is required in query parameter to get the permission details of the document.

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
