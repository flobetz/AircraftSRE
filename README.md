# AircraftSRE
simple golang flights webservice, running in Docker on k8s on Azure, built and deployed with Jenkins. All packed in this repo

### What this project consists of
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

ToDo:
- :white_check_mark: Create go application which serves a REST API
- :white_check_mark: Create flights endpoints
- :white_check_mark: put App in a docker container
- :white_check_mark: Create Jenkins Docker image with docker, kubectl, azure cli
- :white_check_mark: deploy app image on Azure registry with local Jenkins instance
- :white_check_mark: deploy db image on Azure registry with local Jenkins instance
- :white_check_mark: create K8S Cluster on AKS, done with terraform and Jenkins in git
- :white_check_mark: Run images (app and DB) on K8S AKS Cluster  
- :white_check_mark: implement basic auth for all endpoints
- :white_check_mark: create metrics endpoint
- :white_check_mark: configure pronmetheus
- :white_check_mark: configure grafana
- :white_check_mark: put prometheus and grafana into docker-compose
- :white_check_mark: check departure time of a new flight
- :white_check_mark: remove flight number from POST request
- :white_check_mark: automatically create flight number when creating a new flight 
- :white_check_mark: create Jenkins build jobs for prometheus and grafana container
- :white_check_mark: add prometheus and grafana instance into kubernetes cluster
- :white_check_mark: Share Grafana URL within Repo:  
  [http://flightoperator-grafana.eastus.cloudapp.azure.com](http://flightoperator-grafana.eastus.cloudapp.azure.com:3000)  
  [http://flightoperator.eastus.cloudapp.azure.com/v1/flights](http://flightoperator.eastus.cloudapp.azure.com/v1/flights)  
- :white_check_mark: enable auto deployments when merging to master / develop (Jenkins)
- [optional] enable plugins.txt for Jenkins (and save in repo)
- [optional] enable config as code for Jenkins (and sae in repo)
- [optional] add some users to grafana
