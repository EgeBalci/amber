FROM golang:1.20 as builder

RUN apt-get update && apt-get -y install \
    build-essential \    
    cmake \
    g++-multilib \
    gcc-multilib \
    git \
    libcapstone-dev \
    python3 \
    time
WORKDIR /root/
RUN git clone https://github.com/EgeBalci/keystone
RUN mkdir keystone/build
WORKDIR /root/keystone/build

RUN ../make-lib.sh
RUN cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF -DLLVM_TARGETS_TO_BUILD="AArch64;X86" -G "Unix Makefiles" ..
RUN make -j8
RUN make install && ldconfig

# RUN mkdir /root/amber
WORKDIR /root
RUN git clone https://github.com/egebalci/amber
WORKDIR /root/amber
RUN go build -trimpath -buildvcs=false -ldflags="-extldflags=-static -s -w" -o /root/bin/amber  main.go

FROM scratch
COPY --from=builder /root/bin/amber /amber
ENTRYPOINT ["/amber"]
