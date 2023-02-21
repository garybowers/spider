<img src="assets/logo.avif"  width="300" height="300">

# Spiderweb

Scale your team quickly by moving your IDE to the cloud using Kubernetes!


## What is it?

Spiderweb is a authenticated connection broker for your web based development environments.
Scale up/down your developer machines as your need them using the power of Kubernetes.

## How does it work?

Spiderweb is typically deployed to the same cluster your IDE's would be leveraging and makes use of in-cluster RBAC to manage the deployments.

Users enter the domain pointing to spiderweb, if the session has a valid session cookie then spiderweb
looks for a valid deployment based on the user's e-mail address returned by the autheticated session.
If there's no valid cookie, spiderweb redirects to the 3rd party authentication system (e.g. Google Oauth).

If no valid deployment in kubernetes exists, one is created and the user is presented with a 'please wait' screen.

Once a valid deployment exists and the user is authenticated then spiderweb opens a reverse proxy to that deployment/pod.

When the user logs off by going to the url endpoint https://spider.mydomain.tld/logout the deployment is deleted and your cluster's autoscaling should reduce the number of nodes needed, saving money.

## Spider 

Spider is the accompanying image in docker format that your end users will use, it contains the IDE itself and a host of other tools.
Theia [https://theia-ide.org/](https://theia-ide.org/) is chosen as the default image to deploy as it's a flexible and powerful open platform that has compatibility with VSCode Extensions.
You can create your own image as long as it serves a web endpoint.

As this has primarily been built around Google Kubernetes Engine, but will work on any Kubernetes distribution on any cloud, Artifact registry is preffered as it supports image straming for faster boootstrap times.

### Spider - Storage and Persistance

When users logout their instance of the deployment is deleted we use shared storage (NFS) mounted to the container to persist sessions and files, when a user logs back in their session is restored as-is.

## Current limitations and roadmap

Right now the only authentication system spiderweb supports is Google OAuth, feel free to raise a issue or contribute to add other authentication backends e.g. Okta or Azure AD

## Deploying 

For a example deployment of spiderweb see the deployment file here [deploy/spiderweb-deployment.yaml](deploy/spiderweb-deployment.yaml)

For information on creating a Google OAuth application see [https://support.google.com/cloud/answer/6158849?hl=en](https://support.google.com/cloud/answer/6158849?hl=en)


Some environment variables that spiderweb needs:

GOOGLE_CLIENT_ID - 

GOOGLE_CLIENT_SECRET - 

OAUTH_CALLBACK_URL

SPIDER_IMAGE

SPIDER_NAMESPACE

SPIDER_APPNAME

SPIDERWEB_LISTEN_PORT

SPIDER_NFS_SERVER

SPIDER_FQDN


URI: https://spider.domain.com 
CallBack URI: https://spider.ssp.immersion.dev/auth/google/callback
	      http://spider.ssp.immersion.dev/auth/google/callback
