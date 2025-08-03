# Spark 🔥

**Spark** is an HTTP API for [Xray-core](https://github.com/XTLS/Xray-core), enabling distributed control of tunnels with a simple RESTful interface.

---

## ✨ Key Features

- Add/Delete inbounds
- Health check endpoint for monitoring
- Built-in TLS support (HTTPS)
- Simple shared-secret authentication
- Minimal and production-ready

---

## ⚙️ Configuration

Spark uses a TOML-based configuration file.

### Example `config.toml`:

```toml
[tls]
enabled = true
cert_file = "/etc/letsencrypt/live/sub.domain.tld/fullchain.pem"
key_file = "/etc/letsencrypt/live/sub.domain.tld/privkey.pem"

[server]
addr = ":10100"
prefork = true

[xray]
addr = "127.0.0.1:8080"

[auth]
secret = "your-shared-secret"
```

---

## 🔐 Authentication

All requests must include a valid `Authorization` header.

### Format:
```
Authorization: rolling ts:hash
```

- `ts`: Current Unix timestamp (in seconds)
- `hash`: SHA256 hash of `{timestamp}:{secret}`

> ⚠️ The hash expires **1 minute** after generation. Clients must generate a new one for each request.

### Example (bash):
```bash
SECRET="your-shared-secret"
TS=$(date +%s)
HASH=$(echo -n "${TS}:${SECRET}" | sha256sum | cut -d ' ' -f1)
curl -H "Authorization: rolling ${TS}:${HASH}" http://localhost:10100/health
```

---

## 🚀 Deployment

You can deploy Spark using **Docker** or as a **standalone binary with systemd**.

---

### 🐳 Deploy with Docker

1. Create or edit `config.toml`:
   ```bash
   vim config.toml
   ```

2. Edit `compose.yml` to match your port settings:
   ```bash
   vim compose.yml
   ```

3. Start the service:
   ```bash
   docker compose up -d
   ```

---

### 🧩 Deploy with Binary & systemd

Use the official install script:

```bash
bash <(curl -Ls https://raw.githubusercontent.com/vayzur/spark/main/install.sh) sub.domain.tld
```

> 🌐 Ensure your subdomain points to your server's public IP.

After installation:

- Edit default config:
  ```bash
  vim /etc/spark/config.toml
  ```

- Check service status:
  ```bash
  systemctl status spark
  ```

- Default port is `10100`. Make sure it's not in use:
  ```bash
  ss -tlpn
  ```

---

## 📚 API Documentation

Check the `docs/` directory for example requests and usage.

---

## 🛠️ Notes

- Spark listens on all interfaces by default (`0.0.0.0`)
- TLS is recommended for production use
- Always keep your `secret` safe and rotate it regularly
