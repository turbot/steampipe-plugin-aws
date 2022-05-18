# Table: aws_securityhub_insight

An AWS Security Hub insight is a collection of related findings. It identifies a security area that requires attention and intervention. For example, an insight might point out EC2 instances that are the subject of findings that detect poor security practices. An insight brings together findings from across finding providers.

## Examples

### Basic info

```sql
select
  name,
  arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight;
```

### List insights by a particular attribute

```sql
select
  name,
  arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight
where
  group_by_attribute = 'ResourceId';
```

### Get workflow status details for insights

```sql
select
  name,
  arn,
  filters ->> 'WorkflowStatus' as workflow_status
from
  aws_securityhub_insight;
```

### Get severity details for insights

```sql
select
  name,
  arn,
  filters ->> 'SeverityLabel' as severity_label
from
  aws_securityhub_insight;
```

### Get malware details for insights

```sql
select
  name,
  arn,
  filters ->> 'MalwareName' as malware_name,
  filters ->> 'MalwarePath' as malware_path,
  filters ->> 'MalwareType' as malware_type
from
  aws_securityhub_insight;
```

### Get network details for insights

```sql
select
  name,
  arn,
  filters ->> 'NetworkSourceDomain' as network_source_domain,
  filters ->> 'NetworkDestinationDomain' as network_destination_domain,
  filters ->> 'NetworkSourceIpV4' as network_source_ip_v4,
  filters ->> 'NetworkDestinationIpV4' as network_destination_ip_v4,
  filters ->> 'NetworkSourceIpV6' as network_source_ip_v6,
  filters ->> 'NetworkDestinationIpV6' as network_destination_ip_v6,
  filters ->> 'NetworkSourcePort' as network_source_port,
  filters ->> 'NetworkDestinationPort' as network_destination_port
from
  aws_securityhub_insight;
```

### Get compliance details for insights

```sql
select
  name,
  arn,
  filters ->> 'ComplianceStatus' as compliance_status
from
  aws_securityhub_insight;
```

### Get record state details for a custom insight named 'sp'

```sql
select
  name,
  arn,
  filters ->> 'RecordState' as record_state
from
  aws_securityhub_insight
where
  name = 'sp';
```