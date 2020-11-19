package models

type Mappings struct {
	NamespaceId int64 `json:"namespace_id"`
	KeyId	int64 `json:"key_id"`
	ValExt string `json:"val_ext"`
	ValInt string `json:"val_int"`
	Payload string `json:"payload"`
} 