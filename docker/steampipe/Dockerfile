# ---------------------------------------------------
# Stage 1: Build your custom AWS plugin
# ---------------------------------------------------
FROM golang:1.23 AS plugin-builder

WORKDIR /workspace

# Copy the local plugin files to the workspace
COPY . .

# Build the Steampipe plugin binary
RUN go build -tags "netgo" -o aws.plugin *.go

# ---------------------------------------------------
# Stage 2: Final image with Steampipe and your plugin
# ---------------------------------------------------
FROM ubuntu:22.04

# Install minimal runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates \
      curl \
      unzip \
      sudo \
      gnupg && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root 'steampipe' user
RUN useradd -ms /bin/bash steampipe && \
    echo 'steampipe ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER steampipe
WORKDIR /home/steampipe
ENV HOME=/home/steampipe
ENV STEAMPIPE_LOG_LEVEL=TRACE

# Install the Steampipe engine
RUN sudo /bin/sh -c "$(curl -fsSL https://steampipe.io/install/steampipe.sh)"

# Define plugin installation path
ENV PLUGIN_DIR=/home/steampipe/.steampipe/plugins/hub.steampipe.io/plugins/turbot/aws@latest

# Create plugin directory
RUN mkdir -p $PLUGIN_DIR

# Copy the custom-built AWS plugin into the plugin directory
COPY --from=plugin-builder \
     /workspace/aws.plugin \
     $PLUGIN_DIR/aws.plugin

# Ensure the plugin binary is executable
RUN sudo chmod +x $PLUGIN_DIR/aws.plugin

# Create versions.json to register the plugin
RUN echo '{\
  "plugins": {\
    "hub.steampipe.io/plugins/turbot/aws@latest": {\
      "name": "aws",\
      "version": "0.1.0",\
      "struct_version": 20230502\
    }\
  },\
  "struct_version": 20220411\
}' > /home/steampipe/.steampipe/plugins/versions.json

# Launch Steampipe in the foreground
CMD ["steampipe", "service", "start", "--foreground"]
