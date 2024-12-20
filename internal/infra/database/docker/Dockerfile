FROM mongo:6.0

ENV MONGO_INITDB_ROOT_USERNAME=admin \
    MONGO_INITDB_ROOT_PASSWORD=admin \
    MONGO_INITDB_DATABASE=auctions

COPY ./mongo-init.js /docker-entrypoint-initdb.d/

EXPOSE 27017

CMD ["mongod"]
