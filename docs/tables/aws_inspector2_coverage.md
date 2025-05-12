---
title: "Steampipe Table: aws_inspector2_coverage - Query AWS Inspector Coverage using SQL"
description: "Allows users to query AWS Inspector Coverage data that provides details on the assessment targets and the assessment templates that are associated with the AWS Inspector service."
folder: "Inspector2"
---

# Table: aws_inspector2_coverage - Query AWS Inspector Coverage using SQL

The AWS Inspector is a service that helps you improve the security and compliance of applications deployed on AWS. It automatically assesses applications for vulnerabilities or deviations from best practices, including impacted networks, and insecure configurations. Inspector provides a detailed list of security findings prioritized by level of severity, enabling you to identify potential security issues and areas for improvement effectively.

## Table Usage Guide

The `aws_inspector2_coverage` table in Steampipe provides you with information about the coverage of AWS Inspector within your AWS account. This table allows you, as a DevOps engineer, to query details about the assessment targets and the assessment templates that are associated with the AWS Inspector service. You can utilize this table to gather insights on the coverage of the AWS Inspector service, such as the number of assessment targets and templates, the ARN of the assessment targets and templates, and more. The schema outlines the various attributes of the AWS Inspector Coverage for you, including the ARN, name, duration, rules package ARNs, and user attributes for the assessment target and template.

## Examples

### Basic info
Explore the status and details of security inspections within your AWS resources to understand where potential vulnerabilities may exist. This query is useful for gaining insights into the security health of your resources and identifying areas for improvement.

```sql+postgres
select
  source_account_id,
  resource_id,
  resource_type,
  scan_type,
  scan_status_reason,
  scan_status_code
from
  aws_inspector2_coverage;
```

```sql+sqlite
select
  source_account_id,
  resource_id,
  resource_type,
  scan_type,
  scan_status_reason,
  scan_status_code
from
  aws_inspector2_coverage;
```

### List coverages that are inactive
This example can be used to identify the inactive coverage areas within your AWS Inspector service. It helps in pinpointing the specific locations where the scan status is inactive, allowing you to focus on reactivating these areas to ensure comprehensive security coverage.

```sql+postgres
select
  source_account_id,
  resource_id,
  resource_type,
  scan_status_code
from
  aws_inspector2_coverage
where
  scan_status_code = 'INACTIVE';
```

```sql+sqlite
select
  source_account_id,
  resource_id,
  resource_type,
  scan_status_code
from
  aws_inspector2_coverage
where
  scan_status_code = 'INACTIVE';
```

### List EC2 resource type coverage
Determine the areas in which EC2 resources are covered by examining the types of resources being scanned. This is particularly useful to ensure all necessary resources are included in security inspections.

```sql+postgres
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  ec2_ami_id,
  ec2_platform
from
  aws_inspector2_coverage
where
  resource_type = 'aws_EC2_INSTANCE';
```

```sql+sqlite
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  ec2_ami_id,
  ec2_platform
from
  aws_inspector2_coverage
where
  resource_type = 'aws_EC2_INSTANCE';
```

### List coverages by EC2 instance tags
Discover the segments that are covered by EC2 instance tags in your AWS account. This is useful for understanding your resource configuration and identifying any instances that may be tagged incorrectly or inconsistently.

```sql+postgres
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  ec2_ami_id,
  ec2_platform,
  ec2_instance_tags
from
  aws_inspector2_coverage
where
  ec2_instance_tags = '{"foo": "bar", "foo1": "bar1"}';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List coverages by lambda function tags
This query is used to examine the coverages associated with specific Lambda function tags within the AWS Inspector service. It can be useful for pinpointing the specific instances where these tagged functions are utilized, facilitating more efficient resource management and inspection.

```sql+postgres
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  lambda_function_name,
  lambda_function_runtime,
  lambda_function_tags
from
  aws_inspector2_coverage
where
  lambda_function_tags = '{"foo": "bar", "foo1": "bar1"}';
