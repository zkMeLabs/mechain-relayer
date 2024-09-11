package main

import (
	"fmt"
	"os"
	"text/template"
)

type PortConfig struct {
	AddressPort int
	P2PPort     int
	GRPCPort    int
	GRPCWebPort int
	RPCPort     int
	EVMRPCPort  int
	EVMWSPort   int
}
type NodeConfig struct {
	NodeIndex int
	PortConfig
}

type ComposeConfig struct {
	Nodes           []NodeConfig
	Image           string
	ProjectBasePath string
	BasePorts       PortConfig
}

const dockerComposeTemplate = `
services:
  init:
    container_name: init-relayer
    image: "{{$.Image}}"
    volumes:
      - "{{$.ProjectBasePath}}/deployment/dockerup:/workspace/deployment/dockerup:Z"
      - "local-env:/workspace/deployment/dockerup/.local"
    working_dir: "/workspace/deployment/dockerup"
    command: >
      bash -c "
      rm -f init_done &&
      bash localup.sh config 4 && 
      touch init_done && 
      sleep infinity
      "
    healthcheck:
      test: ["CMD-SHELL", "test -f /workspace/deployment/dockerup/init_done && echo 'OK' || exit 1"]
      interval: 10s
      retries: 5
    restart: "on-failure"
{{- range .Nodes }}
  rnode{{.NodeIndex}}:
    container_name: mechain-relayer-{{.NodeIndex}}
    depends_on:
      init:
        condition: service_healthy
    image: "{{$.Image}}"
    volumes:
      - "local-env:/app"
    command: >
      greenfield-relayer run --config-type local --config-path /app/relayer{{.NodeIndex}}/config.json --log_dir json
{{- end }}
`

func main() {
	bp := PortConfig{
		AddressPort: 28750,
		P2PPort:     27750,
		GRPCPort:    9090,
		GRPCWebPort: 1317,
		RPCPort:     26657,
		EVMRPCPort:  8545,
		EVMWSPort:   8546,
	}

	numNodes := 4

	var nodes []NodeConfig
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, NodeConfig{
			NodeIndex: i,
			PortConfig: PortConfig{
				AddressPort: bp.AddressPort + i,
				P2PPort:     bp.P2PPort + i,
				GRPCPort:    bp.GRPCPort + i,
				GRPCWebPort: bp.GRPCWebPort + i,
				RPCPort:     bp.RPCPort + i,
				EVMRPCPort:  bp.EVMRPCPort + i*2,
				EVMWSPort:   bp.EVMWSPort + i*2,
			},
		})
	}

	config := ComposeConfig{
		Nodes:           nodes,
		Image:           "zkmelabs/mechain-relayer",
		ProjectBasePath: ".",
		BasePorts:       bp,
	}

	file, err := os.Create("docker-compose.yml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	tpl, err := template.New("docker-compose").Parse(dockerComposeTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tpl.Execute(file, config)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Docker Compose file generated successfully!")
}
