# Table: aws_eks_cluster

Amazon Elastic Kubernetes Service (EKS) is a managed Kubernetes service that makes it easy to run Kubernetes on AWS and on-premises.

## Examples

### Basic info

```sql
select
  name,
  arn,
  endpoint,
  identity,
  status
from
  aws_eks_cluster;
```


### List the Vpc configuration used by cluster

```sql
select
  name,
  resources_vpc_config ->> 'ClusterSecurityGroupId' as cluster_security_group_id,
  resources_vpc_config ->> 'EndpointPrivateAccess' as endpoint_private_access,
  resources_vpc_config ->> 'EndpointPublicAccess' as endpoint_public_access,
  resources_vpc_config ->> 'PublicAccessCidrs' as public_access_cidrs,
  resources_vpc_config ->> 'SecurityGroupIds' as security_group_ids,
  resources_vpc_config -> 'SubnetIds' as subnet_ids,
  resources_vpc_config ->> 'VpcId' as vpc_id
from
  aws_eks_cluster;
```


### List of log Types that are not enabled for logging

```sql
select
  name,
  i -> 'Enabled' as enabled,
  i -> 'Types' as types
from
  aws_eks_cluster,
  jsonb_array_elements(logging -> 'ClusterLogging') as i
where
  i ->> 'Enabled' = 'false';
```


### List of Kubernetes Service IP address range of cluster

```sql
select
  name,
  kubernetes_network_config ->> 'ServiceIpv4Cidr' as service_ipv4_cidr
from
  aws_eks_cluster;
```


### List of certificate authority data for cluster

```sql
select
  name,
  certificate_authority ->> 'Data' as data
from
  aws_eks_cluster;
```