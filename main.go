package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	zera_protobuf "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/prototext"
)

// GetBlockDetails retrieves raw block details from Zera indexer
func GetBlockDetails(blockHeight int64) (*zera_protobuf.Block, error) {
	url := fmt.Sprintf("https://indexer.zera.vision/store?requestType=getBlockDetailsRaw&blockHeight=%d", blockHeight)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	authToken := os.Getenv("INDEXER_AUTH")
	if authToken == "" {
		return nil, fmt.Errorf("INDEXER_AUTH environment variable not set")
	}

	req.Header.Set("Target", "indexer")
	req.Header.Set("authorization", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// The response is a JSON string containing the protobuf text format
	var responseStr string
	if err := json.Unmarshal(body, &responseStr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	// Now parse the protobuf text format
	block := &zera_protobuf.Block{}
	if err := prototext.Unmarshal([]byte(responseStr), block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf text: %v", err)
	}

	return block, nil
}

func init() {
	godotenv.Load(".env")
}

// sendBlockViaGRPC sends a block to the specified gRPC server using the Broadcast method
func sendBlockViaGRPC(block *zera_protobuf.Block, address string, port int) error {
	grpcAddr := fmt.Sprintf("%s:%d", address, port)

	// Set up a connection to the server
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create a new client (leverage the validator service configs to essentiallygossip blocks)
	client := zera_protobuf.NewValidatorServiceClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the block using the Broadcast method
	_, err = client.Broadcast(ctx, block)
	return err
}

func main() {
	// Get block height from environment
	blockHeightStr := os.Getenv("BLOCK_HEIGHT")
	if blockHeightStr == "" {
		panic("Error: BLOCK_HEIGHT environment variable not set")
	}

	blockHeight, err := strconv.ParseInt(blockHeightStr, 10, 64)
	if err != nil {
		panic("Error: invalid BLOCK_HEIGHT")
	}

	// Get gRPC target from environment or use default
	grpcAddress := os.Getenv("GRPC_ADDRESS")

	grpcPort := 50051
	if portStr := os.Getenv("GRPC_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			grpcPort = p
		}
	}

	// Get block details
	block, err := GetBlockDetails(int64(blockHeight))
	if err != nil {
		fmt.Printf("Error getting block details: %v\n", err)
		return
	}

	// Send block via gRPC
	if err := sendBlockViaGRPC(block, grpcAddress, grpcPort); err != nil {
		fmt.Printf("Error sending block via gRPC: %v\n", err)
		return
	}

	fmt.Printf("Successfully sent block %d to %s:%d\n", blockHeight, grpcAddress, grpcPort)
}
