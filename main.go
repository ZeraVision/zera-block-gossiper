package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	zera_protobuf "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/joho/godotenv"
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

func main() {
	// Example usage:
	blockHeight := 2763
	block, err := GetBlockDetails(int64(blockHeight))
	if err != nil {
		fmt.Printf("Error getting block details: %v\n", err)
		return
	}

	fmt.Printf("Block %d details: %+v\n", blockHeight, block)
}
