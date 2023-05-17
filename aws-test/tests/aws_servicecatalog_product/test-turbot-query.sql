select name, title, akas, region, account_id
from aws.aws_servicecatalog_product
where name = '{{ resourceName }}' and region = '{{ output.aws_region.value }}';
