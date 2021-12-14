- [Description](#description)
- [Requirement](#requirement)
- [Deploy the test service](#deploy-the-test-service)
- [Deploy stack with Helm](#deploy-stack-with-helm)
  - [Get Repo info](#get-repo-info)
  - [Install grafana 与 prometheus with options](#install-grafana-与-prometheus-with-options)
  - [Update test service yaml file if needed](#update-test-service-yaml-file-if-needed)
  - [Deploy ServiceMonitor](#deploy-servicemonitor)
  - [Access the service and generate test data](#access-the-service-and-generate-test-data)
  - [On prometheus UI, verify the service is shown in the targets](#on-prometheus-ui-verify-the-service-is-shown-in-the-targets)
  - [Deploy the dashboard for test service](#deploy-the-dashboard-for-test-service)
  - [On Grafana dashboard, Verify prometheus scrapped the service metrics and dashboard](#on-grafana-dashboard-verify-prometheus-scrapped-the-service-metrics-and-dashboard)
# Description
Deploy prometheus operator, prometheus and grafana.
# Requirement
1. K8s cluster is setup
2. Helm is installed
# Deploy the test service
The finished yaml file is located at https://github.com/liuliwh/k8s-ansible/blob/main/deployment/myhttpserver.yaml
```bash
kubectl create -f k8s-ansible/deployment/myhttpserver.yaml 
```
# Deploy stack with Helm
## Get Repo info
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```
## Install grafana 与 prometheus with options
根据文档 https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack#configuration 描述，prometheus operator 默认只discover 与 prometheus operator 相同namespace且lable release tag与其相同的service与pod。所以修改以下两个设置项为false.
```bash
helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false,prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false
```
## Update test service yaml file if needed
由于prometheus operator不是通过annotation的方式做service discovery，而是使用PodMonitor or ServiceMonitor CRD 配置，所以原本的service yaml无需改变（但如果原本service未加label的话，需要加上label，以便servicemonitor 的selector使用)
原yaml 已有label,故无需修改
## Deploy ServiceMonitor
Service是通过ServiceMonitor对象来发现，所以仅需要将spec/selector/matchLabels与Service对应即可。
The finished yaml file is located at https://github.com/liuliwh/k8s-ansible/blob/main/deployment/myhttpserver-servicemonitor.yaml
## Access the service and generate test data
Simplify the process, I use port forward
```bash
kubectl  port-forward service/httpserver-service 8080:80
for i in {1..20}; do curl -H 'Host:httpserver.localdev.me' localhost:8080/ok ; done;
ok
ok
ok
ok
ok
ok
...
```
## On prometheus UI, verify the service is shown in the targets
```bash
kubectl port-forward service/kube-prometheus-stack-prometheus 9090:9090
```
Launch the web browser, to access localhost:9090/targets, serviceMonitor/default/httpserver/0 (1/1 up) should be shown.
## Deploy the dashboard for test service
根据https://github.com/grafana/helm-charts/blob/main/charts/grafana/README.md 以及以下命令可看出，grafana通过sidecar导入dashboar默认是enabled的，以及label是"grafana_dashboard"。
```bash
k8s-ansible % helm get values kube-prometheus-stack -a | grep -A 30 sidecar 
  sidecar:
    dashboards:
      ...
      enabled: true
      ...
      label: grafana_dashboard
```
所以，我们仅需要通过configmap将dashboard导入即可
```bash
kubectl apply -f https://raw.githubusercontent.com/liuliwh/k8s-ansible/main/deployment/myhttpserver-grafana-dashboard.yaml
```
## On Grafana dashboard, Verify prometheus scrapped the service metrics and dashboard
Login grafana dashboard login password
```bash
# Get grafana dashboard password
kubectl get secret kube-prometheus-stack-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
# Port forward to access grafana dashboard
kubectl port-forward service/kube-prometheus-stack-grafana 3000:80
# Login the dashboard on the browser (username:admin, password: retrieved at the previous step)
```
1. Launch the browser, Logon to http://localhost:3000/
2. From left navagation menu -> Explore -> Select Prometheus from dropdown list.
3. Click Metrics browser -> Type "httpserver_durations_histogram_seconds_bucket" , then click "Use query".
4. You should be able to see some data from the Graph and Table section. 
5. Navigate to Dashboard -> Browse, you will see Httpserver dashboard with tag myhttpserver.
