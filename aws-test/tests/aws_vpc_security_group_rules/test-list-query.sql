select security_group_rule_id, group_id
from aws.aws_vpc_security_group_rules
where group_id = '{{ output.resource_id.value }}'
