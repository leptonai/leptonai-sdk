provider "aws" {
  region = var.region
}

provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    # This requires the awscli to be installed locally where Terraform is executed
    args = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
  }
}

provider "helm" {
  kubernetes {
    host                   = module.eks.cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      # This requires the awscli to be installed locally where Terraform is executed
      args = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
    }
  }
}

data "aws_availability_zones" "available" {}
data "aws_caller_identity" "current" {}
data "aws_iam_group" "dev_members" {
  group_name = "dev"
}

locals {
  vpc_cidr = "10.0.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  cluster_name = coalesce(var.cluster_name, "eks-${random_string.suffix.result}")

  vpc_cidr = "10.0.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.19.0"

  name = "vpc-${local.cluster_name}"

  cidr = local.vpc_cidr
  azs  = local.azs

  # VPC CIDR 10.0.0.0/16 ranges from 10.0.0.0 to 10.0.255.255 with 65534 IPs.
  #
  # We need more private subnets since we can just deploy load balancers in public subnets
  # and still route traffic to pods in private subnets.
  # ref. https://docs.aws.amazon.com/eks/latest/userguide/creating-a-vpc.html
  #
  # e.g., a reigon of 4 AZs will have:
  # 10.0.0.0/20  from 10.0.0.0  to 10.0.15.255 with 4094 IPs.
  # 10.0.16.0/20 from 10.0.16.0 to 10.0.31.255 with 4094 IPs.
  # 10.0.32.0/20 from 10.0.32.0 to 10.0.47.255 with 4094 IPs.
  # 10.0.48.0/20 from 10.0.48.0 to 10.0.63.255 with 4094 IPs.
  # ref. https://developer.hashicorp.com/terraform/language/functions/cidrsubnet
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]
  #
  # e.g., a reigon of 4 AZs will have:
  # 10.0.60.0/24 from 10.0.60.0 to 10.0.60.255 with 254 IPs.
  # 10.0.61.0/24 from 10.0.61.0 to 10.0.61.255 with 254 IPs.
  # 10.0.62.0/24 from 10.0.62.0 to 10.0.62.255 with 254 IPs.
  # 10.0.63.0/24 from 10.0.63.0 to 10.0.63.255 with 254 IPs.
  #
  # NOTE: use 60 to avoid CIDR range conflicts in a region of >=4 AZs.
  # ref. https://developer.hashicorp.com/terraform/language/functions/cidrsubnet
  public_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 60)]

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = 1
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"             = 1
  }
}

resource "aws_security_group" "eks" {
  name_prefix = local.cluster_name
  description = "EKS cluster security group."
  vpc_id      = module.vpc.vpc_id

  tags = {
    "Name" = "${local.cluster_name}-eks_cluster_sg"
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.5.1"

  cluster_name    = local.cluster_name
  cluster_version = "1.26"

  vpc_id                         = module.vpc.vpc_id
  subnet_ids                     = module.vpc.private_subnets
  cluster_endpoint_public_access = true

  eks_managed_node_group_defaults = {
    ami_type = "AL2_x86_64"
  }

  # https://docs.aws.amazon.com/eks/latest/APIReference/API_Nodegroup.html
  eks_managed_node_groups = {
    one = {
      use_custom_launch_template = false
      name                       = "t3xlarge"

      instance_types = ["t3.xlarge"]
      disk_size      = 100

      min_size     = 1
      max_size     = 2
      desired_size = 2
    }
  }

  create_cluster_security_group = false
  cluster_security_group_id     = aws_security_group.eks.id

  manage_aws_auth_configmap = true

  aws_auth_users = [
    for user in data.aws_iam_group.dev_members.users : {
      userarn  = "${user.arn}"
      username = "${user.user_name}"
      groups   = ["system:masters"]
    }
  ]
}

locals {
  oidc_id    = substr(module.eks.cluster_oidc_issuer_url, length(module.eks.cluster_oidc_issuer_url) - 32, 32)
  account_id = data.aws_caller_identity.current.account_id
}
