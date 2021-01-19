############################################## ARCH ##############################################
FROM archlinux/base as base


RUN pacman --noconfirm -Syy glibc busybox go && \
    pacman --noconfirm -Scc

# basic root dir
RUN \
    mkdir -pv /newroot/tmp && \
    mkdir -pv /newroot/lib && \
    mkdir -pv /newroot/lib64 && \
    \
    cp -pv /usr/lib/libdl* /usr/lib/libc-*.so* /usr/lib/libc.so* /usr/lib/libpthread*.so* /newroot/lib && \
    cp -pv /usr/lib64/ld-*.so /usr/lib64/ld-linux-x86-64.so.2 /newroot/lib64


# locales
RUN \
    mkdir -pv /newroot/etc && \
    echo LANG=de_DE.UTF-8 > /newroot/etc/locale.conf && \
    echo KEYMAP=de-latin1-nodeadkeys > /newroot/etc/vconsole.conf && \
    echo "en_US.UTF-8 UTF-8" > /newroot/etc/locale.gen && \
    echo "de_DE.UTF-8 UTF-8" >> /newroot/etc/locale.gen && \
    \
    mkdir -pv /newroot/usr/share/i18n/locales && \
    cp -v /usr/share/i18n/locales/de_DE /newroot/usr/share/i18n/locales/de_DE && \
    cp -v /usr/share/i18n/locales/en_US /newroot/usr/share/i18n/locales/en_US && \
    \
    mkdir -pv /newroot/usr/share/locale && \
    cp -vR /usr/share/locale/de /newroot/usr/share/locale/de && \
    cp -vR /usr/share/locale/en_GB /newroot/usr/share/locale/en_GB && \
    \
    mkdir -pv /newroot/usr/share/i18n/charmaps && \
    cp -v /usr/share/i18n/charmaps/UTF-8.gz /newroot/usr/share/i18n/charmaps/UTF-8.gz && \
    \
    locale-gen && \
    mkdir -pv /newroot/usr/lib/locale && \
    cp -v /usr/lib/locale/locale-archive /newroot/usr/lib/locale/locale-archive

# zoneinfo
RUN mkdir -pv /newroot/usr/share/zoneinfo && \
    cp /usr/share/zoneinfo/UTC /newroot/usr/share/zoneinfo && \
    \
    mkdir -pv /newroot/usr/share/zoneinfo/right && \
    cp /usr/share/zoneinfo/right/UTC /newroot/usr/share/zoneinfo/right

# hosts-file
RUN echo "127.0.0.1 localhost" > /newroot/etc/hosts

# This is that binarys can determine the username from an id
# for example whoami or id ( without this libs, no binary can get the username from an uid )
# libnss_dns -> DNS resolving hostnames to ips
RUN mkdir -pv /newroot/usr/lib && \
    cp -vp \
    /usr/lib/libnss_compat* \
    /usr/lib/libnss_db* \
    /usr/lib/libnss_dns* \
    /usr/lib/libnss_files* \
    /usr/lib/libnss_myhostname* \
    /usr/lib/libnss_resolve* \
    /usr/lib/libnsl* \
    /newroot/usr/lib

# default user
RUN echo "root:x:0:0::/root:/bin/bash" >> /newroot/etc/passwd && \
    echo "nobody:x:65534:65534:Nobody:/:/sbin/nologin" >> /newroot/etc/passwd




######################################################## healthcheck ########################################################
FROM archlinux/base as healthcheck

RUN pacman --noconfirm -Syy glibc busybox go && \
    pacman --noconfirm -Scc

COPY cmd/healthcheck/main.go cmd/healthcheck/go.mod cmd/healthcheck/go.sum /go/src/
WORKDIR /go/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o healthcheck . && \
    ./healthcheck -h



######################################################## mgtt ########################################################
FROM archlinux/base as mgtt

RUN pacman --noconfirm -Syy glibc busybox go && \
    pacman --noconfirm -Scc

WORKDIR /

# download modules - this is cached when go.mod and go.sum not changes
COPY go.mod go.sum /
#RUN mkdir /vendor
#COPY mgtt/vendor/modules.txt /vendor/
RUN go mod download 

# Copy the rest
COPY cmd/  /cmd
COPY internal/ /internal

RUN ls -ahl /

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o mgtt gitlab.com/mgtt/cmd/mgtt && \
    ./mgtt -h

RUN mkdir -p /data && \
    chown 65534:65534 /data

######################################################## The finale image ########################################################
FROM scratch as prod

ENV LD_LIBRARY_PATH /usr/local/lib:/usr/lib:/lib:/usr/lib64:/lib64

# copy new root
COPY --from=base /newroot/etc /etc
COPY --from=base /newroot/usr /usr
COPY --from=base /newroot/lib /lib
COPY --from=base /newroot/lib64 /lib64
COPY --from=base --chown=nobody:0 /newroot/tmp /tmp

# healthceck
COPY --from=healthcheck --chown=nobody:0 /go/src/healthcheck /healthcheck

# mgtt
COPY --from=mgtt --chown=nobody:0 /mgtt /mgtt
COPY --from=mgtt --chown=nobody:0 /data /data

# for debugging
#COPY --from=base --chown=nobody:0 /bin/sh /bin/sh
#COPY --from=base --chown=nobody:0 /usr/lib/libreadline* /usr/lib/
#COPY --from=base --chown=nobody:0 /usr/lib/libncurses* /usr/lib/

USER nobody
CMD ["/mgtt", "--debug", "--terminal", "serve", "--db-filename=/data/messages.db", "--ca-file=/data/ca.crt", "--cert-file=/data/server.crt", "--key-file=/data/server.key", "--config-path=/data/" ]

