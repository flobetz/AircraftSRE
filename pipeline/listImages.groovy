#!groovy​

node("master"){

    stage("List Images") {
        withCredentials([string(credentialsId: 'azuretenant', variable: 'azuretenant'), usernamePassword(credentialsId: 'azureServicePrincipal', passwordVariable: 'password', usernameVariable: 'username')]) {
            sh """
               echo "Login to Azure registry flightoperatorreg.azurecr.io"
               az login --service-principal -u ${username} -p ${password} --tenant ${azuretenant}
               
               az acr login --name flightOperatorReg
               az acr repository list --name flightOperatorReg --output table
           """
        }
    }
}