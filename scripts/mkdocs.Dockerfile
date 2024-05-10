FROM python:3.9-slim

COPY docs/gh-pages/requirements.txt /docs/requirements.txt
RUN pip install -r /docs/requirements.txt

COPY docs/gh-pages /docs/
RUN ls /docs/

WORKDIR /docs
ENTRYPOINT ["mkdocs"]
CMD ["serve"]
EXPOSE 8000