select name, product_id
from aws_servicecatalog_product
where product_id = '{{ output.resource_id.value }}';