```

```sql+sqlite
select
  source_account_id,
  scan_type,
  resource_id as ec2_instance_id,
  resource_type,
  lambda_function_name,
  lambda_function_runtime,
  lambda_function_tags
from
  aws_inspector2_coverage
where
  json_extract(lambda_function_tags, '$.foo') = 'bar' and json_extract(lambda_function_tags, '$.foo1') = 'bar1';
```

### List coverage details of a package scan
Determine the areas in which a package scan has been performed within your AWS account. This can be useful for understanding the scope and reach of your security measures.

```sql+postgres
select
  source_account_id,
  resource_id,
  resource_type,
  scan_type
from
  aws_inspector2_coverage
where
  scan_type = 'PACKAGE';
```

```sql+sqlite
select
  source_account_id,
  resource_id,
  resource_type,
  scan_type
from
  aws_inspector2_coverage
where
  scan_type = 'PACKAGE';
```

### Get ECR repository details of each coverage
Determine the areas in which your Elastic Container Registry (ECR) repositories are covered by AWS Inspector. This is useful for understanding the extent of your security assessments and identifying any repositories that may be missing coverage.

```sql+postgres
select
  c.resource_id,
  c.resource_type,
  c.ecr_repository_name,
  r.registry_id,
  r.repository_uri,
  r.encryption_configuration
from
  aws_inspector2_coverage as c,
  aws_ecr_repository as r
where
  r.repository_name = c.ecr_repository_name
and
  c.resource_type = 'AWS_ECR_REPOSITORY';
```

```sql+sqlite
select
  c.resource_id,
  c.resource_type,
  c.ecr_repository_name,
  r.registry_id,
  r.repository_uri,
  r.encryption_configuration
from
  aws_inspector2_coverage as c,
  aws_ecr_repository as r
where
  r.repository_name = c.ecr_repository_name
and
  c.resource_type = 'AWS_ECR_REPOSITORY';
```

### Get lambda function details of each coverage
This query is used to gain insights into the details of each Lambda function within your AWS Inspector coverage. It allows you to understand the specifications of your functions, such as the runtime environment and code size, which can be useful in optimizing your resources and ensuring the security of your cloud infrastructure.

```sql+postgres
select
  c.resource_id,
  c.resource_type,
  c.lambda_function_name,
  f.arn as lambda_function_arn,
  c.lambda_function_runtime,
  f.code_sha_256,
  f.code_size,
  f.kms_key_arn,
  f.package_type
from
  aws_inspector2_coverage as c,
  aws_lambda_function as f
where
  f.name = c.lambda_function_name;
```

```sql+sqlite
select
  c.resource_id,
  c.resource_type,
  c.lambda_function_name,
  f.arn as lambda_function_arn,
  c.lambda_function_runtime,
  f.code_sha_256,
  f.code_size,
  f.kms_key_arn,
  f.package_type
from
  aws_inspector2_coverage as c,
  aws_lambda_function as f
where
  f.name = c.lambda_function_name;
```

### Get EC2 instance details of each coverage
Determine the areas in which the EC2 instances are covered by the AWS Inspector. This allows you to understand the security and compliance status of your instances, helping to maintain optimal configurations and avoid potential vulnerabilities.

```sql+postgres
select
  c.resource_id as ec2_instance_id,
  c.resource_type,
  c.ec2_ami_id,
  i.instance_type,
  i.instance_state,
  i.disable_api_termination,
  i.ebs_optimized
from
  aws_inspector2_coverage as c,
  aws_ec2_instance as i
where
  i.instance_id = c.resource_id
and
  c.resource_type = 'AWS_EC2_INSTANCE';
```

```sql+sqlite
select
  c.resource_id as ec2_instance_id,
  c.resource_type,
  c.ec2_ami_id,
  i.instance_type,
  i.instance_state,
  i.disable_api_termination,
  i.ebs_optimized
from
  aws_inspector2_coverage as c
join
  aws_ec2_instance as i
on
  i.instance_id = c.resource_id
where
  c.resource_type = 'AWS_EC2_INSTANCE';
```

