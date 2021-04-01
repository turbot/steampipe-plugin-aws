select name, description, tags_src
from aws.aws_wafv2_regex_pattern_set
where id = '{{ output.resource_id_regional.value }}';