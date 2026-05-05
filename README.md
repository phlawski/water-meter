# Water Meter

A simple home utility billing app running on Kubernetes (K3s on Raspberry Pi 5).

Tracks monthly water meter readings, calculates usage and cost in PLN, and accounts for a fixed internet cost deduction.

## Stack

- Go + chi + html/template + Pico CSS + Chart.js
- SQLite (persisted via PVC)
- K3s on Raspberry Pi 5 (Ubuntu 24.04, aarch64)

## Local development

```bash
go run ./cmd/server
# open http://localhost:8080
```

## Build Docker image

```bash
# For Raspberry Pi (arm64)
docker build --platform linux/arm64 -t phlawski/water-meter:latest .

# Push to Docker Hub
docker push phlawski/water-meter:latest
```

## Deploy to K3s

```bash
# First time
kubectl apply -f k8s/deployment.yaml

# Update after pushing a new image
kubectl rollout restart deployment/water-meter
```

## Access

```
http://peter-desktop.local
```

## Useful kubectl commands

```bash
# Check pod status
kubectl get pods

# Check logs
kubectl logs -l app=water-meter -f

# Check persistent volume
kubectl get pvc

# Delete and redeploy everything
kubectl delete -f k8s/deployment.yaml
kubectl apply -f k8s/deployment.yaml
```
