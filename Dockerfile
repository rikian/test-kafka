FROM golang
WORKDIR /app
COPY . .
EXPOSE 9095
RUN go mod tidy
RUN go build -o goapp .
# change to none root user
# RUN useradd -u 10001 appuser
# RUN chown -R appuser:appuser /app
# USER appuser
CMD ["./goapp"]
