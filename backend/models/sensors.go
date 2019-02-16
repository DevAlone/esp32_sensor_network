package models

type SensorNode struct {
	MacAddress uint64 `sql:",pk" json:"mac_address"`
	Name       string `sql:",notnull" json:"name"`
}

/*type SensorDataBase struct {
}*/

type Sensor struct {
	Id                   uint64 `sql:",pk" json:"id" api:"ordering,filter"`
	SensorNodeMacAddress uint64 `sql:",notnull" json:"sensor_node_mac_address" api:"ordering,filter"`
	// type of sensor, say temperature
	Type string `sql:",notnull" json:"type"`
	Pin  uint8  `sql:",notnull" json:"pin"`
}

type SensorData struct {
	SensorId  uint64        `sql:",pk" json:"sensor_id" api:"ordering,filter"`
	Timestamp TimestampType `sql:",pk" json:"timestamp" api:"ordering,filter"`
	Value     float32       `sql:",notnull" json:"value"`
}

/*
type PikabuUser struct {
	PikabuId uint64 `sql:",pk" json:"pikabu_id" api:"ordering,filter"`

	Username            string        `sql:",notnull" gen_versions:"" json:"username" api:"ordering,filter"`
	Gender              string        `sql:",notnull" gen_versions:"" json:"gender" api:"ordering"`
	Rating              int32         `sql:",notnull" gen_versions:"" json:"rating" api:"ordering,filter"`
	NumberOfComments    int32         `sql:",notnull" gen_versions:"" json:"number_of_comments" api:"ordering,filter"`
	NumberOfSubscribers int32         `sql:",notnull" gen_versions:"" json:"number_of_subscribers" api:"ordering,filter"`
	NumberOfStories     int32         `sql:",notnull" gen_versions:"" json:"number_of_stories" api:"ordering"`
	NumberOfHotStories  int32         `sql:",notnull" gen_versions:"" json:"number_of_hot_stories" api:"ordering"`
	NumberOfPluses      int32         `sql:",notnull" gen_versions:"" json:"number_of_pluses" api:"ordering"`
	NumberOfMinuses     int32         `sql:",notnull" gen_versions:"" json:"number_of_minuses" api:"ordering"`
	SignupTimestamp     TimestampType `sql:",notnull" gen_versions:"" json:"signup_timestamp" api:"ordering"`
	AvatarURL           string        `sql:",notnull" gen_versions:"" json:"avatar_url" api:"ordering"`
	ApprovedText        string        `sql:",notnull" gen_versions:"" json:"approved_text" api:"ordering"`
	AwardIds            []uint64      `sql:",notnull,array" gen_versions:"" json:"award_ids" api:"ordering"`
	CommunityIds        []uint64      `sql:",notnull,array" gen_versions:"" json:"community_ids" api:"ordering"`
	BanHistoryItemIds   []uint64      `sql:",notnull,array" gen_versions:"" json:"ban_history_item_ids" api:"ordering"`
	BanEndTimestamp     TimestampType `sql:",notnull" gen_versions:"" json:"ban_end_timestamp" api:"ordering"`
	IsRatingHidden      bool          `sql:",notnull" gen_versions:"" json:"is_rating_hidden" api:"ordering"`
	IsBanned            bool          `sql:",notnull" gen_versions:"" json:"is_banned" api:"ordering"`
	IsPermanentlyBanned bool          `sql:",notnull" gen_versions:"" json:"is_permanently_banned" api:"ordering"`

	// ?
	// IsDeleted bool `sql:",notnull,default:false"`

	AddedTimestamp      TimestampType `sql:",notnull" json:"added_timestamp" api:"ordering"`
	LastUpdateTimestamp TimestampType `sql:",notnull" json:"last_update_timestamp" api:"ordering"`
	NextUpdateTimestamp TimestampType `sql:",notnull" json:"next_update_timestamp" api:"ordering"`
}
*/

func init() {
	for _, item := range []interface{}{
		&SensorNode{},
		&Sensor{},
		&SensorData{},
	} {
		Tables = append(Tables, item)
	}

	addIndex("sensors", "sensor_node_mac_address", "")
	addIndex("sensors", "type", "")
	addIndex("sensors", "pin", "")
	// addIndex("sensor_data", "timestamp", "")

	// addIndex("sensor_dht_data", "temperature", "")
	// addIndex("sensor_dht_data", "humidity", "")
	addUniqueIndex("sensors", []string{
		"sensor_node_mac_address",
		"type",
		"pin",
	}, "")
}
