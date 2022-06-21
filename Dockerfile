FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go get -u github.com/mailru/easyjson/...
ENV PATH=$GOPATH/bin:$PATH

COPY . .
RUN go mod tidy
RUN easyjson -all -pkg src/models
RUN go build -o api.run ./cmd/.

FROM ubuntu

RUN apt -y update && \
    echo "tzdata "Geographic area" select 8" | debconf-set-selections && \
    apt install -y tzdata && apt install -y postgresql

ENV USER user
ENV DB_NAME forum

USER postgres

RUN /etc/init.d/postgresql start && \
    psql --command "CREATE USER $USER WITH SUPERUSER PASSWORD 'test_user';" && \
    createdb --owner=$USER $DB_NAME && \
    /etc/init.d/postgresql stop

ENV ARTIFACT api.run

WORKDIR /app
COPY --from=build /app/$ARTIFACT $ARTIFACT
COPY ./db/db.sql db.sql

ENV MODE release

CMD service postgresql start && \
    psql -h localhost -p 5432 -d $DB_NAME -U $USER -w -q -f db.sql \
    && ./$ARTIFACT

VOLUME ["/var/lib/postgresql/data"]
EXPOSE 5000