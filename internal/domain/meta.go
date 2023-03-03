package domain

type MetaProvider interface {
	Get() *BusinessMeta
	Set(meta BusinessMeta)
}

type BusinessMeta struct {
	DeliveryPunishmentThreshold int64 `json:"deliveryPunishmentThreshold" bson:"deliveryPunishmentThreshold"`
	DeliveryPunishmentValue     int64 `json:"deliveryPunishmentValue" bson:"deliveryPunishmentValue"`
}
