select
  decision
from
  aws.aws_iam_policy_simulator
where
    action = 'ec2:DescribeInstances'
  and
  resource_arn = '*'
  and
  principal_arn = '{{ output.user1_arn.value }}'