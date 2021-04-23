# Table: aws_accessanalyzer_analyzer

AWS IAM Access Analyzer helps you identify the resources in your organization and accounts, such as Amazon S3 buckets or IAM roles, that are shared with an external entity. This lets you identify unintended access to your resources and data, which is a security risk. Access Analyzer identifies resources that are shared with external principals by using logic-based reasoning to analyze the resource-based policies in your AWS environment.

## Examples

### Basic info

```sql
select
  name,
  last_resource_analyzed,
  last_resource_analyzed_at,
  status,
  type
from
  aws_accessanalyzer_analyzer;
```


### List analyzers which are enabled

```sql
select
  name,
  last_resource_analyzed,
  last_resource_analyzed_at,
  tags
from
  aws_accessanalyzer_analyzer
where
  status = 'ACTIVE';
```


### List analyzers based on type

```sql
select
  name,
  last_resource_analyzed,
  last_resource_analyzed_at,
  tags
from
  aws_accessanalyzer_analyzer
where
  type = 'ACCOUNT';
```