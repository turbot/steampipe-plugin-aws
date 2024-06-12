## In order to run the integration test for the AWS plugin, you need to have the following:

### Install AWS CLI

#### For Mac, you can install AWS CLI using Homebrew:

```bash
brew install awscli
aws configure
```

#### For Linux:

You can install AWS CLI using the following command:

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
aws configure
```

#### For Windows:

You can download the binary from the [official website](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)

### Install Terraform

#### For Mac, you can install Terraform using Homebrew:

```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

or you can download the binary from the [official website](https://developer.hashicorp.com/terraform/install#darwin)

#### For Windows:

You can download the binary from the [official website](https://developer.hashicorp.com/terraform/install#windows)

#### For Linux:

```bash
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform
```

### Install Steampipe

#### For Mac, you can install Steampipe using Homebrew:

```bash
brew install steampipe
```

#### For Linux:

You can install Steampipe using the following command:

```bash
sudo /bin/sh -c "$(curl -fsSL https://steampipe.io/install/steampipe.sh)"
```

#### For Windows:

You can download the binary from the [official website](https://steampipe.io/downloads?install=windows)

### Install Steampipe Plugin

```bash
steampipe plugin install steampipe
```

### Install Steampipe AWS Plugin

```bash
steampipe plugin install aws
```

### Git clone steampipe-plugin-aws repository

```bash
git clone https://github.com/turbot/steampipe-plugin-aws.git
```

### Change the directory to steampipe-plugin-aws > aws-test directory

```bash
cd steampipe-plugin-aws/aws-test
```

### Install the node modules

```bash
npm install
```

### Run the integration test

```bash
./tint.js <table_name>
```