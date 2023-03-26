FROM ubuntu:22.04 as build

RUN apt-get update \
    && apt-get install -y --no-install-recommends libsdl2-dev libsdl2-gfx-dev libsdl2-image-dev libsdl2-ttf-dev libsdl2-gfx-dev fonts-dejavu curl ca-certificates git gcc \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

# Install Go
ARG GO_VERSION=1.18.10
RUN curl -L https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz --output go${GO_VERSION}.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /csgoverview

COPY . .

RUN go build

# Download map images
RUN chmod 755 overviews.sh && ./overviews.sh


FROM ubuntu:22.04 as runtime

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends dumb-init x11vnc xvfb novnc libsdl2-dev libsdl2-gfx-dev libsdl2-image-dev libsdl2-ttf-dev libsdl2-gfx-dev fonts-dejavu zenity \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -rm -d /home/ubuntu -s /bin/bash -g root -G sudo -u 1001 ubuntu
USER ubuntu
WORKDIR /home/ubuntu

COPY --from=build /csgoverview/csgoverview csgoverview/

# Copy map overviews into the correct location
RUN mkdir -p /home/ubuntu/.local/share/csgoverview/assets/maps
COPY --from=build /csgoverview/overviews/*.jpg /home/ubuntu/.local/share/csgoverview/assets/maps/

# Script that will start a x11 server (xvfb), vnc server (x11vnc), and expose vnc over htmml (novnc) 
RUN echo "#!/bin/bash\n\
x11vnc -create -env FD_PROG='/home/ubuntu/csgoverview/csgoverview' -env X11VNC_FINDDISPLAY_ALWAYS_FAILS=1 -env X11VNC_CREATE_GEOM=${1:-1624x1024x16} -nopw -loop & \n\
/usr/share/novnc/utils/launch.sh --listen 6081 --vnc localhost:5900" > ~/run.sh
RUN chmod 755 ~/run.sh

# Create volume mount where users can put their demos
RUN mkdir -p /home/ubuntu/demos
VOLUME /home/ubuntu/demos

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/home/ubuntu/run.sh"]