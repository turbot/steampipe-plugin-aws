select name, title, akas
from aws.aws_accessanalyzer_analyzer
where name = '{{ resourceName }}';
