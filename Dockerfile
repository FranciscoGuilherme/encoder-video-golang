FROM golang:1.27.3-alpine
ENV CGO_ENABLED=1 \
    PATH="$PATH:/bin/bash" \
    BENTO4_BIN="/opt/bento4/bin" \
    PATH="$PATH:/opt/bento4/bin"

# FFMPEG
RUN apk add --update ffmpeg bash make build-base

# -------------------------
# -----[Install Bento]-----
# -------------------------

WORKDIR /tmp/bento4
ENV BENTO4_BASE_URL="https://www.bok.net/Bento4/binaries" \
    BENTO4_VERSION="1-5-0-615" \
    BENTO4_CHECKSUM="5e5567d474f18be72db4745983e7ec149768e5a7e0053babfe30bf66a437e7a1" \
    BENTO4_TARGET=".x86_64-unknown-linux" \
    BENTO4_PATH="/opt/bento4" \
    BENTO4_TYPE="SDK"

# -------------------------------------
# -----[download and unzip bento4]-----
# -------------------------------------

RUN apk add --no-cache unzip gcompat libstdc++ libgcc && \
    wget -q ${BENTO4_BASE_URL}/Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip && \
    echo "${BENTO4_CHECKSUM}  Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip" | sha256sum -c - && \
    mkdir -p ${BENTO4_PATH} && \
    unzip -q Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip \
    "Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}/bin/*" -d ${BENTO4_PATH} && \
    mv ${BENTO4_PATH}/Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}/bin ${BENTO4_PATH}/bin && \
    rm -rf ${BENTO4_PATH}/Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET} \
    Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip && \
    apk del unzip

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]
