# Use an official Ubuntu image as the base image
FROM anakin007/service-mesh:latest

# Set environment variables to non-interactive (avoid prompts during package installation)
ENV DEBIAN_FRONTEND=noninteractive

# Set the working directory in the container
WORKDIR /app

# Copy the C source code, Makefile, and the run.sh script into the container
COPY . .

# Make the run.sh script executable
RUN chmod +x run.sh
RUN chmod +x build.sh

# Execute the run.sh script to install dependencies and build the application
RUN ./build.sh

CMD [ "./run.sh" ]