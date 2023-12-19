resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "aws securityhub delete-insight --insight-arn {{ output.insight_arn.value }} --profile default"
  }
}
