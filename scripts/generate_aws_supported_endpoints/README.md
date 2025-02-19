
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
