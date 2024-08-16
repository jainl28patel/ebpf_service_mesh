# Use an official Ubuntu image as the base image
FROM ubuntu:latest

# Set environment variables to non-interactive (avoid prompts during package installation)
ENV DEBIAN_FRONTEND=noninteractive

# Set the working directory in the container
WORKDIR /app

# Copy the C source code, Makefile, and the run.sh script into the container
COPY . .

# Make the run.sh script executable
RUN chmod +x run.sh

# Execute the run.sh script to install dependencies and build the application
RUN ./run.sh

# Command to run the application
CMD ["./app/udp_echo_server"]
