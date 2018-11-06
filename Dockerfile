FROM registry.fit2cloud.com/jumpserver/python:v3
MAINTAINER Fit2anything Team <ibuler@qq.com>

WORKDIR /opt/fit2ansible

COPY ./requirements /tmp/requirements

RUN rm -f /usr/bin/python && ln -s /usr/bin/python2 /usr/bin/python

RUN yum -y install epel-release && cd /tmp/requirements && \
    yum -y install $(cat rpm_requirements.txt)
RUN cd /opt && python3 -m venv py3
RUN cd /tmp/requirements && /opt/py3/bin/pip install --upgrade pip setuptools && \
    /opt/py3/bin/pip install -r requirements.txt -i https://mirrors.ustc.edu.cn/pypi/web/simple
RUN sed -i "s@'uri': True@'uri': False@g" /opt/py3/lib/python3.6/site-packages/django/db/backends/sqlite3/base.py
ENV LANG=zh_CN.UTF-8
ENV LC_ALL=zh_CN.UTF-8
ENV VENV=/opt/py3
ENV APP_DIR=/opt/fit2ansible
ENV PYTHONOPTIMIZE=1
ENV C_FORCE_ROOT=1

RUN mkdir -p /root/.ssh/
RUN echo "ClusterHost *\n  StrictHostKeyChecking no\n  UserKnownHostsFile=/dev/null" > /root/.ssh/config

COPY . /opt/fit2ansible
VOLUME /opt/fit2ansible/data

EXPOSE 8000
CMD ["bash", "-c", "python3 entrypoint.py start"]