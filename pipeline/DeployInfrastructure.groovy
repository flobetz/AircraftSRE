#!groovy​

node("master"){
    ansiColor('xterm') {
        stage("Checkout") {
            checkout scm
        }

        dir("INFRA") {
            dir("terraform") {
                stage("Deploy Azure Infrastructure") {
                    withCredentials([string(credentialsId: 'azuretenant', variable: 'azuretenant'),
                                     usernamePassword(credentialsId: 'azureServicePrincipal', passwordVariable: 'password', usernameVariable: 'username'),
                                     usernamePassword(credentialsId: 'azureStorageAccount', passwordVariable: 'storageAccountPassword', usernameVariable: 'storageAccountName')]) {
                        sh """
                            az login --service-principal -u ${username} -p ${password} --tenant ${azuretenant}
                            echo "Setting environment variables for Terraform"
                            export ARM_SUBSCRIPTION_ID=d354bff2-0995-4811-a06e-7800ca7b5d37
                            export ARM_CLIENT_ID=${username}
                            export ARM_CLIENT_SECRET=${password}
                            export ARM_TENANT_ID=${azuretenant}
                            terraform init -backend-config="storage_account_name=${storageAccountName}" -backend-config="container_name=tfstate" -backend-config="access_key=${storageAccountPassword}" -backend-config="key=codelab.microsoft.tfstate"
                    
                            export TF_VAR_client_id=${username}
                            export TF_VAR_client_secret=${password}

                            terraform plan -out out.plan
                            terraform apply out.plan
                            
                            echo "\$(terraform output kube_config)" > terraformOut.txt
                        """
                    }

                    archiveArtifacts 'terraformOut.txt'
                }
            }
        }
    }
}