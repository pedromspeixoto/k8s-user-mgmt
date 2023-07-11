variable "kind_cluster_name" {
  description = "Name for the kind cluster"
  type        = string
  default     = "kind-user-mgmt-cluster"
}

variable "kubectl_config_path" {
  description = "Path to kubectl config file"
  type        = string
  default     = "~/.kube/config"
}