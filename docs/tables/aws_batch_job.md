# Table: aws_batch_job

A unit of work (such as a shell script, a Linux executable, or a Docker container image) that you submit to AWS Batch. It has a name, and runs as a containerized application on AWS Fargate or Amazon EC2 resources in your compute environment, using parameters that you specify in a job definition. Jobs can reference other jobs by name or by ID, and can be dependent on the successful completion of other job

**Note:** Only one filter can be used at a time. When the filter is used, jobStatus is ignored

## Examples

### Basic Info

```sql
select
  job_id,
  job_name,
  status,
  created_at
from
  aws_batch_job;
```

### List jobs created in the last 30 days

```sql
select
  job_id,
  job_name,
  status,
  created_at
from
  aws_batch_job
where
  created_at > now() - interval '30 days';
```

### List all successful jobs

```sql
select
  job_id,
  job_name,
  created_at
from
  aws_batch_job
where
  status = 'succeeded';
```

### List all cancelled jobs

```sql
select
  job_id,
  job_name,
  status
from
  aws_batch_job
where
  is_cancelled = true;
```

###  List all failed jobs from the last week

```sql
select
  job_id,
  job_name,
  created_at
from
  aws_batch_job
where
  status = 'failed'
  and created_at > now() - interval '7 days';
```

### Get the multi-node properties of a particular job

```sql
select
  job_id,
  job_name,
  created_at.
  multi_node_properties -> 'MainNode' as main_node,
  multi_node_properties -> 'NodeRangeProperties' as node_range_properties,
  multi_node_properties -> 'NumNodes' as number_of_nodes
from
  aws_batch_job
where
  job_name = 'job-123';
```