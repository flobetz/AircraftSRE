#!groovy​

node("master"){

    stage("Checkout") {
        checkout scm
    }

    stage("Build") {
        sh "docker build -t flightOperator:latest ."
    }

    stage("Push") {
        sh "docker push flightOperator:latest"
    }
}