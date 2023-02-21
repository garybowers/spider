FROM golang:1.20 AS build
WORKDIR /app
ADD . ./
RUN go build -o /spiderweb


FROM gcr.io/distroless/base-debian11
COPY --from=build /spiderweb /spiderweb
EXPOSE 8080
CMD ["/spiderweb"]
