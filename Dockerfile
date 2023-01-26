FROM golang:1.19.0-buster as build

# set go module mode without GOPATH
ENV GO111MODULE=on

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

WORKDIR /usr/src/app/cmd/gameserver

RUN go build -o our-playground-game

FROM golang:1.19.0-buster as stage

# remove apt install lists
RUN apt-get update && \
apt-get -y install uuid-runtime && \
rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app/

COPY --from=build /usr/src/app/cmd/gameserver/our-playground-game /usr/src/app/cmd/gameserver/our-playground-game

RUN chmod +x /usr/src/app/cmd/gameserver/our-playground-game

EXPOSE 6112

CMD [ "/usr/src/app/cmd/gameserver/our-playground-game" ]
