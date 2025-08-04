import time
import hashlib
import requests
import json

# --- Configuration ---
# Replace with your actual Worker URL and shared secret
WORKER_URL = "https://sub.domain.tld:10100"
SHARED_SECRET = "your_shared_secret_here"

# --- Authentication Helper ---
def generate_auth_header(secret: str) -> str:
    """
    Generates the 'Authorization' header value using SHA256 hashing.
    Format: "rolling <timestamp>:<signature>"
    """
    timestamp = str(int(time.time()))
    data_to_sign = f"{timestamp}:{secret}"
    signature = hashlib.sha256(data_to_sign.encode()).hexdigest()
    return f"rolling {timestamp}:{signature}"

# --- Prepare Authentication ---
auth_header = generate_auth_header(SHARED_SECRET)
headers = {
    "Authorization": auth_header,
    "Content-Type": "application/json"
}

# --- Define the Inbound Configuration ---
# This is a VLESS over WebSocket inbound.
# See: https://xtls.github.io/en/config/inbound.html#inboundobject
inbound_config = {
    "listen": "0.0.0.0",
    "port": 11882,
    "protocol": "vless",
    "settings": {
        "clients": [
            {
                "id": "6211f9a6-a8e6-421d-ab27-7091655345e6" # UUID for the client
            }
        ],
        "decryption": "none",
        "fallbacks": []
    },
    "streamSettings": {
        "network": "ws",
        "security": "none",
        "wsSettings": {
            "acceptProxyProtocol": False,
            "headers": {},
            "heartbeatPeriod": 0,
            "host": "",
            "path": "/"
        }
    },
    "tag": "inbound-11882", # Unique tag for this inbound
    "sniffing": {
        "enabled": False,
        "destOverride": ["http", "tls", "quic", "fakedns"],
        "metadataOnly": False,
        "routeOnly": False
    },
    "allocate": {
        "strategy": "always",
        "refresh": 5,
        "concurrency": 3
    }
}

# --- API Calls ---

# 1. Health Check
print("--- Health Check ---")
try:
    health_resp = requests.get(f"{WORKER_URL}/healthz", headers=headers)
    print(f"Status Code: {health_resp.status_code}")
except requests.exceptions.RequestException as e:
    print(f"Error during health check: {e}")

# 2. Add New Inbound
print("\n--- Add New Inbound ---")
try:
    inbound_json = json.dumps(inbound_config)
    add_resp = requests.post(f"{WORKER_URL}/inbounds", headers=headers, data=inbound_json)
    print(f"Status Code: {add_resp.status_code}")
    print(f"Response: {add_resp.text}")
except requests.exceptions.RequestException as e:
    print(f"Error adding inbound: {e}")
except json.JSONEncodeError as e:
    print(f"Error encoding inbound config to JSON: {e}")

# 3. Remove Inbound (using the tag defined above)
print("\n--- Remove Inbound ---")
tag = inbound_config["tag"] # Use the tag from the config
try:
    delete_resp = requests.delete(f"{WORKER_URL}/inbounds/{tag}", headers=headers)
    print(f"Status Code: {delete_resp.status_code}")
    print(f"Response: {delete_resp.text}")
except requests.exceptions.RequestException as e:
    print(f"Error removing inbound: {e}")
