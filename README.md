# Embedding Exporter

An openTelemetry exporter that sends data to OpenAI's Embedding service.


## Running

[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/droosma/opentelemetry-embedding-exporter)

The following commands should be available within the `devcontainer`:

- debug

  Starts a new debug session with the `dlv` debugger. This will start the exporter and wait for a debugger to attach.
- run

  Starts the exporter without a debugger attached.
- build

  Builds the exporter binary.

## Tips

- To exit debug mode, open new terminal and type `killall dlv`