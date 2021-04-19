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
  created,
  deleted
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
  volume_size,
  volume_type,
  ebs_enabled
from
  aws_elasticsearch_domain
where
  ebs_enabled = true;
```


### Get details of vpc associated with domain

```sql
select
  domain_name,
  availability_zones,
  security_group_ids,
  subnet_ids,
  vpc_id
from
  aws_elasticsearch_domain
where
  availability_zones is not null;
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