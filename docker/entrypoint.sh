export http_proxy=${HTTP_PROXY}
export https_proxy=${HTTPS_PROXY}

if [ -f /dependencies/apt-requirements.txt ]; then
  apt-get update
  apt-get install -y $(cat /dependencies/apt-requirements.txt)
fi

exec /main
