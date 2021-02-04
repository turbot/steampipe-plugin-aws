select
  decision
from
  aws.aws_iam_policy_simulator
where
    action = 'ec2:Describe*'
  and
  resource_arn = '*'
  and
  principal_arn = '{{ output.user2_arn.value }}'