# TFTP-HTTP Proxy

This tool serves as a proxy for TFTP requests, forwarding them to an HTTP endpoint. It is particularly useful in PXE (Preboot Execution Environment) environments, where the client uses TFTP to download files but the server hosting the files only supports HTTP.

## Usage

$ tftp-http-proxy [options]

Copy code

### Options

- `--listen`: IP address to listen on for TFTP requests (default: "0.0.0.0")
- `--port`: The port to listen on for TFTP requests (default: 69) 
- `--url`: URL to forward TFTP requests to
- `--log-level`: Sets the default log level (default: "info")

### Example

$ tftp-http-proxy --url http://example.com/tftp/

This will start the proxy and listen for TFTP requests on all available network interfaces on port 69. Any requests received will be forwarded to `http://example.com/tftp/`

## Installation

This tool is a single go binary, no particular installation instruction.

## Notes

- The tool only supports a subset of the TFTP protocol and is only intended to be used in PXE environments.
- It's only a basic example and you should consider security, performance and scalability issues before using it in production.

## Contributing

We welcome contributions to this project. Please submit a pull request with your changes.

## Licensing

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE-2.0.txt) file for details.