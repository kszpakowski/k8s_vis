FROM scratch
COPY ./out/app /app
ENTRYPOINT /app