if [ -f /dependencies/apt-requirements.txt ]; then
  apt-get update
  apt-get install -y $(cat /dependencies/apt-requirements.txt)
fi

exec /main
