# Table: aws_ssm_document

AWS Systems Manager Document defines the actions that SSM performs on managed instances. SSM provides more than 100 pre-configured documents that used by specifying parameters at runtime.

## Examples

### SSM document basic info

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


### List of documents which does not have default version

```sql
select
    name,
    document_version,
    status,
    document_format,
    document_type,
    platform_types
from
    aws_ssm_document
where
    document_version != '1';
```


### List of documents which are not owned by Amazon

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
    owner != 'Amazon';
```


### List of documents those have default target type

```sql
select
    name,
    owner,
    document_format,
    document_type,
    target_type
from
    aws_ssm_document
where
    target_type is null;
```


### List of documents with empty platform types

```sql
select
    name,
    owner,
    document_format,
    document_type,
    target_type,
    platform_types
from
    aws_ssm_document
where
    platform_types = '[]';
```
