FROM library/golang AS build

MAINTAINER tailinzhang1993@gmail.com

ENV APP_DIR /go/src/helm-client
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o helm-client .
ENTRYPOINT ./helm-client

# Create a minimized Docker mirror
FROM scratch AS prod

COPY --from=build /go/src/helm-client/helm-client /helm-client
CMD ["/helm-client"]
