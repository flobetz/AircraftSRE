# AircraftSRE
simple golang flights webservice, running in Docker on k8s on Azure, built and deployed with Jenkins.  
All packed in this project.
  
Grafana can be found here:  
http://flightoperator-grafana.eastus.cloudapp.azure.com:3000  

flights API can be found here:  
http://flightoperator.eastus.cloudapp.azure.com/v1/flights

URL to Github project:  
https://github.com/flobetz/AircraftSRE

### Contents of this project
- local Jenkins Dockerfile, build and startup script with docker socket and var/jenkins_home dir mounted
- PostgreSQL Dockerfile and a startup script
- flights.go a simple go webapp which serves a REST API for flight management and a endpoint for metrics
- Dockerfile for the flights app
- Prometheus Dockerfile and config to scrape metrics from the flights app
- Grafana Dockerfile, config, dashboard and data sources to visualize the metrics prometheus scrapes
- Azure AKS infrastructure written in terraform deployed with terraform and Jenkins
- Kubernetes Deployment file for the whole stack (App, DB, prometheus, grafana)
- Jenkins Pipeline scripts to build and upload all docker images into private Azure registry
- Jenkins Pipeline script to deploy the Azure infrastructure using terraform
- Jenkins Pipeline script which automatically updates the whole stack when changes happen on the master branch
- docker-compose file for the whole stack which can be used for local testing
- some testing scripts which call the flights endpoints

### Project structure
- APP: webservice golang application and a Dockerfile which compiles the webservice and stores it in a new Docker image
- DB: simple Dockerfile of a PostgreSQL Database to track the version of the used Docker base image
- INFRA: Terraform code which describes the underlying Infrastructure of the webservice and it's components
(This Terraform code creates a Kubernetes Cluster on three Linux VMs on Azure AKS)
- Jenkins: Jenkins is only used locally within this project. "buildJenkins.sh" creates a Jenkins Docker image with preinstalled plugins from "plugins.txt".
"startJenkins.sh" starts a local Jenkins instance (and mounts the Jenkins Home directory and the clients docker socket to keep data persistent and to be able to use docker within Jenkins)  
Additional: Add Jenkins configuration as code to the repo - This makes it possible to persist the same Jenkins configuration across a development team  
Additional: Set Up a team wide Vault instance where Jenkins pipelines can get secrets from - This makes Jenkins even more moveable, 
everything is tracked as code and secrets are managed globally
- kubernetes: Kubernetes deployment file which consists of deployments and services for the webservice, the database, Prometheus and Grafana
- local: Docker compose file to start the application stack locally (for testing and debugging purpose). "testall.sh" calls all endpoints our service provides.
- monitoring:  
    - grafana: Dockerfile for creating a custom Grafana image which includes a prometheus datasource, some plugins and a dashboard for visualizing the metrics of our webservice
    - prometheus: Dockerfile for creating a custom Prometheus image with configuration to our webservices metrics endpoint
- pipeline:
    - AutoDeploy.groovy: Jenkinsfile which is getting triggered as soon as a new commit gets created on the master branch. 
    It checks in which directories changes have been made and recreates only those services. Last step is to create or update the project on Azure AKS.
    - BuildApp, BuildDB, BuildGrafana, BuildPrometheus: Jenkinsfile which builds the according service and uploads the resulting Docker image to the private Azure Docker registry.
    - DeployInfrastructure.groovy: Applies the given Terraform code to Azure.
    - ListImages.groovy: Jenkinsfile which lists all Docker images of the private Azure Docker registry.
- testing: Bash script which calls all endpoints of the webservice which is running in Azure 