# Table: aws_elasticsearch_domain

Amazon ES is a managed service that helps to deploy,operate,scale Elasticsearch clusters in the AWS Cloud. ES domain is synonymous with an Elasticsearch cluster.

## Example

### Basic info

```sql
select
  domain_name,
  domain_id,
  arn,
  elasticsearch_version,
  created
from
  aws_elasticsearch_domain;
```


### List domains that are not encrypted at rest

```sql
select
  domain_name,
  domain_id,
  encryption_at_rest_options ->> 'Enabled' as enabled,
  encryption_at_rest_options ->> 'KmsKeyId' as kms_key_id
from
  aws_elasticsearch_domain
where
  encryption_at_rest_options ->> 'Enabled' = 'false';
```


### Get details of ebs associated with each domain

```sql
select
  domain_name,
  domain_id,
  ebs_options ->> 'VolumeSize' as volume_size,
  ebs_options ->> 'VolumeType' as volume_type,
  ebs_options ->> 'EBSEnabled' as ebs_enabled
from
  aws_elasticsearch_domain
where
  ebs_options ->> 'EBSEnabled' = 'true';
```


### Get network details associated with domain

```sql
select
  domain_name,
  vpc_options ->> 'AvailabilityZones' as availability_zones,
  vpc_options ->> 'SecurityGroupIds' as security_group_ids,
  vpc_options ->> 'SubnetIds' as subnet_ids,
  vpc_options ->> 'VPCId' as vpc_id
from
  aws_elasticsearch_domain
where
  vpc_options ->> 'AvailabilityZones' is not null;
```


### Get the details of instance associated with domain cluster

```sql
select
  domain_name,
  domain_id,
  elasticsearch_cluster_config ->> 'InstanceType' as instance_type,
  elasticsearch_cluster_config ->> 'InstanceCount' as instance_count
from
  aws_elasticsearch_domain;
```


### List of domains policy statements that grant anonymous access

```sql
select
  domain_name,
  p as principal,
  a as action,
  s ->> 'Effect' as effect
from
  aws_elasticsearch_domain,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  p = '*'
  and s ->> 'Effect' = 'Allow';
```


### List of domains which are plan for deletion

```sql
select
  domain_name,
  domain_id,
  deleted
from
  aws_elasticsearch_domain
where
  deleted is not false;
```