FROM python:3.9-slim

COPY docs/gh-pages/requirements.txt /docs/requirements.txt
RUN pip install -r /docs/requirements.txt

COPY docs/gh-pages /docs/
RUN ls /docs/

WORKDIR /docs
ENTRYPOINT ["mkdocs"]
CMD ["serve"]
EXPOSE 8000


docker run --rm -it -v ./docs/gh-pages:/docs -w /docs python:3.9-slim sh -c "pip install -r requirements.txt && mkdocs build --config-file mkdocs.yml"