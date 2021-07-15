# Table: aws_codecommit_repository

AWS CodeCommit repository is the fundamental version control object in CodeCommit. It's where you securely store code and files for your project. It also stores your project history, from the first commit through the latest changes. You can share your repository with other users so you can work together on a project.

## Examples

### Basic info

```sql
select
  repository_name,
  repository_id,
  arn,
  creation_date,
  region
from
  aws_codecommit_repository;
```
