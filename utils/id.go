package utils

import (
    "github.com/bwmarrin/snowflake"
    "os"
    "strconv"
)

type Node struct {
    node *snowflake.Node
}

var local *Node

func init() {
    nodeId := 1
    if str := os.Getenv("NodeId"); str != "" {
        nodeId, _ = strconv.Atoi(str)
        if nodeId <= 0 {
            nodeId = 1
        }
    }
    local, _ = NewNode(nodeId)
}

func GenerateId() string {
    return local.Generate()
}

func NewNode(id int) (*Node, error) {
    node, err := snowflake.NewNode(int64(id % 1024))
    return &Node{node: node}, err
}

func (n *Node) Generate() string {
    return Base62Encode(n.node.Generate().Int64())
}
