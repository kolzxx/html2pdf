FROM 372349696520.dkr.ecr.us-west-2.amazonaws.com/docker-hub/library/golang:1.23.4-bookworm as build

ARG version

RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y libtool musl-dev
RUN apt-get -y install subversion g++ make gcc

WORKDIR /app

COPY . .

RUN make build VERSION=$version


FROM 372349696520.dkr.ecr.us-west-2.amazonaws.com/docker-hub/chromedp/headless-shell:127.0.6533.120 as navigator

WORKDIR /app

ENV GIN_MODE=release
ENV SWAGGER_ENABLED=false

COPY --from=build /app/server .

RUN mkdir -p /etc/ssl/certs/

RUN apt-get update && apt-get install -y ca-certificates dumb-init chromium && apt-get clean

CMD chromium-browser --password-store=basic --disable-gpu --disable-web-security --disable-background-timer-throttling --disable-popup-blocking --disable-prompt-on-repost --disable-renderer-backgrounding --metrics-recording-only --force-color-profile=srgb --no-first-run --hide-scrollbars --disable-background-networking --disable-default-apps --no-default-browser-check --mute-audio --disable-breakpad --enable-automation --password-store=basic --disable-client-side-phishing-detection --disable-dev-shm-usage --disable-extensions --safebrowsing-disable-auto-update --use-mock-keychain --headless --disable-hang-monitor --disable-sync --enable-features=NetworkService,NetworkServiceInProcess --disable-backgrounding-occluded-windows --disable-features=site-per-process,Translate,BlinkGenPropertyTrees --disable-ipc-flooding-protection --remote-debugging-port=9222 about:blank --no-sandbox

ENV PATH /headless-shell:$PATH

RUN echo "/app/server" >> entrypoint.sh \
    && chmod a+x entrypoint.sh

ENTRYPOINT [ "/bin/sh", "-c", "./entrypoint.sh" ]