# ARGOCD

## Running argocd

Here are a few basic commands to get started with ARGOCD. Full documentation available at [https://argoproj.github.io/](https://argoproj.github.io/).

```bash
az aks list -o table

az aks get-credentials -n <aks-cluster-name> -g <resource-group-name>


kubectl  # (if not installed then follow these steps)

az aks install-cli

sudo rm /usr/local/kubectl

mkdir /usr/local/kubectl

alias k=kubectl

# (continue)


# Review for argocd resources running in your subscription

k get pods --all-namespaces

k get svc --all-namespaces


argocd # (if not installed follow these steps)

VERSION=$(curl --silent "https://api.github.com/repos/argoproj/argo-cd/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

sudo curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64

sudo chmod +x /usr/local/bin/argocd

# (continue)


# Get ARGOCD instance password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

argocd account update-password
# >> password
# >> new password
# >> confirm password

# In a new terminal window, run loop for handling localhost connection for 8080
k port-forward svc/argo-argocd-server -n argocd 8080:443

  # Web browser navigate to localhost:8080 > login

# Install sample app into ARGOCD
k apply -f /workspaces/symphony/caf/app_argocd/level4/argocd/traefik/traefik.yaml
```

