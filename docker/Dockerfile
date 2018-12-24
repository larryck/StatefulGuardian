FROM busybox

COPY bin/linux_amd64/* /

RUN  mkdir /etc/kubeconfig/

ENTRYPOINT ["/statefulguardian"]
