#!/usr/bin/sh
echo docker run -u "$(id -u):$(id -g)" --rm $@
docker run --rm $@
