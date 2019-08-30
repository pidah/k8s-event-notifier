# k8s-event-notifier
A simple k8s events notifier to slack.

# Overview
The kubernetes Events API streams a lot of information, useful to Operations folk and application Developers. This simple tool allows you to filter the events object based on [reason](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#event-v1-core) as follows:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-event-notifier
  namespace: kube-system
  labels:
    app: k8s-event-notifier
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8s-event-notifier
  template:
    metadata:
      labels:
        app: k8s-event-notifier
    spec:
      containers:
      - name: k8s-event-notifier
        env:
        - name: EVENT_FILTER
          value: BackOff,Failed,FailedScheduling,Killing
        - name: SLACK_API_URL
          value: https://hooks.slack.com/services/RANDOM12/BFMADBNDS/0cYHvK1a9xviCsjakjawTMaa
        image: pearsontechnology/k8s-event-notifier:1
        command: ["/k8s-event-notifier"]
        args:
        - --logtostderr
        - -v=4
         2>&1
```

You can pass a comma-separated list of event reasons to the EVENT_FILTER as shown above. Note that by default, if no EVENT_FILTER is specified, no event will be captured.

Once an interesting event is detected, a notification will be sent to slack using the specified slack incoming webhook url `SLACK_API_URL` as shown above.

# Install

Checkout this repo and you can use this example [deployment manifest file](https://github.com/pearsontechnology/k8s-event-notifier/blob/master/manifests/deployment.yaml) as follows:

```
kubectl create -f manifests/deployment.yaml
```

and a sample [RBAC manifest file](https://github.com/pearsontechnology/k8s-event-notifier/blob/master/manifests/rbac.yaml) if required:

```
kubectl create -f manifests/rbac.yaml
```
