#!/bin/bash

set -e

# Function to install for Linux
install_linux() {
    # Download TensorFlow library
    wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz

    # Extract TensorFlow library
    tar -C /usr/local -xzf ./libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz

    # Update linker cache
    /sbin/ldconfig -v

    # Update package lists
    apt-get update
    apt-get upgrade -y
    apt-get dist-upgrade -y

    # Install image processing libraries with automatic yes to prompts
    apt-get install -y libjpeg-dev libpng-dev libtiff-dev libgif-dev libwebp-dev
    go get github.com/kolesa-team/go-webp

    echo "Installation for Linux completed successfully."
}

# Function to install for Windows
install_windows() {
    # Define cross-compiler prefix
    CROSS_PREFIX=x86_64-w64-mingw32-

    # Install cross-compile tools
    apt-get update
    apt-get install -y mingw-w64

    # Download and extract libwebp
    wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.2.0.tar.gz
    tar -xzf libwebp-1.2.0.tar.gz
    cd libwebp-1.2.0

    # Configure for cross-compilation for Windows
    ./configure --host=${CROSS_PREFIX%?} --prefix=/usr/${CROSS_PREFIX%?} --enable-shared --enable-static

    # Compile and install
    make
    make install

    # Return to the original directory
    cd ..

    echo "Installation for Windows completed successfully."
}

# Function to install for macOS
install_macos() {
    LIBWEBP_URL="https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.1.0-rc2-mac-10.15.tar.gz"

    # Define the base directory where you want to install the library
    BASE_DIR="/usr/local/libwebp"

    # The actual directory where libwebp contents will be
    INSTALL_DIR="${BASE_DIR}/libwebp-1.1.0-rc2-mac-10.15"

    # Create the base directory
    mkdir -p "$BASE_DIR"

    # Download and extract the libwebp binaries
    curl -L "$LIBWEBP_URL" | tar -xz -C "$BASE_DIR"

    # Add the lib and include paths to the environment variables
    export LIBRARY_PATH="${INSTALL_DIR}/lib:${LIBRARY_PATH}"
    export C_INCLUDE_PATH="${INSTALL_DIR}/include:${C_INCLUDE_PATH}"
    export CPLUS_INCLUDE_PATH="${INSTALL_DIR}/include:${CPLUS_INCLUDE_PATH}"
    export LD_LIBRARY_PATH="${INSTALL_DIR}/lib:${LD_LIBRARY_PATH}"
    export PATH="${INSTALL_DIR}/bin:${PATH}"

    # Check if libsharpyuv is included and add it to the linker flags
    if [ -f "${INSTALL_DIR}/lib/libsharpyuv.a" ]; then
        export CGO_LDFLAGS="-L${INSTALL_DIR}/lib -lwebp -lsharpyuv ${CGO_LDFLAGS}"
    else
        export CGO_LDFLAGS="-L${INSTALL_DIR}/lib -lwebp ${CGO_LDFLAGS}"
    fi

    # Create a pkg-config file if it does not exist
    PKG_CONFIG_PATH="${INSTALL_DIR}/lib/pkgconfig"
    mkdir -p "$PKG_CONFIG_PATH"
    PKG_CONFIG_FILE="$PKG_CONFIG_PATH/libwebp.pc"

    if [ ! -f "$PKG_CONFIG_FILE" ]; then
        echo "prefix=$INSTALL_DIR" > "$PKG_CONFIG_FILE"
        echo "exec_prefix=\${prefix}" >> "$PKG_CONFIG_FILE"
        echo "libdir=\${exec_prefix}/lib" >> "$PKG_CONFIG_FILE"
        echo "includedir=\${prefix}/include" >> "$PKG_CONFIG_FILE"
        echo "" >> "$PKG_CONFIG_FILE"
        echo "Name: libwebp" >> "$PKG_CONFIG_FILE"
        echo "Description: WebP library" >> "$PKG_CONFIG_FILE"
        echo "Version: 1.1.0-rc2" >> "$PKG_CONFIG_FILE"
        echo "Libs: -L\${libdir} -lwebp" >> "$PKG_CONFIG_FILE"
        if [ -f "${INSTALL_DIR}/lib/libsharpyuv.a" ]; then
            echo "Libs: -L\${libdir} -lsharpyuv" >> "$PKG_CONFIG_FILE"
        fi
        echo "Cflags: -I\${includedir}" >> "$PKG_CONFIG_FILE"
    fi

    # Set PKG_CONFIG_PATH to the directory containing 'libwebp.pc'
    export PKG_CONFIG_PATH="$PKG_CONFIG_PATH:${PKG_CONFIG_PATH}"

    echo "libwebp version 1.1.0-rc2 installed successfully in ${INSTALL_DIR}"

}

echo "building for HOST: $HOST"

# Check HOST environment variable
case "$HOST" in
    *darwin*)
        install_macos
        ;;
    *mingw*|*win*)
        install_windows
        ;;
    *linux*)
        install_linux
        ;;
    *)
        echo "Unsupported HOST: $HOST"
        exit 1
        ;;
esac
