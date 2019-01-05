#!/usr/bin/env bash
docker rm -f jenkins
docker run -p 8080:8080 -p 50000:50000 -v /var/run/docker.sock:/var/run/docker.sock -v ~/data:/var/jenkins_home -d --name jenkins localjenkins:latest