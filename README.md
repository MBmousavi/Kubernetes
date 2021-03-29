# Stage 1: Application
The sample Go application should be dockerized. We do this by creating a simple multi stage Dockerfile. The Dockerfile is in docker directory in our repo.
Next step is to build our image with proper tag and push it into docker hub.
The dockerfile and the main.go files are in the same directory.

`docker build . -t smbmousavi/testgoapp:<tag-name>`

You have to be authenticated to push to docker hub. 

`docker push smbmousavi/testgoapp:<tag-name>`

We use 2 methods for application health check:

1. Kubernetes `livenessProbe`:

It will be configured in `deployment.yaml` file.
In the `deployment.yaml` file, you can see that the Pod has a single container. The `periodSeconds` field specifies that the kubelet should perform a liveness probe every 5 seconds. The `initialDelaySeconds` field tells the kubelet that it should wait 5 seconds before performing the first probe. To perform a probe, the kubelet sends an HTTP GET request to the server that is running in the container and listening on port 8020. If the handler for the server's `/` path returns a success code, the kubelet considers the container to be alive and healthy. If the handler returns a failure code, the kubelet kills the container and restarts it.

2. Application health check endpoint:

In main.go we defined a `/health` endpoint. If it returns 200 http status code, application is up and running. 

# Stage 2: Kubernetes Cluster 
We setup our kubernetes with minikube. Minikube is an application that provide a single node kubernetes on your local system.
It can be installed on MacOS, Ubuntu and windows operating system. Our operating system is windows 10. Be aware that Minikube will install only kubernetes and not `kubectl`, we should install it manually. This fact is also true about `Helm`.

