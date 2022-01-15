select group_id, security_group_rule_id, akas, tags, title
from aws.aws_vpc_security_group_rules
where group_id = '{{ output.resource_id.value }}' and security_group_rule_id = '{{ output.ingress_sgr_id.value }}'
