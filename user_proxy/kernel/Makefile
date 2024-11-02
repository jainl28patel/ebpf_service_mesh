# Variables
CLANG = clang
CFLAGS = -target bpf -I/usr/include/linux -I/usr/include/uapi
BPF_FLAGS = -O2 -g -Wall
SOURCE = xdp_main.c
OBJECT = xdp_main.o
OUTPUT = xdp_main

# Targets
all: $(OUTPUT)

$(OUTPUT): $(OBJECT)
	# Move the object file to the final output name
	mv $(OBJECT) $(OUTPUT)

$(OBJECT): $(SOURCE)
	$(CLANG) $(CFLAGS) $(BPF_FLAGS) -c $(SOURCE) -o $(OBJECT)

clean:
	rm -f $(OBJECT) $(OUTPUT)

.PHONY: all clean
