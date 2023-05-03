# Table: aws_ssm_document

AWS Systems Manager Document defines the actions that SSM performs on managed
instances. SSM provides more than 100 pre-configured documents that used by
specifying parameters at runtime.

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

### Get a document by region

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
  arn = 'arn:aws:ssm:ap-south-1:112233445566:document/AWS-ASGEnterStandby'
and
  region = 'ap-south-1';
```