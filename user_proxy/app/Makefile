# Variables
CC = gcc
CFLAGS = -Wall -Wextra -O2
LDFLAGS =
SOURCES = udp_echo_server.c
OBJECTS = $(SOURCES:.c=.o)
TARGET = udp_echo_server

# Targets
all: $(TARGET)

$(TARGET): $(OBJECTS)
	$(CC) $(LDFLAGS) -o $@ $(OBJECTS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJECTS) $(TARGET)

.PHONY: all clean
