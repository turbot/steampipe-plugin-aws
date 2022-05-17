# Table: aws_securityhub_insight

An AWS Security Hub insight is a collection of related findings. It identifies a security area that requires attention and intervention. For example, an insight might point out EC2 instances that are the subject of findings that detect poor security practices. An insight brings together findings from across finding providers.

## Examples

### Basic info

```sql
select
  name,
  insight_arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight;
```

### List insights by a particular attribute

```sql
select
  name,
  insight_arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight
where
  group_by_attribute='ResourceId';
```

### Get workflow status details for a particular insight

```sql
select
  name,
  insight_arn,
  filters ->> 'WorkflowStatus' as workflow_status
from
  aws_securityhub_insight
where
  name = 'sp';
```

### Get record state details for a particular insight

```sql
select
  name,
  insight_arn,
  filters ->> 'RecordState' as record_state
from
  aws_securityhub_insight
where
  name = 'sp';
```