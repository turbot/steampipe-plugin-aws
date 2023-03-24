# Table: aws_appstream_image

Amazon AppStream 2.0 contains applications you can stream to your users and default system and application settings to enable your users to get started with those applications quickly. However, after you create an image, you can't change it. To add other applications, update existing applications, or change image settings, you must start and reconnect to the image builder that you used to create the image. If you deleted the image builder, launch a new image builder that is based on your image.

## Examples

### Basic info

```sql
select
  name,
  arn,
  base_image_arn,
  description,
  created_time,
  display_name,
  image_builder_name,
  tags
from
  aws_appstream_image;
```

### List available images

```sql
select
  name,
  arn,
  display_name,
  platform,
  state
from
  aws_appstream_image
where
  state = 'AVAILABLE';
```

### List Windows based images

```sql
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_supported,
  image_builder_name
from
  aws_appstream_image
where
  platform = 'WINDOWS';
```

### List images that support image builder

```sql
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_supported,
  image_builder_name
from
  aws_appstream_image
where
  image_builder_supported;
```

### List private images

```sql
select
  name,
  created_time,
  base_image_arn,
  display_name,
  image_builder_name,
  visibility
from
  aws_appstream_image
where
  visibility = 'PRIVATE';
```

### Get application details of images

```sql
select
  name,
  arn,
  a ->> 'AppBlockArn' as app_block_arn,
  a ->> 'Arn' as app_arn,
  a ->> 'CreatedTime' as app_created_time,
  a ->> 'Description' as app_description,
  a ->> 'DisplayName' as app_display_name,
  a ->> 'Enabled' as app_enabled,
  a ->> 'IconS3Location' as app_icon_s3_location,
  a ->> 'IconURL' as app_icon_url,
  a ->> 'InstanceFamilies' as app_instance_families,
  a ->> 'LaunchParameters' as app_launch_parameters,
  a ->> 'LaunchPath' as app_launch_path,
  a ->> 'Name' as app_name,
  a ->> 'Platforms' as app_platforms,
  a ->> 'WorkingDirectory' as app_WorkingDirectory
from
  aws_appstream_image,
  jsonb_array_elements(applications) as a;
```

### Get the permission model of the images

```sql
select
  name,
  arn,
  image_permissions ->> 'AllowFleet' as allow_fleet,
  image_permissions ->> 'AllowImageBuilder' as allow_image_builder
from
  aws_appstream_image;
```

### Get error details of failed images

```sql
select
  name,
  arn,
  e ->> 'ErrorCode' as error_code,
  e ->> 'ErrorMessage' as error_message,
  e ->> 'ErrorTimestamp' as error_timestamp
from
  aws_appstream_image,
  jsonb_array_elements(image_errors) as e;
```
