select title, akas, tags
from aws.aws_wafv2_regex_pattern_set
where id = '{{ output.resource_id_regional.value }}';