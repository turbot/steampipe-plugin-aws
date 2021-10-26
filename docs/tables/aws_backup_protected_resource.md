# Table: aws_backup_protected_resource

AWS Backup is AWS’ native data protection platform that offers centralized, fully managed, and policy-based service to protect customer data and ensure compliance and business continuity across AWS services. With AWS Backup, you can centrally configure data protection (backup) policies and monitor backup activity across your AWS resources, including: Amazon EBS volumes, Amazon Relational Database Service (Amazon RDS) databases (including Aurora clusters), Amazon DynamoDB tables, Amazon Elastic File System (Amazon EFS), Amazon EC2 instances, AWS Storage Gateway volumes and Amazon FSx file systems.

## Examples

### Basic Info

```sql
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource;
```

### List EBS volumes that are backed up

```sql
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource
where
  resource_type = 'EBS';
```
