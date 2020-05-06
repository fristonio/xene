FROM golang

# For generating certificates
RUN git clone https://github.com/cloudflare/cfssl/ /go/src/github.com/cloudflare/cfssl && \
    cd /go/src/github.com/cloudflare/cfssl/ && \
    make && \
    mv bin/* /usr/local/bin/

ADD . /go/src/github.com/fristonio/xene
WORKDIR /go/src/github.com/fristonio/xene

# Build and generate certificates for xene.
RUN make build && \
    make -C contrib/certs/ certs && \
    mkdir -p /etc/xene/certs && \
    mv contrib/certs/*.gen /etc/xene/certs/ && \
    mv /go/bin/* /usr/local/bin/

EXPOSE 6060

ENTRYPOINT ["xene", "apiserver", "-n"]
