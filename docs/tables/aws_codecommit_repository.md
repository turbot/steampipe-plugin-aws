# Table: aws_codecommit_repository

AWS CodeCommit is a secure, highly scalable, managed source control service that hosts private Git repositories. It makes it easy for teams to securely collaborate on code with contributions encrypted in transit and at rest. CodeCommit eliminates the need for you to manage your own source control system or worry about scaling its infrastructure. You can use CodeCommit to store anything from code to binaries. It supports the standard functionality of Git, so it works seamlessly with your existing Git-based tools.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  creation_date,
  region
from
  aws_codecommit_repository;
```
