from alpine

#WORKDIR /poker.image

ENV NATS_URL='nats://host.docker.internal:4222'
ENV POSTGRES_USER=Rick
ENV POSTGRES_PASSWORD=
ENV POKER_DB_HOST=host.docker.internal

ENV AUTH0_CLIENT_ID=L9efltOhB5V1vleWwUWCOyMwfIoXQ6zd
ENV AUTH0_CLIENT_SECRET=fiXjso1CFRWTAbe8ULttF7J2iJol-maeGrnzuDu87yDTviT5z4p6ePyweYOM_hbk
ENV AUTH0_DOMAIN='https://dummydomain1.auth0.com/'

COPY server.linux /
COPY initdb/* /initdb/
COPY web/* /web/
COPY migrations/* /migrations/
#COPY wait-for-it.sh /wait-for-it.sh
#RUN chmod +x /wait-for-it.sh
#RUN chmod +x /docker-entrypoint.sh

#COPY docker-entrypoint.sh /usr/bin/
#RUN chmod +x /usr/bin/docker-entrypoint.sh
#ENTRYPOINT ["docker-entrypoint.sh"]


#ENTRYPOINT ["./docker-entrypoint.sh"]
#CMD ["node", "main.js"]

CMD ["sh", "-c", "/server.linux"]
#CMD ["sh", "-c", "/bin/sh"]

