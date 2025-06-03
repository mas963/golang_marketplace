FROM ubuntu:latest
LABEL authors="yasar"

ENTRYPOINT ["top", "-b"]