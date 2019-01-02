#!groovy​

node("master"){

    stage("Checkout") {
        checkout([$class: 'GitSCM', branches: [[name: '*/master']],
                  doGenerateSubmoduleConfigurations: false,
                  extensions: [[$class: 'LocalBranch'], [$class: 'CleanBeforeCheckout'],
                               [$class: 'IgnoreNotifyCommit']],
                  submoduleCfg: [],
                  userRemoteConfigs: [[credentialsId: 'github-credentials', url: 'https://github.com/flobetz/AircraftSRE.git']]])
    }

    stage("Build") {
        sh "ls -lsa"
    }
}