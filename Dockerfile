FROM golang:1.7

RUN apt-get update
RUN apt-get install -y \
    dnsmasq \
    file \
    host \
    iptables \
    lsof \
    netcat \
    nginx \
    vim

RUN echo "server {\n listen              443 ssl;\n ssl_certificate     /var/vcap/jobs/gorouter/config/cert.pem;\n ssl_certificate_key /var/vcap/jobs/gorouter/config/key.pem;\n}" > /etc/nginx/conf.d/pcfdev.conf
RUN mkdir -p /var/vcap/packages/uaa/tomcat/conf
RUN echo "<web-app></web-app>" > /var/vcap/packages/uaa/tomcat/conf/web.xml
RUN mkdir -p /var/pcfdev/api
RUN mkdir -p /var/vcap/monit/job
RUN mkdir -p /var/vcap/jobs/garden/bin
RUN mkdir -p /var/vcap/jobs/uaa/config
RUN echo "exec /var/vcap/packages/garden-linux/bin/garden-linux \\ \n  -dnsServer=some-dns-server \\ \n  1>>\$LOG_DIR/garden.stdout.log \\ \n  2>>\$LOG_DIR/garden.stderr.log" > /var/vcap/jobs/garden/bin/garden_ctl
RUN echo "scim:\n  users:\n  - admin|admin|scim.write,scim.read,openid" > /var/vcap/jobs/uaa/config/uaa.yml
RUN ln -s /bin/true /usr/local/bin/resolvconf

RUN mkdir -p /var/vcap/bosh
RUN echo "" > /var/vcap/bosh/agent_state.json
RUN mkdir -p /var/vcap/jobs/gorouter/config
RUN mkdir -p /var/vcap/jobs/cloud_controller_ng/bin
RUN echo "some-gorouter-cert" > /var/vcap/jobs/gorouter/config/cert.pem
RUN echo "some-pcfdev-trusted-ca" > /var/pcfdev/trusted_ca.crt
RUN echo "" > /var/vcap/jobs/cloud_controller_ng/bin/cloud_controller_worker_ctl

RUN mkdir -p /var/vcap/jobs/cflinuxfs2-rootfs-setup/bin
RUN echo '#!/bin/bash\necho "some-cflinuxfs2-rootfs-setup-prestart"' > /var/vcap/jobs/cflinuxfs2-rootfs-setup/bin/pre-start
RUN chmod +x /var/vcap/jobs/cflinuxfs2-rootfs-setup/bin/pre-start
