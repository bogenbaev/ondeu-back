FROM golang:1.18
WORKDIR /app
RUN addgroup appgroup && \
	adduser appuser --ingroup appgroup && \
	chown appuser:appgroup /app
COPY ./ /app
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 4000
CMD ["/app/main"]