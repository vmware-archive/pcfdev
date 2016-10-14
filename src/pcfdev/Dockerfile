FROM golang:1.6

RUN apt-get update
RUN apt-get install -y nginx vim host dnsmasq

RUN echo "server {\n listen              443 ssl;\n ssl_certificate     /var/vcap/jobs/gorouter/config/cert.pem;\n ssl_certificate_key /var/vcap/jobs/gorouter/config/key.pem;\n}" > /etc/nginx/conf.d/pcfdev.conf
RUN mkdir -p /var/vcap/packages/uaa/tomcat/conf
RUN echo "<web-app></web-app>" > /var/vcap/packages/uaa/tomcat/conf/web.xml
RUN ln -s /bin/true /usr/local/bin/resolvconf