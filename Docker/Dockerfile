# build stage
FROM golang:alpine AS builder
WORKDIR /src
ADD . /src
RUN go build -o exam main.go

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /src/exam /app/
CMD ["./exam"]
# ENTRYPOINT ./exam
