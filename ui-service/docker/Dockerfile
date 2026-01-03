# Bar Restaurant UI Service Dockerfile
FROM nginx:alpine

# Install curl for healthcheck
RUN apk --no-cache add curl

# Copy nginx configuration
COPY docker/nginx.conf /etc/nginx/nginx.conf

# Copy UI files
COPY *.html /usr/share/nginx/html/
COPY *.js /usr/share/nginx/html/
COPY shared/ /usr/share/nginx/html/shared/

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
