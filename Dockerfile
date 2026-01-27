FROM golang:1.25

RUN mkdir -p /opt/skyflow
WORKDIR /opt/skyflow 

RUN  make build

COPY ./bin/skyflow /opt/skyflow/skyflow
COPY ./kratos_config.yaml /opt/skyflow/kratos_config.yaml
COPY ./build /opt/skyflow

WORKDIR /opt/skyflow

CMD [ "bash" ]