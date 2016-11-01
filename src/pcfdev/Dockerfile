RUN apt-get update
RUN apt-get install -y \
    dnsmasq \
    host \
    iptables \
    nginx \
    vim \
    lsof \
    netcat

RUN echo "server {\n listen              443 ssl;\n ssl_certificate     /var/vcap/jobs/gorouter/config/cert.pem;\n ssl_certificate_key /var/vcap/jobs/gorouter/config/key.pem;\n}" > /etc/nginx/conf.d/pcfdev.conf
RUN mkdir -p /var/vcap/packages/uaa/tomcat/conf
RUN echo "<web-app></web-app>" > /var/vcap/packages/uaa/tomcat/conf/web.xml
RUN mkdir -p /var/vcap/jobs/garden/bin
RUN echo "exec /var/vcap/packages/garden-linux/bin/garden-linux \\ \n  -dnsServer=some-dns-server \\ \n  1>>\$LOG_DIR/garden.stdout.log \\ \n  2>>\$LOG_DIR/garden.stderr.log" > /var/vcap/jobs/garden/bin/garden_ctl
RUN ln -s /bin/true /usr/local/bin/resolvconf
