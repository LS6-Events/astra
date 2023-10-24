package types

type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
)

type Operation struct {
	Type       OperationType `json:"operation" yaml:"operation"`
	EntityType string        `json:"entity_type" yaml:"entityType"`
	EntityID   int           `json:"entity_id" yaml:"entityID"`
}
