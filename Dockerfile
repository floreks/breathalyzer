FROM floreks/kubepi-base

COPY build/breathalyzer-linux-arm-6 /usr/bin/breathalyzer

ENTRYPOINT ["/usr/bin/breathalyzer"]

EXPOSE 3000
