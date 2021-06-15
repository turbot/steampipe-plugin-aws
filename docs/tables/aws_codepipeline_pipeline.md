# Table: aws_codebuild_project

AWS CodePipeline is a fully managed continuous delivery service that helps you automate your release pipelines for fast and reliable application and infrastructure updates. ... You can easily integrate AWS CodePipeline with third-party services such as GitHub or with your own custom plugin.

## Examples

### Basic info

```sql
select
  name,
  arn,
  encryption_key,
  role_arn,
  stages,
  tags_src,
  title,
  created,
  updated,
  version
from
  aws_codepipeline_pipeline;
```

### List unencrypted pipelines

```sql
select
  name,
  arn,
  encryption_key
from
  aws_codepipeline_pipeline
where
  encryption_key is null;
```
