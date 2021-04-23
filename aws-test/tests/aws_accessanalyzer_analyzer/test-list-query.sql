select name, status, type, region, tags, title
from aws.aws_accessanalyzer_analyzer
where akas::text = '["{{ output.resource_aka.value }}"]';
