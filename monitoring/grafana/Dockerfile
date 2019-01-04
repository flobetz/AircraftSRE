FROM grafana/grafana:5.4.0
ADD ./provisioning /etc/grafana/provisioning
ADD ./dashboards /var/lib/grafana/dashboards
ADD ./plugins /var/lib/grafana/plugins
USER root
RUN chown -R grafana:grafana /var/lib/grafana/
USER grafana