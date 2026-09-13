select 
  name, 
  arn, 
  tags
from aws_tagging_resource
where tag_filters = '[{"key":"Name","values":["{{ output.resource_name.value }}"]}]'::jsonb;
