select akas, name, region, tags, title
from aws.aws_accessanalyzer_analyzer
where name = '{{ resourceName }}';
