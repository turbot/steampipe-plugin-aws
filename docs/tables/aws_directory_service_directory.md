# Table: aws_directory_service_directory

AWS Directory Service for Microsoft Active Directory, also known as AWS Managed Microsoft Active Directory (AD), enables your directory-aware workloads and AWS resources to use managed Active Directory (AD) in AWS.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  directory_id
from
  aws_directory_service_directory;
```

### List Directories of  type MicrosoftAD

```sql
select
  name,
  arn,
  directory_id,
  type
from
  aws_directory_service_directory
where 
  type = 'MicrosoftAD';
```