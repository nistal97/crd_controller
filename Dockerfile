FROM hub.tess.io/sherlockio/debian:stretch

RUN mkdir -p /platform

RUN pwd
RUN ls

COPY output/crd_controller /bin/crd_controller

WORKDIR /platform

ENTRYPOINT [ "/bin/crd_controller" ]
CMD []