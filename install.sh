#!/bin/bash

set -e

GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
RED="\033[0;31m"
RESET="\033[0m"

SPARK_DIR="/etc/spark"
CONFIG_FILE="/etc/spark/config.toml"
SERVICE_FILE="/etc/systemd/system/spark.service"
DOMAIN="$1"

detect_arch() {
  case "$(uname -m)" in
    x86_64) echo "amd64" ;;
    aarch64) echo "arm64" ;;
    *) echo "unsupported" ;;
  esac
}

detect_distro() {
  if command -v apt-get >/dev/null 2>&1; then
    echo "debian"
  elif command -v yum >/dev/null 2>&1; then
    echo "redhat"
  else
    echo "unsupported"
  fi
}

ARCH=$(detect_arch)
DISTRO=$(detect_distro)

if [[ "$ARCH" == "unsupported" || "$DISTRO" == "unsupported" ]]; then
  echo "Unsupported architecture or distribution"
  exit 1
fi

echo -e "${GREEN}[INFO] Installing required packages...${RESET}"

if [[ "$DISTRO" == "debian" ]]; then
  apt-get update
  apt-get install -yqq unzip wget certbot
elif [[ "$DISTRO" == "redhat" ]]; then
  yum install -yq unzip wget certbot
fi

echo -e "${GREEN}[INFO] Getting TLS certificates...${RESET}"

certbot certonly -q --standalone --keep --preferred-challenges http -d "${DOMAIN}" --non-interactive --agree-tos --register-unsafely-without-email

cd /tmp

echo -e "${GREEN}[INFO] Downloading binary...${RESET}"

wget -q -O /tmp/spark.zip "https://github.com/vayzur/spark/releases/latest/download/spark-linux-${ARCH}.zip"

unzip -qo /tmp/spark.zip -d /tmp

install -m 755 /tmp/spark /usr/local/bin/spark

mkdir -p "${SPARK_DIR}"

wget -q -O "${SERVICE_FILE}" https://raw.githubusercontent.com/vayzur/spark/main/spark.service
wget -q -O "${CONFIG_FILE}" https://raw.githubusercontent.com/vayzur/spark/main/config.toml

SECRET=$(cat /proc/sys/kernel/random/uuid)

sed -i "s|sub\.domain\.tld|${DOMAIN}|g" "${CONFIG_FILE}"
sed -i "s|secret = \".*\"|secret = \"${SECRET}\"|g" "${CONFIG_FILE}"

systemctl daemon-reload
systemctl enable --now spark > /dev/null 2>&1

echo -e "${GREEN}[INFO] Done. Spark is running!${RESET}"
