# https://docs.aws.amazon.com/eks/latest/APIReference/API_Nodegroup.html
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_node_group
# https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/modules/eks-managed-node-group/main.tf
resource "aws_eks_node_group" "al2_x86_64_cpu_t3xlarge" {
  cluster_name    = module.eks.cluster_name
  node_group_name = "${var.cluster_name}-al2-x86_64-t3xlarge"
  node_role_arn   = aws_iam_role.mng_iam_role.arn

  # no need to be in public subnets
  # when all services are exposed via ingress/LB
  subnet_ids = module.vpc.private_subnets

  # can only specify kubernetes version without an image id
  version = "1.26"

  scaling_config {
    min_size     = 1
    desired_size = 4
    max_size     = 10
  }
  update_config {
    max_unavailable = 1
  }

  # Force version update if existing pods are unable to be drained due to a pod disruption budget issue.
  force_update_version = false

  capacity_type  = "ON_DEMAND"
  instance_types = ["t3.xlarge"]

  # TODO: make this configurable
  disk_size = 100

  # https://docs.aws.amazon.com/eks/latest/APIReference/API_Nodegroup.html
  ami_type = "AL2_x86_64"

  # set same as "eks-managed-node-group" module
  lifecycle {
    create_before_destroy = true

    # Optional: Allow external changes without Terraform plan difference
    # useful when scaling up/down with cluster-autoscaler
    ignore_changes = [
      scaling_config[0].desired_size,
    ]
  }

  # NOTE: no need for "launch_template"
  # EKS managed node group provides one for AL2 nodes

  depends_on = [
    module.vpc,
    aws_security_group.eks,
    aws_security_group.nodes,
    aws_iam_role_policy_attachment.role_policy_attachment_AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.role_policy_attachment_AmazonEC2ContainerRegistryReadOnly,
    aws_iam_role_policy_attachment.role_policy_attachment_AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.efs_iam_role_policy_attachment_node,
    module.eks
  ]
}

resource "aws_autoscaling_group_tag" "al2_x86_64_cpu_t3xlarge" {
  autoscaling_group_name = aws_eks_node_group.al2_x86_64_cpu_t3xlarge.resources[0].autoscaling_groups[0].name

  # add extra label in case we run cluster-autoscaler in parallel with others
  # e.g., karpenter
  tag {
    key                 = "autoscaler-kind"
    value               = "cluster-autoscaler"
    propagate_at_launch = true
  }
}
