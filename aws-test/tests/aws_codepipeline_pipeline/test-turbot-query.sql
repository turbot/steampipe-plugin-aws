select name, title, account_id, region, akas
from aws_codepipeline_pipeline
where name = '{{ resourceName }}';