FROM redis:7.4.2-alpine

COPY redis.conf /usr/local/etc/redis/redis.conf

CMD ["redis-server", \
     "--save", "60", "1", \
     "--loglevel", "warning"]