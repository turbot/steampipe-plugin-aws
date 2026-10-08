
# **AWS Endpoint Generator**

## **Overview**
The **AWS Endpoint Generator** is a python script that **automatically fetches AWS service endpoint data** from the **AWS SDK for Go v2** and generates Go files containing:

- **Service ID constants** for programmatic access to AWS service metadata.

This tool helps determine **AWS service availability** across various regions and partitions.

---

## **How It Works**
The tool retrieves endpoint data from the following AWS SDK JSON file:  
🔗 **[AWS Endpoints JSON](https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json)**  

It then **parses the JSON**, extracts relevant data, and generates a Go file:

📌 **`aws/endpoint_service_ids_gen.go`** → Contains service ID constants for all AWS services.

---

## **Installation & Usage**
### **1️⃣ Navigate to the `steampipe-plugin-aws` Directory**
```sh
cd steampipe-plugin-aws
```

### **2️⃣ Install Dependencies**
Ensure **Python 3** is installed:
```sh
brew install python3  # For macOS users
```

### **3️⃣ Run the AWS Endpoint Generator**
```sh
./scripts/generate_aws_supported_endpoints/generate.sh
```

---

## **🔧 Troubleshooting**
If you encounter the following **error**:
```plaintext
error: externally-managed-environment

× This environment is externally managed
╰─> To install Python packages system-wide, try brew install
    xyz, where xyz is the package you are trying to install.
```
Try **activating a virtual environment** before running the script:

```sh
python3 -m venv venv  # Create a virtual environment
source venv/bin/activate  # Activate the virtual environment (For macOS/Linux)
```
For **Windows users**, activate the virtual environment with:
```sh
venv\Scripts\activate
```

Once activated, run the script again:
```sh
./scripts/generate_aws_supported_endpoints/generate.sh
```

Exit the virtual environment(venv)
```sh
deactivate
```

---

## **Custom Endpoints**
AWS SDK for Go v2 no longer adds newer services to its `endpoints.json` (the SDK team treats that file as internal), so some services the plugin supports have no upstream entry. Those are maintained by hand in `custom_endpoints.json`, keyed by service ID, partition, and region:

```json
"bedrock-agentcore": {
  "aws": {
    "us-east-1": ["bedrock-agentcore-control.us-east-1.amazonaws.com"]
  }
}
```

The generator merges this file into the downloaded `endpoints.json` (each merged hostname is tagged `custom`) before generating the Go constants. The merge skips hostnames that are already present, so re-running it is safe.

To add a service, get its region list from SSM and add an entry:
```sh
aws ssm get-parameters-by-path --path /aws/service/global-infrastructure/services/<service-id>/regions --query 'Parameters[].Value' --output text
```

Also add a service here if upstream removes it while the plugin still has a table for it, otherwise its `AWS_<SERVICE>_SERVICE_ID` constant disappears and the build breaks. If the service itself has been retired by AWS, remove its table instead.

---

## **Generated Files**
The tool generates the following files:

### **📌 1. `aws/endpoint_service_ids_gen.go` (Service ID Constants)**
- Defines **constants** for AWS services.
- These constants allow referencing AWS services programmatically.

Example:
```go
const AWS_S3_SERVICE_ID = "s3"
const AWS_EC2_SERVICE_ID = "ec2"
const AWS_LAMBDA_SERVICE_ID = "lambda"
```

### **📌 2. `internal/aws_endpoint_generator/endpoints.json`**
- This file contains AWS **endpoint metadata**.
- It is **automatically downloaded** and embedded into the project via `embed_endpoint.go`.

---
