select security_group_rule_id, arn, group_id, description, group_owner_id
from aws.aws_vpc_security_group_rules
where group_id = '{{ output.resource_id.value }}' and security_group_rule_id = '{{ output.ingress_sgr_id.value }}'
