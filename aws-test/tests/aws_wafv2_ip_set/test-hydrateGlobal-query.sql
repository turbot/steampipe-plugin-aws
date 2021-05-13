select name, description, tags_src
from aws.aws_wafv2_ip_set
where id = '{{ output.resource_id_global.value }}';