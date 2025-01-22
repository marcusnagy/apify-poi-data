FROM postgis/postgis:16-3.5-alpine

# Update packages and install prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates \
    gnupg \
    wget \
    build-essential \
    postgresql-server-dev-16 \
    gcc \
    make \
    pgxnclient \
    libssl-dev \
    libcurl4-openssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Enable bullseye-backports for a newer cmake (3.20+)
RUN echo "deb http://deb.debian.org/debian bullseye-backports main" \
    > /etc/apt/sources.list.d/bullseye-backports.list

RUN apt-get update && apt-get install -y --no-install-recommends -t bullseye-backports \
    cmake \
    && rm -rf /var/lib/apt/lists/*

# Install postgresql-16-h3 separately to leverage caching
RUN apt-get update && apt-get install -y --no-install-recommends \
    postgresql-16-h3 \
    && rm -rf /var/lib/apt/lists/*

# Expose PostgreSQL port
EXPOSE 5432

# Set the default command to run PostgreSQL
CMD ["postgres"]