# ArgoCD Services Namespace
resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

# Prod Services Namespace
resource "kubernetes_namespace" "prod" {
  metadata {
    name = "prod"
  }
}

# ArgoCD Helm Release
resource "helm_release" "argocd_helm_release" {
  depends_on = [
    kubernetes_namespace.argocd,
    kubernetes_namespace.prod
  ]

  name       = "argocd"
  namespace  = kubernetes_namespace.argocd.metadata[0].name
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "5.38.0"
}

# ArgoCD App of Apps for Prod Services
resource "helm_release" "argocd_prod_apps" {
  depends_on = [
    helm_release.argocd_helm_release
  ]

  name       = "argocd-prod-apps"
  namespace  = kubernetes_namespace.argocd.metadata[0].name
  chart      = "argocd/apps/"
}

# NGINX Ingress Controller for Kind
data "kubectl_path_documents" "nginx_ingress_manifests" {
  pattern  = "${path.module}/nginx/config.yaml"
}

resource "kubectl_manifest" "nginx_ingress" {
  count     = length(data.kubectl_path_documents.nginx_ingress_manifests.documents)
  yaml_body = element(data.kubectl_path_documents.nginx_ingress_manifests.documents, count.index)
}