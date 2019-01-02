#!groovy​

node("master"){

    stage("Checkout") {
        checkout scm
    }

    dir("APP") {

        stage("Build") {
            sh "docker build -t flightoperatorreg.azurecr.io/flightoperator:latest ."
        }

        stage("Push") {
            withCredentials([string(credentialsId: 'azuretenant', variable: 'azuretenant'), usernamePassword(credentialsId: 'azureServicePrincipal', passwordVariable: 'password', usernameVariable: 'username')]) {
                sh """
                    echo "Login to Azure registry flightoperatorreg.azurecr.io"
                    az login --service-principal -u ${username} -p ${password} --tenant ${azuretenant}
                
                    az acr login --name flightOperatorReg
                    docker push flightoperatorreg.azurecr.io/flightoperator:latest
                """
            }
        }
    }
}