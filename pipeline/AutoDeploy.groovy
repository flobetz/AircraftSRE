#!groovy​

node("master"){

    // set scm trigger to scan for changes on repo every minute
    properties([pipelineTriggers([pollSCM('* * * * *')])])

    echo "branchname: ${env.BRANCH_NAME}"

    // check branchname
    if (env.BRANCH_NAME.equals("master")) {

        // checkout repo
        stage("Checkout") {
            checkout scm
        }

        // Check in which subdirectories there have been changes within the last commit
        stage("Check where changes have been made") {
            lastCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
            filesChanged = sh(returnStdout: true, script: "git diff-tree --no-commit-id --name-only -r ${lastCommit}").trim()
            echo "lastcommit: ${lastCommit}"
            echo "fileschanged:"
            echo "${filesChanged}"
            lines = filesChanged.split('\n')
        }

        // trigger other jobs - depending on where changes have been made
        lines.each {
            echo "file changed: ${it}"
            def dirs = it.split('/')

            switch(dirs[0])
            {
                case 'APP':
                    stage("Build app Image") {
                        echo "Building APP docker image"
                        build 'docker-build/BuildApp'
                    }
                    break
                case 'DB' :
                    stage("Build DB Image") {
                        echo "Building PostgreSQL DB docker image"
                        build 'docker-build/BuildDB'
                    }
                    break
                case 'INFRA' :
                    stage("Update Infrastructure in Azure") {
                        echo "Update Azure Infrastructure"
                        build 'DeployAzureInfrastructure'
                    }
                    break
                case 'monitoring' :
                    if (dirs[1].equals("grafana")) {
                        stage("Build Grafana Image") {
                            echo "Building Grafana docker image"
                            build 'docker-build/BuildGrafana'
                        }
                    } else if (dirs[1].equals("prometheus")) {
                        stage("Build Prometheus Image") {
                            echo "Building Prometheus docker image"
                            build 'docker-build/BuildPrometheus'
                        }
                    }
                    break
                default:
                    break
            }
        }

        // (re)deploy everythin to AKS
        stage("Deploy to AKS") {
            dir("kubernetes") {
                // set KUBECONFIG
                withCredentials([file(credentialsId: 'kubeconf-azure-aks', variable: 'kubeconf')]) {
                    def kubectl = "kubectl  --kubeconfig=\$kubeconf "
                    // check if services are running
                    svc = sh(returnStdout: true, script: "${kubectl} get svc | tail -n +2 | wc -l").trim()
                    if (svc == 0) {
                        // no service is running, recreate everything
                        sh "${kubectl} create -f ./flightOperatorAzure.yaml"
                    } else {
                        // at least one service is running, just apply updates
                        sh "${kubectl} apply -f ./flightOperatorAzure.yaml"
                    }
                }
            }
        }
    }
}