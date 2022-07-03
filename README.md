# Kube Sync
Kube Sync is a simple service that copies data from a target Pod, then commits and pushes to a target GIT repository.

[![CodeQL](https://github.com/steve-winter/kube-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/steve-winter/kube-sync/actions/workflows/codeql-analysis.yml)  [![Go](https://github.com/steve-winter/kube-sync/actions/workflows/go.yml/badge.svg)](https://github.com/steve-winter/kube-sync/actions/workflows/go.yml)

## Targeted Use Cases
The initial use case was focussed on [Home Assistant](https://artifacthub.io/packages/helm/k8s-at-home/home-assistant) and [Zibgee2Mqtt](https://artifacthub.io/packages/helm/k8s-at-home/zigbee2mqtt) which run on my K3s setup. I wanted a consistent mechanism to backup and version control changes.

## Getting Started

### Docker
```shell
docker run \
  -v <LOCAL_CONFIG>:/etc/kube-sync \
  -v <LOCAL_KUBE_CONTEXT>:/.kube \
  kube-sync
```

## Design

### High Level Sequence
```mermaid
sequenceDiagram
participant KubeSync as Kube Sync
participant Git
participant KubeAPI
participant Pod
Activate KubeSync
KubeSync->>+Git: Clone repo
Git-->>-KubeSync: Git repo
KubeSync->>+KubeAPI: Locate Pod
KubeAPI-->>-KubeSync: Pod name
KubeSync->>+Pod: Initiate tar stream
Pod-->>-KubeSync: Download tar stream
KubeSync->>KubeSync: Commit changes
KubeSync->>+Git: Push changes
Git-->>-KubeSync: Push success
Deactivate KubeSync
```


# Known Gaps/Issues
1. Only tested with AT token access to GitHub
2. Does not currently use secrets for tokens, should be a very simple migration
3. Isolate to specific role and user

## References
- Solution uses the Kube Copy operations in Kubectl - https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/cp/cp.go
- Leverages the "Go-Git" library for Git integration - https://github.com/go-git/go-git
