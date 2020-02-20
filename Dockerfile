FROM golang:alpine
WORKDIR /book_server
EXPOSE 80
CMD ["sh","/book_server/build.sh"]

