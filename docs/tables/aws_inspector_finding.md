# Table: aws_inspector_finding

Amazon Inspector offers several ways to sort, group, and manage your findings. These features help you tailor findings to your environment, aggregate findings by different views, and focus on vulnerabilities to your specific AWS environment.

## Examples

### Basic info

```sql
select
  id,
  arn,
  agent_id as instance_id,
  asset_type,
  confidence,
  severity
from
  aws_inspector_finding;
```

### List findings with high severity

```sql
select
  id,
  arn,
  agent_id as instance_id,
  asset_type,
  confidence,
  severity
from
  aws_inspector_finding
where
  severity = 'High';
```

### Count the number of findings by severity

```sql
select
  severity,
  count(severity)
from
  aws_inspector_finding
group by
  severity
order by
  severity;
```

### List last 10 days findings

```sql
select
  title,
  id,
  confidence,
  severity
from
  aws_inspector_finding
where
  created_at >= now() - interval '10' day;
```

### List attributes for each finding

```sql
select
  title,
  id, 
  jsonb_pretty(attributes) as attributes
from
  aws_inspector_finding;
```

### Get asset attributes for each finding

```sql
select
  id,
  title,
  asset_attributes ->> 'AgentId' as agent_id,
  asset_attributes ->> 'AmiId' as ami_id,
  asset_attributes ->> 'Hostname' as hostname,
  asset_attributes ->> 'Tags' as tags
from
  aws_inspector_finding;
```

### List EC2 instances with high severity

```sql
select
  distinct i.instance_id,
  i.instance_state,
  i.instance_type,
  f.title,
  f.service,
  f.severity,
  f.confidence
from
  aws_ec2_instance as i,
  aws_inspector_finding as f
where
  severity = 'High'
and
  i.instance_id = f.agent_id;
```

### Get service attributes for each finding

```sql
select
  id,
  title,
  service_attributes ->> 'AssessmentRunArn' as assessment_run_arn,
  service_attributes ->> 'RulesPackageArn' as rules_package_arn,
  service_attributes ->> 'SchemaVersion' as schema_version,
from
  aws_inspector_finding;
```

### Get assessment run details for findings

```sql
select
  f.id,
  r.title,
  f.service_attributes ->> 'AssessmentRunArn' as assessment_run_arn,
  r.assessment_template_arn,
  r.finding_counts
from
  aws_inspector_finding as f,
  aws_inspector_assessment_run as r
where
  f.service_attributes ->> 'AssessmentRunArn' = r.arn;
```

### List findings order by confidence

```sql
select
  id,
  arn,
  agent_id as instance_id,
  asset_type,
  confidence,
  severity
from
  aws_inspector_finding
order by
  confidence;
```
