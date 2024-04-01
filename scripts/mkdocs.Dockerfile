FROM python:3.9-slim

COPY docs/requirements.txt /docs/requirements.txt
RUN pip install -r /docs/requirements.txt

COPY docs /docs/

WORKDIR /docs
ENTRYPOINT ["mkdocs"]
CMD ["serve"]
EXPOSE 8000