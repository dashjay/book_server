FROM golang:alpine

WORKDIR /book_server
COPY . .

EXPOSE 80
CMD ["sh","/book_server/build.sh"]

