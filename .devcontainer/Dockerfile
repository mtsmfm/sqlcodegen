FROM golang:1.15.5

ARG USERNAME=app
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
  && apt-get update \
  && apt-get install -y sudo postgresql-client \
  && echo "$USERNAME ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/$USERNAME \
  && chmod 0440 /etc/sudoers.d/$USERNAME

RUN mkdir -p /go/src/app /home/app/.config
RUN chown -R $USERNAME /go

USER $USERNAME

WORKDIR /go/src/app

ENV SHELL=/bin/bash
ENV LANG=C.UTF-8 \
  LANGUAGE=C.UTF-8 \
  LC_ALL=C.UTF-8
