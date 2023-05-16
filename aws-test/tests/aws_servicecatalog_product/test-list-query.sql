select name
from aws.aws_servicecatalog_product
where name = '{{ resourceName }}';
