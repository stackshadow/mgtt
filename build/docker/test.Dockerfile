FROM archlinux/base as base


RUN pacman --noconfirm -Syu && \
    pacman --noconfirm -Syy glibc archlinux-keyring gawk && \
    pacman --noconfirm -Scc

RUN pacman --noconfirm -Syy busybox make go && \
    pacman --noconfirm -Scc

RUN pacman --noconfirm -Syy grep python python-pip && \
    pacman --noconfirm -Scc

RUN pip install --user anybadge
