When region specified in aws.spc - *

Welcome to Steampipe v0.18.5
For more information, type .help
> .cache clear
> .cache off
> .search_path a_role_aab_without_mfa
> select
  placement_availability_zone as az,
  instance_type,
  count(*)
from
  aws_ec2_instance
group by
  placement_availability_zone,
  instance_type;
+----+---------------+-------+
| az | instance_type | count |
+----+---------------+-------+
+----+---------------+-------+
> select
  instance_id,
  monitoring_state
from
  aws_ec2_instance
where
  monitoring_state = 'disabled';
+-------------+------------------+
| instance_id | monitoring_state |
+-------------+------------------+
+-------------+------------------+
> select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed;

Error: operation error IAM: ListPolicies, failed to sign request: failed to retrieve credentials: failed to refresh cached credentials, operation error STS: AssumeRole, failed to resolve service endpoint, an AWS region is required, but was not found (SQLSTATE HV000)

+------+-----+
| name | arn |
+------+-----+
+------+-----+
> select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed
  and path = '/turbot/' limit 10;

Error: operation error IAM: ListPolicies, failed to sign request: failed to retrieve credentials: failed to refresh cached credentials, operation error STS: AssumeRole, failed to resolve service endpoint, an AWS region is required, but was not found (SQLSTATE HV000)

+------+-----+
| name | arn |
+------+-----+
+------+-----+