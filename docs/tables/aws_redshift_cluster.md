# Table: aws_redshift_cluster

A cluster is a fully managed data warehouse that consists of a set of compute nodes.

## Examples

### Basic info

```sql
select
  cluster_identifier,
  akas,
  node_type,
  region
from
  aws_redshift_cluster;
```


### List clusters that are publicly accessible

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  publicly_accessible
from
  aws_redshift_cluster
where
  publicly_accessible;
```


### List clusters that are not in a VPC

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  vpc_id
from
  aws_redshift_cluster
where
  vpc_id is null;
```


### List clusters whose storage is not encrypted

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  encrypted
from
  aws_redshift_cluster
where
  not encrypted;
```
