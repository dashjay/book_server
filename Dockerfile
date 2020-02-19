FROM golang:1.13.7

WORKDIR /book_server
COPY . .
EXPOSE 80


CMD ["/bin/bash","/book_server/build.sh"]

