# AircraftSRE
simple golang flights webservice, running in Docker on k8s on Azure, built and deployed with Jenkins. All packed in this repo

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
- check departure time of a new flight
- remove flight number from POST request
- automatically create flight number when creating a new flight 
- create Jenkins build jobs for prometheus and grafana container
- add prometheus and grafana instance into kubernetes cluster
- Share Grafana URL within Repo
- enable auto deployments when merging to master / develop (Jenkins)
https://docs.microsoft.com/de-de/azure/container-service/kubernetes/
