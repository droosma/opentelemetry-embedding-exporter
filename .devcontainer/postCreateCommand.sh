GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@latest
go install github.com/go-delve/delve/cmd/dlv@latest
echo "alias debug='dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient --log exec /tmp/dist/otelcol-custom -- --config config.yaml'" >> ~/.bashrc
echo "alias run='/tmp/dist/otelcol-custom --config=config.yaml'" >> ~/.bashrc
echo "alias build='builder --config=otelcol-builder.yaml'" >> ~/.bashrc