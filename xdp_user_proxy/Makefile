# Root Makefile

.PHONY: all clean app kernel user

# Default target
all: app kernel user

# Build targets for each subdirectory
app:
	$(MAKE) -C app

kernel:
	$(MAKE) -C kernel


# Clean targets for each subdirectory
clean:
	$(MAKE) -C app clean
	$(MAKE) -C kernel clean
