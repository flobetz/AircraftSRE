#!groovy​
node("master"){

    registry    = "flightoperatorreg.azurecr.io"
    image       = "prometheus"
    tag         = "latest"

    stage("Checkout") {
        checkout scm
    }

    dir("monitoring") {
        dir("prometheus") {
            stage("Build") {
                sh "docker build -t ${registry}/${image}:${tag} ."
            }

            stage("Push") {
                withCredentials([string(credentialsId: 'azuretenant', variable: 'azuretenant'), usernamePassword(credentialsId: 'azureServicePrincipal', passwordVariable: 'password', usernameVariable: 'username')]) {
                    sh """
                        echo "Login to Azure registry ${registry}"
                        az login --service-principal -u ${username} -p ${password} --tenant ${azuretenant}
                        
                        az acr login --name flightOperatorReg
                        docker push ${registry}/${image}:${tag}
                    """
                }
            }
        }
    }
}