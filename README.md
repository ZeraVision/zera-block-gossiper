# Zera Block Gossiper

A development tool for testing block processing functionality for certain applications. This utility allows developers to fetch specific blocks from the ZV Indexer and send them to any gRPC endpoint, making it ideal for testing various logic.

## Features

- Fetch block details by block height from the ZV Indexer
- Send blocks to any gRPC endpoint for processing
- Configurable target address and port

## Prerequisites
- Access to a ZV Indexer API

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ZeraVision/zera-block-gossiper.git
   cd zera-block-gossiper
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Required: Your Zera indexer API key
INDEXER_AUTH=

# Required: Block height to fetch and send
BLOCK_HEIGHT=

# Required: gRPC server address
GRPC_ADDRESS=

# Optional: gRPC server port (default: 50051)
# GRPC_PORT=
```

## Usage

1. Set up your `.env` file with the required configuration
2. Run the application:
   ```bash
   go run main.go
   ```

The application will:
1. Fetch the specified block from the Zera indexer
2. Connect to the gRPC server at the specified address and port
3. Send the block using the `Broadcast` method

## Error Handling

The application will panic with descriptive error messages for:
- Missing required environment variables
- Invalid block height format
- Connection errors to the gRPC server
- Failed block retrieval from the indexer

## Development

### Building

To build the application:

```bash
go build -o zera-block-gossiper
```

### Testing

1. Ensure you have a test gRPC server running
2. Update the `.env` file with your test server details
3. Run the application with test parameters

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

