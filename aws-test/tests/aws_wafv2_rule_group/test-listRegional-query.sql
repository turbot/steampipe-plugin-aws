select name, arn, id, scope, region
from aws.aws_wafv2_rule_group
where akas::text = '["{{output.resource_aka_regional.value}}"]';