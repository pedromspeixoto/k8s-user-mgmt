terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.2"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}

provider "kubernetes" {
  config_path            = var.kubectl_config_path
  config_context_cluster = var.kind_cluster_name
}

provider "helm" {
  kubernetes {
    config_path            = var.kubectl_config_path
    config_context_cluster = var.kind_cluster_name
  }
}

provider "kubectl" {
  config_path            = var.kubectl_config_path
  config_context_cluster = var.kind_cluster_name
  load_config_file       = true
}