FROM hub.tess.io/sherlockio/debian:stretch

RUN mkdir -p /platform

COPY ./output/crd_controller /bin/crd_controller

WORKDIR /platform

ENTRYPOINT [ "/bin/crd_controller" ]
CMD []