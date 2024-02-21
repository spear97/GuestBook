# Use Ubuntu 20.04 as base image
FROM ubuntu:20.04

# Set environment variable to prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Update package lists and install necessary packages
RUN apt-get update -qq && \
    apt-get install -y -q --no-install-recommends \
    openssl python3.7 python3-pip

# Upgrade pip and setuptools
RUN pip3 install --upgrade pip setuptools

# Install required Python packages
RUN pip3 install requests Flask Flask-Cache pyopenssl ibm-watson==6.0.0

# Copy application source code into the container (assuming it's in a 'src' directory)
ADD src /app/src

# Expose port 5000 for the Flask application
EXPOSE 5000

# Set the entrypoint command to run the Flask application
ENTRYPOINT ["python3", "/app/src/app.py"]
