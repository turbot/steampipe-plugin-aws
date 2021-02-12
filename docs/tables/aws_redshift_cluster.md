# Table: aws_redshift_cluster

A cluster is a fully managed data warehouse that consists of a set of compute nodes.

## Examples

### List of clusters which are publicly accessible

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

### List of clusters which are not in a VPC

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

### List of clusters whose storage is not encrypted

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

### Endpoint info of each cluster

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  endpoint_address,
  endpoint_port,
  endpoint_vpc_endpoints
from
  aws_redshift_cluster;
```
