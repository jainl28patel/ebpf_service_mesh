# Root Makefile

.PHONY: all clean app kernel user

# Default target
all: app kernel user

# Build targets for each subdirectory
app:
	# @echo "Building app"
	$(MAKE) -C app

kernel:
	# @echo "Building kernel"
	$(MAKE) -C kernel

user:
	# @echo "Building user"
	$(MAKE) -C user

# Clean targets for each subdirectory
clean:
	# @echo "Cleaning app"
	$(MAKE) -C app clean
	# @echo "Cleaning kernel"
	$(MAKE) -C kernel clean
	# @echo "Cleaning user"
	$(MAKE) -C user clean
