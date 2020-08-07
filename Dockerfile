FROM debian:buster

RUN apt-get update && \
    apt-get upgrade -y

# Deploy component
RUN mkdir /app
COPY server /app

CMD [ "/app/server" ]