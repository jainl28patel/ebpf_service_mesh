# Use an official Ubuntu image as the base image
FROM ubuntu:latest

WORKDIR /tmp

COPY . .

# Make the run.sh script executable
RUN chmod +x install.sh

# Execute the run.sh script to install dependencies and build the application
RUN ./install.sh
