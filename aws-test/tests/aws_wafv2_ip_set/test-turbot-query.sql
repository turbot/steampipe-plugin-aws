select title, akas, tags
from aws.aws_wafv2_ip_set
where id = '{{ output.resource_id_regional.value }}';