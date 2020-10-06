# README

This is a Keptn SLI Provider built for my talk at Neotys (Jurassic) PAC 2020: https://www.neotys.com/performance-advisory-council/andreas-grabner

It was built using the [keptn-service-template-go](https://github.com/keptn-sandbox/keptn-service-template-go/generate) repository instructions which I kept here for reference!


# pac-sliprovider
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/pac-sliprovider)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/pac-sliprovider)](https://goreportcard.com/report/github.com/keptn-sandbox/pac-sliprovider)

This implements a pac-sliprovider for Keptn. If you want to learn more about Keptn visit us on [keptn.sh](https://keptn.sh)

## Compatibility Matrix

| Keptn Version    | [pac-sliprovider Docker Image](https://hub.docker.com/r/grabnerandi/pac-sliprovider/tags) |
|:----------------:|:----------------------------------------:|
|       0.7.x      | grabnerandi/pac-sliprovider:0.1.0 |

## Full Installation Instructions

As I have presented this for Neotys PAC event I want to give you detailed instructions on how you can replicate what I have done in my talk.
In my talk I took a t2.medium Amazon Linux 2 EC2 machine where I
* Installed Keptn based on [Keptn 0.7.1 on K3s](https://github.com/keptn-sandbox/keptn-on-k3s/tree/release-0.7.1)
* Installed my *pac-sliprovider*
* Created a Keptn *pac-project* 
* Configured the *pac-sliprovider* as SLI provider for that project
* Created a Keptn Service *pacservice*
* Uploaded an SLI.yaml and SLO.yaml
* Executed a couple of Keptn Quality Gates

Now - lets go into the details of each step so you can replicate this!

### Step 1 - Install Keptn

As I said - I just go with the simplest option which is Keptn on k3s. At the time of the conference Keptn 0.7.1 was the latest Keptn version so I decided to use that [0.7.1 release](https://github.com/keptn-sandbox/keptn-on-k3s/tree/release-0.7.1) on the [Keptn on K3s](https://github.com/keptn-sandbox/keptn-on-k3s) github repo. If there are newer versions available make sure to pick the latest!

In my case I launched an Amazon Linux 2 EC2 size t2.medium. Keptn on k3s only needs 1vcpu and 4GB of RAM and has been tested on a variety of platforms. Check out the [prerequisits](https://github.com/keptn-sandbox/keptn-on-k3s#prerequisites) on the Keptn on k3s github repo!

To install keptn on k3s on an AWS EC2 I just executed the following command:
```console
sudo curl -Lsf https://raw.githubusercontent.com/keptn-sandbox/keptn-on-k3s/0.7.1/install-keptn-on-k3s.sh | bash -s - --provider=aws
```
The output of that command after its finished looks something like this
```console
#######################################>
# Deployment Summary
#######################################>
API URL   :      https://172.31.x.y/api
Bridge URL:      https://172.31.x.y/bridge
Bridge Username: keptn
Bridge Password: PASSWORDFORBRIDGE
API Token :      KEPTNAPITOKEN
To use keptn:
- Install the keptn CLI: curl -sL https://get.keptn.sh | sudo -E bash
- Authenticate: keptn auth  --api-token "KEPTNAPITOKEN" --endpoint "https://172.31.x.y/api"
```

To finish the installation just follow the two additional instructions to install the Keptn CLI and then authenticate it!
You should see an output similar to this:
```console
$ keptn auth  --api-token xxxxxxx
keptn creates the folder /home/ec2-user/.keptn/ to store logs and possibly creds.
Starting to authenticate
Successfully authenticated
Using a file-based storage for the key because the password-store seems to be not set up.
```

Now we are ready to use Keptn through the CLI.

### Step 2 - Install my PAC SLI Provider

Keptn is an event-driven control plane which means it issues events to trigger different activities, e.g: deploy, test, get sli data, validate, ...
A Keptn Service - such as my PAC SLI Provider - needs to be installed on the Keptn k8s cluster and needs to subscribe to the events that the servie wants to handle.
As of Keptn 0.7.x we do this by simply applying the deployment yaml which will deploy my pac-sliprovider as a pod.

The easiest way to get the deployment yaml - and some other files we will need later - on your machine is a simple git pull:
```console

```


The *pac-sliprovider* can be installed as a part of [Keptn's uniform](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the *pac-sliprovider* in your Keptn Kubernetes cluster, apply the [`deploy/service.yaml`](deploy/service.yaml) file:

```console
kubectl apply -f deploy/service.yaml
```

This should install the `pac-sliprovider` together with a Keptn `distributor` into the `keptn` namespace, which you can verify using

```console
kubectl -n keptn get deployment pac-sliprovider -o wide
kubectl -n keptn get pods -l run=pac-sliprovider
```

### Up- or Downgrading

Adapt and use the following command in case you want to up- or downgrade your installed version (specified by the `$VERSION` placeholder):

```console
kubectl -n keptn set image deployment/pac-sliprovider pac-sliprovider=grabnerandi/pac-sliprovider:$VERSION --record
```

### Uninstall

To delete a deployed *pac-sliprovider*, use the file `deploy/*.yaml` files from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

## Development

Development can be conducted using any GoLang compatible IDE/editor (e.g., Jetbrains GoLand, VSCode with Go plugins).

It is recommended to make use of branches as follows:

* `master` contains the latest potentially unstable version
* `release-*` contains a stable version of the service (e.g., `release-0.1.0` contains version 0.1.0)
* create a new branch for any changes that you are working on, e.g., `feature/my-cool-stuff` or `bug/overflow`
* once ready, create a pull request from that branch back to the `master` branch

When writing code, it is recommended to follow the coding style suggested by the [Golang community](https://github.com/golang/go/wiki/CodeReviewComments).

### Where to start

If you don't care about the details, your first entrypoint is [eventhandlers.go](eventhandlers.go). Within this file 
 you can add implementation for pre-defined Keptn Cloud events.
 
To better understand Keptn CloudEvents, please look at the [Keptn Spec](https://github.com/keptn/spec).
 
If you want to get more insights, please look into [main.go](main.go), [deploy/service.yaml](deploy/service.yaml),
 consult the [Keptn docs](https://keptn.sh/docs/) as well as existing [Keptn Core](https://github.com/keptn/keptn) and
 [Keptn Contrib](https://github.com/keptn-contrib/) services.

### Common tasks

* Build the binary: `go build -ldflags '-linkmode=external' -v -o pac-sliprovider`
* Run tests: `go test -race -v ./...`
* Build the docker image: `docker build . -t grabnerandi/pac-sliprovider:dev` (Note: Ensure that you use the correct DockerHub account/organization)
* Run the docker image locally: `docker run --rm -it -p 8080:8080 grabnerandi/pac-sliprovider:dev`
* Push the docker image to DockerHub: `docker push grabnerandi/pac-sliprovider:dev` (Note: Ensure that you use the correct DockerHub account/organization)
* Deploy the service using `kubectl`: `kubectl apply -f deploy/`
* Delete/undeploy the service using `kubectl`: `kubectl delete -f deploy/`
* Watch the deployment using `kubectl`: `kubectl -n keptn get deployment pac-sliprovider -o wide`
* Get logs using `kubectl`: `kubectl -n keptn logs deployment/pac-sliprovider -f`
* Watch the deployed pods using `kubectl`: `kubectl -n keptn get pods -l run=pac-sliprovider`
* Deploy the service using [Skaffold](https://skaffold.dev/): `skaffold run --default-repo=your-docker-registry --tail` (Note: Replace `your-docker-registry` with your DockerHub username; also make sure to adapt the image name in [skaffold.yaml](skaffold.yaml))


### Testing Cloud Events

We have dummy cloud-events in the form of [RFC 2616](https://ietf.org/rfc/rfc2616.txt) requests in the [test-events/](test-events/) directory. These can be easily executed using third party plugins such as the [Huachao Mao REST Client in VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

## Automation

### GitHub Actions: Automated Pull Request Review

This repo uses [reviewdog](https://github.com/reviewdog/reviewdog) for automated reviews of Pull Requests. 

You can find the details in [.github/workflows/reviewdog.yml](.github/workflows/reviewdog.yml).

### GitHub Actions: Unit Tests

This repo has automated unit tests for pull requests. 

You can find the details in [.github/workflows/tests.yml](.github/workflows/tests.yml).

### Travis-CI: Build Docker Images

This repo uses [Travis-CI](https://travis-ci.org) to automatically build docker images. This process is optional and needs to be manually 
enabled by signing in into [travis-ci.org](https://travis-ci.org) using GitHub and enabling Travis for your repository.

After enabling Travis-CI, the following settings need to be added as secrets to your repository on the Travis-CI Repository Settings page:

* `REGISTRY_USER` - your DockerHub username
* `REGISTRY_PASSWORD` - a DockerHub [access token](https://hub.docker.com/settings/security) (alternatively, your DockerHub password)

Furthermore, the variable `IMAGE` needs to be configured properly in the respective section:
```yaml
env:
  global:
    - IMAGE=grabnerandi/pac-sliprovider # PLEASE CHANGE THE IMAGE NAME!!!
```
You can find the implementation of the build-job in [.travis.yml](.travis.yml).

## How to release a new version of this service

It is assumed that the current development takes place in the master branch (either via Pull Requests or directly).

To make use of the built-in automation using Travis CI for releasing a new version of this service, you should

* branch away from master to a branch called `release-x.y.z` (where `x.y.z` is your version),
* write release notes in the [releasenotes/](releasenotes/) folder,
* check the output of Travis CI builds for the release branch, 
* verify that your image was built and pushed to DockerHub with the right tags,
* update the image tags in [deploy/service.yaml], and
* test your service against a working Keptn installation.

If any problems occur, fix them in the release branch and test them again.

Once you have confirmed that everything works and your version is ready to go, you should

* create a new release on the release branch using the [GitHub releases page](https://github.com/keptn-sandbox/pac-sliprovider/releases), and
* merge any changes from the release branch back to the master branch.

## License

Please find more information in the [LICENSE](LICENSE) file.
