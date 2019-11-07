# First stage, build the executable
FROM golang:1.13

ENV GOPROXY=https://goproxy.cn
ENV HOME=/opt/app

WORKDIR $HOME

COPY go.mod $HOME
COPY go.sum $HOME
RUN go mod download

COPY . $HOME
# Use static linking to get rid of the error below
# exec user process caused "no such file or directory"
RUN GOOS=linux GOARCH=amd64 go build -a -ldflags "-linkmode external -extldflags '-static' -s -w"

# Second stage for executable only image
FROM scratch

RUN mkdir /opt/app/cache
RUN if [ -f /opt/app/cache/emb_short.json ]; then echo "emb_short.json exist"; else curl -o /opt/app/cache/emb_short.json http://47.93.237.110/emb_short.json; fi
RUN if [ -f /opt/app/cache/embVectors.json ]; then echo "embVectors.json exist"; else curl -o /opt/app/cache/embVectors.json http://47.93.237.110/embVectors.json; fi
# RUN curl -o /opt/app/cache/embVectors.json http://47.93.237.110/embVectors.json 

# Copy executable from the first stage
COPY --from=0 /opt/app/backend /backend
COPY --from=0 /opt/app/ind_keyword.ind /ind_keyword.ind
COPY --from=0 /opt/app/ind_name.ind /ind_name.ind
COPY --from=0 /opt/app/ind_title.ind /ind_title.ind
COPY --from=0 /opt/app/info.json /info.json
COPY --from=0 /tmp/emb_short.json /emb_short.json
COPY --from=0 /tmp/embVectors.json /embVectors.json

EXPOSE 80

CMD ["/backend"]
