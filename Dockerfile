FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go get -u github.com/mailru/easyjson/...
ENV PATH=$GOPATH/bin:$PATH

COPY . .
RUN go mod tidy
RUN easyjson -all -pkg app/models
RUN go build -o api.run ./cmd/.

FROM ubuntu
ENV PGVER 12
RUN apt -y update && \
    echo "tzdata "Geographic area" select 8" | debconf-set-selections && \
    apt install -y tzdata && apt install -y postgresql-$PGVER

ENV PGDEFAULT_USER postgres
ENV PGFORUM_USER forum_user
ENV PGPASSWORD forum_user_password
ENV PGDB_NAME forum
ENV PGPORT 5432
ENV API_PORT 5000

USER $PGDEFAULT_USER

RUN /etc/init.d/postgresql start && \
    psql --command "CREATE USER $PGFORUM_USER WITH SUPERUSER PASSWORD '$PGPASSWORD';" && \
    createdb --owner=$PGFORUM_USER $PGDB_NAME && \
    /etc/init.d/postgresql stop

ENV ARTIFACT api.run

WORKDIR /app
COPY --from=build /app/$ARTIFACT $ARTIFACT
COPY ./db/db.sql db.sql

ENV MODE release
#TODO: docker-compose
CMD service postgresql start && \
    psql -h localhost -p $PGPORT -d $PGDB_NAME -U $PGFORUM_USER -w -q -f db.sql \
    && ./$ARTIFACT

VOLUME ["/var/lib/postgresql/data"]
EXPOSE $API_PORT
