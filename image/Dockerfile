FROM garybowers/spider-base:1.0 

# Install tools
RUN apt-get update -y && \
    apt-get install -y tig vim less dirmngr gnupg unzip zip g++ zsh git tig mc htop w3m && \
    rm -rf /var/cache/apt/* && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*
 
# Install golang
ARG GO_VER=1.14
ARG GO_PATH=/workspace/go
ENV GO_ROOT=/usr/local/go
ENV PATH $PATH:/usr/local/go/bin
ENV PATH $PATH:${GOPATH}/bin

RUN curl -sS https://storage.googleapis.com/golang/go$GO_VER.linux-amd64.tar.gz | tar -C /usr/local -xzf -

# Install Java
RUN apt-get update && \
    apt-get -y install openjdk-11-jre openjdk-11-jdk maven gradle && \ 
    rm -rf /var/cache/apt/* && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*

#C/C++
# public LLVM PPA, development version of LLVM
RUN apt-get update -y && apt-get install -y zlib1g-dev valkyrie valgrind && \
    wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add - && \
    echo "deb http://apt.llvm.org/buster/ llvm-toolchain-buster-10 main" > /etc/apt/sources.list.d/llvm.list && \
    echo "deb-src http://apt.llvm.org/buster/ llvm-toolchain-buster-10 main" >> /etc/apt/sources.list.d/llvm.list && \
    apt-get update -y && \
    apt-get install -y clang* && \ 
    rm -rf /var/cache/apt/* && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*

# Install latest stable CMake
RUN wget "https://cmake.org/files/v3.17/cmake-3.17.0-Linux-x86_64.sh" && \
    chmod a+x cmake-3.17.0-Linux-x86_64.sh && \
    ./cmake-3.17.0-Linux-x86_64.sh --prefix=/usr/ --skip-license && \
    rm cmake-3.17.0-Linux-x86_64.sh

# Install Bazel
ARG BAZEL_VER=3.0.0
RUN wget https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VER}/bazel_${BAZEL_VER}-linux-x86_64.deb && \
    dpkg -i bazel_${BAZEL_VER}-linux-x86_64.deb && \
    rm bazel_${BAZEL_VER}-linux-x86_64.deb && \
    wget https://github.com/bazelbuild/buildtools/releases/download/${BAZEL_VER}/buildifier && mv buildifier /usr/bin/ && chmod +x /usr/bin/buildifier && \
    wget https://github.com/bazelbuild/buildtools/releases/download/${BAZEL_VER}/buildozer && mv buildozer /usr/bin/ && chmod +x /usr/bin/buildozer && \
    wget https://github.com/bazelbuild/buildtools/releases/download/${BAZEL_VER}/unused_deps && mv unused_deps /usr/bin/ && chmod +x /usr/bin/unused_deps && \
    echo "TEST_TMPDIR=/workspace/.cache" >> /etc/environment

# Install Terraform
ARG TF_VER=0.12.25
RUN wget https://releases.hashicorp.com/terraform/${TF_VER}/terraform_${TF_VER}_linux_amd64.zip && \
    unzip terraform_${TF_VER}_linux_amd64.zip && mv terraform /usr/local/bin && \
    rm -f terraform_${TF_VER}_linux_amd64.zip

# Install gcloud sdk
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt-get update && apt-get install -y google-cloud-sdk kubectl

VOLUME /home/coder
  
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/init.sh"]
