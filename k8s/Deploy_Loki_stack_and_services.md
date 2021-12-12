- [Description](#description)
# Description
Deploy Loki(Loki,)
# Requirement
1. K8s cluster is setup
2. Helm is installed
# Deploy Loki stack with Helm
## Get Repo info
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```
## Install grafana 与 prometheus with options
由于我没有集群中未准备PV和PVC，所以就没安装prometheus中与persistentVolume相关的内容
```bash
KUBECONFIG=capi-quickstart.kubeconfig helm upgrade --install loki ./loki-stack  --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```
Login grafana dashboard
```bash
# Get grafana dashboard password
kubectl  --kubeconfig capi-quickstart.kubeconfig get secret loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
# Port forward to access grafana dashboard
kubectl --kubeconfig capi-quickstart.kubeconfig port-forward service/loki-grafana 3000:80
# Login the dashboard on the browser (username:admin, password: retrieved at the previous step)
```
## Deploy the service to allow prometheus probe
### Add meta into pod template spec
```yaml
metadata:
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/path: /metrics
    prometheus.io/port: "8080"
```
The diff as below
```bash
liliu@New-Pro k8s-ansible % git diff deployment/myhttpserver.yaml 
diff --git a/deployment/myhttpserver.yaml b/deployment/myhttpserver.yaml
index 11e18b5..e5ad99c 100644
--- a/deployment/myhttpserver.yaml
+++ b/deployment/myhttpserver.yaml
@@ -13,6 +13,10 @@ spec:
       app.kubernetes.io/name: httpserver
   template:
     metadata:
+        prometheus.io/port: "8080"
+        prometheus.io/scrape: "true"    
       labels:
         app.kubernetes.io/name: httpserver
     spec:
@@ -28,7 +32,7 @@ spec:
       containers:
         - name: httpserver
-          image: sarah4test/httpserver:latest@sha256:e75d85b5086ace7e926e685c80716b657c43a2510485505fef123f99057f1923
+          image: sarah4test/httpserver:latest@sha256:de2a3dffbd3ff510175f2819aadb9d5ad8edb74bde8dc5606baf148ec0d14292
```
### Deploy the service
```bash
kubectl  --kubeconfig capi-quickstart.kubeconfig create -f k8s-ansible/deployment/myhttpserver.yaml 
```
## Access the service
Simplify the process, I port forward, and change the /etc/hosts file
```bash
sudo kubectl  --kubeconfig capi-quickstart.kubeconfig -n ingress-nginx port-forward services/ingress-nginx-controller 443:443
curl https://httpserver.localdev.me/ -k
ok
```
Check /metrics has output
```bash
 curl https://httpserver.localdev.me/metrics -k
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
```
## Generate some test data
```bash
 for i in {1..20}; do curl -k https://httpserver.localdev.me/hello; done;
ok
ok
ok
ok
ok
ok
...
```
## On Grafana dashboard, Verify prometheus scrapped the service metrics 
1. Launch the browser, Logon to http://localhost:3000/
2. From left navagation menu -> Explore -> Select Prometheus from dropdown list.
3. Click Metrics browser -> Type "httpserver_durations_histogram_seconds_bucket" , then click "Use query".
4. You should be able to see some data from the Graph and Table section. 
   
## Create the Httpserver grafana dashboard
The dashboard json is avaible at https://github.com/liuliwh/k8s-ansible/blob/main/deployment/httpserver-grafana-dashboard.json
1. You can import it through Grafana dashboard. (http://localhost:3000/dashboard/import) -> Upload JSON file.
2. Or, you can set the local dashboard when you install the Grapha chart.


