select name, arn, id, scope, region
from aws.aws_wafv2_regex_pattern_set
where akas::text = '["{{output.resource_aka_regional.value}}"]';