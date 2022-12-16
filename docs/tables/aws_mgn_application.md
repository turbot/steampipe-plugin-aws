# Table: aws_mgn_application

AWS Application Migration Service (MGN) is a highly automated lift-and-shift (rehost) solution that simplifies, expedites, and reduces the cost of migrating applications to AWS. It enables companies to lift-and-shift a large number of physical, virtual, or cloud servers without compatibility issues, performance disruption, or long cutover windows.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  application_id,
  creation_date_time,
  is_archived,
  wave_id,
  tags
from
  aws_mgn_application;
```

### List archived applications

```sql
select
  name,
  arn,
  application_id,
  creation_date_time,
  is_archived
from
  aws_mgn_application
where
  is_archived;
```

### Get aggregated status details for a application

```sql
select
  name,
  application_id,
  application_aggregated_status ->> 'HealthStatus' as health_status,
  application_aggregated_status ->> 'ProgressStatus' as progress_status,
  application_aggregated_status ->> 'TotalSourceServers' as total_source_servers
from
  aws_mgn_application;
```

### List application that are created in last 30 days

```sql
select
  name,
  application_id,
  creation_date_time,
  is_archived,
  wave_id
from
  aws_mgn_application
where
  creation_date_time >= now() - interval '30' day;
```