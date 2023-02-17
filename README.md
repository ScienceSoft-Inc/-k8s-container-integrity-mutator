![GitHub contributors](https://img.shields.io/github/contributors/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub last commit](https://img.shields.io/github/last-commit/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub issues](https://img.shields.io/github/issues/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub forks](https://img.shields.io/github/forks/ScienceSoft-Inc/k8s-container-integrity-mutator)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)

# k8s-container-integrity-mutator
This application provides the injection of any patch inside any k8s schemas like sidecar.

When applying a new scheme to a cluster, the application monitors the presence of a "hasher-certificates-injector-sidecar" label and, if available, makes a patch.

## Architecture
### Statechart diagram
![File location: docs/diagrams/mutatorStatechartDiagram.png](/docs/diagrams/mutatorStatechartDiagram.png?raw=true "Statechart diagram")
### Sequence diagram
![File location: docs/diagrams/mutatorSequenceDiagram.png](/docs/diagrams/mutatorSequenceDiagram.png?raw=true "Sequence diagram")
## :hammer: Installing components
### Generate certificates
Generate CA in /tmp :
```
cfssl gencert -initca ./certificates/tls/ca-csr.json | cfssljson -bare /tmp/ca
```

Generate private key and certificate for SSL connection:
```
cfssl gencert \
-ca=/tmp/ca.pem \
-ca-key=/tmp/ca-key.pem \
-config=./certificates/tls/ca-config.json \
-hostname="k8s-webhook-injector,k8s-webhook-injector.default.svc.cluster.local,k8s-webhook-injector.default.svc,localhost,127.0.0.1" \
-profile=default \
./certificates/tls/ca-csr.json | cfssljson -bare /tmp/k8s-webhook-injector
```

Move your SSL key and certificate to the ssl directory:
```
mv /tmp/k8s-webhook-injector.pem ./certificates/ssl/k8s-webhook-injector.pem
mv /tmp/k8s-webhook-injector-key.pem ./certificates/ssl/k8s-webhook-injector.key
```

Update configuration data in the helm-charts/mutator/values.yaml file with your key in the appropriate field `serverKey` and certificate in the appropriate field `serverCrt`:
```
cat ./certificates/ssl/k8s-webhook-injector.key | base64 | tr -d '\n'
cat ./certificates/ssl/k8s-webhook-injector.pem | base64 | tr -d '\n'
```

Update field `caBundle` value in the helm-charts/mutator/values.yaml file with your base64 encoded CA certificate:
```
cat /tmp/ca.pem | base64 | tr -d '\n'
```
### Demo-app
Here is a demo application in which a busybox container in `patch-json-command.json` is injected to a pod with an nginx container

Build docker images mutator:
```
eval $(minikube docker-env)
docker build -t mutator .
```
## Install Helm
Before using helm charts you need to install helm on your local machine.  
You can find the necessary installation information at this link https://helm.sh/docs/intro/install/

### Configuration
To work properly, you first need to sett the configuration files:
+ values in the file `helm-charts/mutator/values.yaml`
+ values in the file `helm-charts/demo-app-to-inject/values.yaml`

### Run helm-charts
Install helm chart with mutator app
```
helm install mutator helm-charts/mutator
```
Install helm chart with demo app
```
helm install demo-app helm-charts/demo-app-to-inject
```