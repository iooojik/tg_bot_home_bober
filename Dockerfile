FROM golang:1.21.3 as build
COPY .. /var/www
WORKDIR /var/www/
ARG JOB_TOKEN
ARG JOB_USER
RUN apt install -y git
RUN echo "machine github.com \nlogin ${JOB_USER} \npassword ${JOB_TOKEN}" >> ~/.netrc
RUN git config --global url."https://${JOB_USER}:${JOB_TOKEN}@github.com/".insteadOf https://github.com/
RUN go env -w GOPRIVATE='github.com/*'
