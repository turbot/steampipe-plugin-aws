# Table: aws_audit_manager_assessment

AWS Audit Manager helps you continuously audit your AWS usage to simplify how you manage risk and compliance with regulations and industry standards. Audit Manager automates evidence collection to make it easier to assess whether your policies, procedures, and activities—also known as controls—are operating effectively. When it is time for an audit, Audit Manager helps you manage stakeholder reviews of your controls, which means you can build audit-ready reports with much less manual effort.

## Examples

### Basic info

```sql
select
  name,
  arn,
  status,
  compliance_type
from
  aws_audit_manager_assessment;
```


### List assessments with public audit bucket

```sql
select
  a.name,
  a.arn,
  a.assessment_report_destination,
  a.assessment_report_destination_type,
  b.bucket_policy_is_public as is_public_bucket
from
  aws_audit_manager_assessment as a
join aws_s3_bucket as b on a.assessment_report_destination = 's3://' || b.Name and b.bucket_policy_is_public;
```


### List inactive assessments

```sql
select
  name,
  arn,
  status
from
  aws_audit_manager_assessment
where
  status <> 'ACTIVE';
```