After installing Minikube installer from this guide [link](https://minikube.sigs.k8s.io/docs/start/), run this command to start minikube:

`minikube start`

Remember that we should install `kubectl`. It can be done with this guide [link](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

After installing run these 2 commands to check cluster status and ip:

`kubectl cluster-info`

`minikubr ip`

![1](https://user-images.githubusercontent.com/38520491/109940286-22ea4480-7ce7-11eb-9523-752246de6393.PNG)

Next step is to setup `Helm`. Helm is "The Kubernetes package manager", and just like `kubectl` it should be install on our host operating system.
We use this [link](https://github.com/helm/helm/releases) to install `Helm` on windows 10.

![2](https://user-images.githubusercontent.com/38520491/109943302-49f64580-7cea-11eb-860f-8ad8aa94e767.PNG)

Next Step in to setup `Helm Chart`. A Helm Chart is a collection of templates and settings that describe a set of Kubernetes resources.

we do this by run a single command. `helm create k8s`. The k8s is the name of my desired directory name. It can be anything you want.

Helm created a bunch of files for you in that directory that are usually important for setup an application in Kubernetes. We can remove a lot of the created files. Let’s go through the only required files.

We only need these files and other files can be deleted:

`k8s/Chart.yaml`

`k8s/templates/deployment.yaml`

`k8s/templates/service.yaml`

`k8s/values.yaml`

All these files are added into our repo. It looks very similar to plain Kubernetes file. For example you can see different placeholders for the replica count environment variable. The `replicas: {{ .Values.replicaCount }}` will be filled with `replicaCount: 3` from values.yaml file.
This feature gives us ability to dynamically change our application attributes, and don't need to hardcode them.

Let's move on forward. we can execute a "dry-run" command to see our final yaml files. In this example we have only one variable in our values.yaml file but in more complicated situations it's important to check what will be our final yaml files that are going to feed into kubernetes.

Run this command to see final yam file: 

`helm template --debug k8s`

If you are happy with the output run this command to initial application.

`helm install mygoapp k8s`. mygoapp  is a custom name that we put on this app and k8s is the name of that Helm chart directory.

Now you can run `helm list` to see installed application. With `kubectl` commands like `kubectl get pods`, `kubectl get svc` and `kubectl get deploy` you can see kubernetes resources that just installed.

You can check the minikube ip with nodeport that we just exposed to see final result.

![3](https://user-images.githubusercontent.com/38520491/109948890-fe469a80-7cef-11eb-8597-695dcfb2d073.PNG)

Next step is to install `prometheus` to monitor kubernetes and our application.

We install `prometheus` with `Helm.` First you need to add `prometheus repo` to your `Helm`.

`helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`

Then you can then run `helm search repo prometheus-community` to see the charts.

![4](https://user-images.githubusercontent.com/38520491/110074596-b6c71980-7d96-11eb-8f13-0d9feda4908f.PNG)

Then we proceed to install `prometheus-operator`. The differences between `prometheus-operator` and `prometheus` in this chart are that `prometheus-operator` comes with `Grafana` too and it's deprecated, but `prometheus` comes with only `prometheus`. In this example we use `prometheus-operator`.
Another tiny thing that we could create a new `namespace` in kubernetes called monitoring and grant some privileges to it and create `prometheus-operator` in that space, but we move on with our `default namespace` in kubernetes.

Run this command to install `prometheus-operator`:

`helm install prometheus-operator prometheus-community/prometheus-operator`

To check the result run `helm list` to see the applications that `helm` installed.

![6](https://user-images.githubusercontent.com/38520491/110076707-2db1e180-7d9a-11eb-8174-657abccfec0f.PNG)

Now we need to access our `Grafana` dashboard. First we need to get pod name of `Grafana` to forward it's port to our host operating system.

![8](https://user-images.githubusercontent.com/38520491/110079415-56d47100-7d9e-11eb-8f60-3aadd722a458.PNG)

The with `kubectl port-forward prometheus-operator-grafana-XXXXXXXXXX-XXXXX 3000:3000` we do this. Now `Grafana` dashboard is accessible in `http://127.0.0.1:3000/`

![7](https://user-images.githubusercontent.com/38520491/110079706-afa40980-7d9e-11eb-9a40-e487a8c4ed27.PNG)

Take your time to explore the dashboard. We can see `prometheus` is monitoring pods, nodes, resource usage and so much other stuff.

![9](https://user-images.githubusercontent.com/38520491/110080369-9780ba00-7d9f-11eb-9f41-d472b96a810e.PNG)

***

* Application Health Check with `prometheus BlackBox Exporter`

We use another helm repo (stable - https://charts.helm.sh/stable) for setup `prometheus BlackBox Exporter` and `Grafana`. It's a deprecated repo but it's fine for our test project.

First step is to add helm repo:

`helm repo add stable https://charts.helm.sh/stable`

it's good to run a `helm repo update`. Now we need 3 applications to install with helm:

`helm install prometheus stable/prometheus`

`helm install prometheus-blackbox-exporter stable/prometheus-blackbox-exporter`

`helm install grafana stable/grafana`

Now we have to tell Prometheus how to get the blackbox exporter. Here we enter the HTTP endpoint that is to be queried by the Blackbox Exporter. We write the additional configuration in a file and deploy it via Helm Upgrade: `prom-blackbox-scrape.yaml` (It's in the code repo)

`helm upgrade --reuse-values -f prom-blackbox.yaml prometheus stable/prometheus`

Now we port forward Prometheus, Grafana and Blackbox Exporter panel to see what just happened.

`kubectl get pods` then forward ports for every desired pods.

`kubectl port-forward prometheus-server-xxxxxxxxx-xxxxx 9090:9090`

`kubectl port-forward prometheus-blackbox-exporter-xxxxxxxxx-xxxxx 9115:9115`

`kubectl port-forward grafana-xxxxxxxxx-xxxxx 3000:3000`

Prometheus:

![11](https://user-images.githubusercontent.com/38520491/110339976-fe230380-803d-11eb-8d7e-fb10ec5d6c56.PNG)


Blackbox Exporter:

![12](https://user-images.githubusercontent.com/38520491/110340134-2874c100-803e-11eb-9070-880d98213df1.PNG)

The metrics that are storing into Prometheus:

![13](https://user-images.githubusercontent.com/38520491/110340246-4b9f7080-803e-11eb-97f6-04896cdd4cb6.PNG)

Now we can login into Grafna dashboard and create some monitoring panel for our application, It's just a example of what we can do:

![14](https://user-images.githubusercontent.com/38520491/110340852-f3b53980-803e-11eb-803f-9db9fed779bd.PNG)

***

* Question: Describe a solution for zero-downtime deployment with Kubernetes.

In General kubernets has 2 types of strategy to update an application. `Recreate` and `Rolling Update`. `Rolling Update` is the default.

`Rolling Update` means a `new deployment` start scaling up or creating and the `old deployment` start scaling down or terminating.

`Rolling Update` is tied with 2 important concepts (or variables) in kubernetes: `Max Surge` and `Max Unavailable`.

Deployment ensures that only a certain number of Pods are down while they are being updated. By default, it ensures that at least 75% of the desired number of Pods are up (25% max unavailable).

Deployment also ensures that only a certain number of Pods are created above the desired number of Pods. By default, it ensures that at most 125% of the desired number of Pods are up (25% max surge).

Deployment first created a new Pod, then deleted some old Pods, and created new ones. It does not kill old Pods until a sufficient number of new Pods have come up, and does not create new Pods until a sufficient number of old Pods have been killed. So if something goes wrong for example wrong image and we have 10 pods, only 2 of them become unavailable not all of them. with `Rolling Update` feature of kubernetes we can have minimum or maybe zero downtime for updating our applications.

Another strategy is `Recreate`. It terminates all pods in the current deployment and start to create new pods for new deployment. 
