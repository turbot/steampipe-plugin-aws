# Table: aws_sim_space_weaver_app

The SimSpace Weaver app SDK provides APIs that you can use to control the entities in your simulation and respond to SimSpace Weaver events.

## Examples

### Basic info

```sql
select
  name,
  simulation,
  domain,
  status,
  target_status
from
  aws_sim_space_weaver_app;
```

### List apps that are in error state

```sql
select
  name,
  simulation,
  domain,
  status
from
  aws_sim_space_weaver_app
where
  status = 'ERROR';
```

### Get simulation details of each app

```sql
select
  a.name,
  a.status
  a.simulation,
  s.arn as simulation_arn,
  s.status as simulation_status,
  s.logging_configuration
from
  aws_sim_space_weaver_app as a,
  aws_sim_space_weaver_simulation as s
where
  a.name = s.name;
```