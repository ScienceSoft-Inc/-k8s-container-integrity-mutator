![GitHub contributors](https://img.shields.io/github/contributors/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub last commit](https://img.shields.io/github/last-commit/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub issues](https://img.shields.io/github/issues/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub forks](https://img.shields.io/github/forks/ScienceSoft-Inc/k8s-container-integrity-mutator)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)

# k8s-container-integrity-mutator

This application provides the injection of any patch inside any k8s schemas like sidecar.

When applying a new scheme to a cluster, the application monitors the presence of a "
hasher-certificates-injector-sidecar" label and, if available, makes a patch.

## Architecture

### Statechart diagram

![File location: docs/diagrams/mutatorStatechartDiagram.png](/docs/diagrams/mutatorStatechartDiagram.png?raw=true "Statechart diagram")

### Sequence diagram

![File location: docs/diagrams/mutatorSequenceDiagram.png](/docs/diagrams/mutatorSequenceDiagram.png?raw=true "Sequence diagram")

## :hammer: Installing components

### Running locally

The code only works running inside a pod in Kubernetes.
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using `minikube`.
Example <https://minikube.sigs.k8s.io/docs/start/>

### Install Helm

Before using helm charts you need to install helm on your local machine.  
You can find the necessary installation information at this link https://helm.sh/docs/intro/install/

### Configuration

To work properly, you first need to set the configuration files:

+ values in the file `helm-charts/integrity-injector/values.yaml`
+ values in the file `helm-charts/demo-app-to-inject/values.yaml`

In order to use make targets to lunch integrity-injector with/without demo app, following environment variables should be set:

- DB_PASSWORD
- DB_USER
- DB_NAME
- RELEASE_NAME_DB

Build docker image:

```
make docker
```

### Run helm-charts

Install helm chart with mutator app

```
make helm-mutator
```

or via helm

```
helm install mutator helm-charts/integrity-injector
```

Install helm chart with demo app
install with db

```
make helm-demo-with-db
```

install without db

```
make helm-demo
```

Install demo with db and syslog server

```
make helm-demo-full
```

or through helm

```
helm install demo-app helm-charts/demo-app-to-inject
```
