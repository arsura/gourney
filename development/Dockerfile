FROM golang:1.18-alpine

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  APP_ENV=development

RUN apk add --no-cache curl git

WORKDIR /app 

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
  && chmod +x install.sh \
  && sh install.sh \
  && cp ./bin/air /bin/air

CMD air