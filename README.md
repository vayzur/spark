# Spark 🔥

**Spark** is a lightweight HTTP API for [Xray-core](https://github.com/XTLS/Xray-core), providing distributed tunnel control through a simple RESTful interface.

---

## ✨ Features
- Add/Delete inbounds dynamically
- Health check endpoint for monitoring
- Built-in TLS (HTTPS)
- Simple shared-token authentication
- Minimal and production-ready

---

## 🛠 Development

### Requirements
- Go 1.23+
- Xray-core installed locally for testing

### Build
```bash
git clone https://github.com/vayzur/spark.git
cd spark

go build -o spark main.go
```
