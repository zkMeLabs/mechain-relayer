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
	MySQLImage      string
	ProjectBasePath string
	BasePorts       PortConfig
}

const dockerComposeTemplate = `
services:
  relayer-mysql:
    container_name: relayer-mysql
    image: {{.MySQLImage}}
    networks:
      - mechain-network
    volumes:
      - db-data:/var/lib/mysql
    environment:
      MYSQL_ROOT_PASSWORD: mechain
      MYSQL_DATABASE: greenfield_relayer
    ports:
      - "3307:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
  init-relayer:
    container_name: init-relayer
    image: "{{$.Image}}"
    networks:
      - mechain-network
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
      relayer-mysql:
        condition: service_healthy
      init:
        condition: service_healthy
    image: "{{$.Image}}"
    networks:
      - mechain-network
    volumes:
      - "local-env:/app"
    command: >
      greenfield-relayer run --config-type local --config-path /app/relayer{{.NodeIndex}}/config.json --log_dir json
{{- end }}
volumes:
  db-data:
  local-env:
networks:
  mechain-network:
    external: true
`

func main() {
	bp := PortConfig{	}

	numNodes := 4

	var nodes []NodeConfig
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, NodeConfig{
			NodeIndex: i,
			PortConfig: PortConfig{
		
			},
		})
	}

	config := ComposeConfig{
		Nodes:           nodes,
		Image:           "zkmelabs/mechain-relayer",
		MySQLImage:      "mysql:8",
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
