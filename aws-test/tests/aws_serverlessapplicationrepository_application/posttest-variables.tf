resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "aws serverlessrepo delete-application --application-id {{ output.arn.value }}"
  }
}
