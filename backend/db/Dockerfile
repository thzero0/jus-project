FROM postgres:16-alpine
RUN apk add --no-cache python3
COPY . /seed
ENTRYPOINT ["/bin/sh", "/seed/seed.sh"]
