FROM ubuntu

RUN apt update && apt install -fy curl docker
RUN curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/nektos/act/master/install.sh | bash

COPY ./scripts/actrc /root/.config/act/actrc

ENTRYPOINT [ "act" ]