select
  arn,
  name,
  description,
  region
from
  aws_glue_job
where 
  akas::text = '["{{ output.resource_aka.value }}"]';
