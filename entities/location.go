package entities

type Province struct {
	ProvID uint `gorm:"primarykey"`
	ProvName string
	LocationID string `gorm:"column:locationid"`
	Status int
}

type City struct {
	CityID uint `gorm:"primarykey"`
	CityName string
	ProvID uint
	Province Province `gorm:"foreignKey:ProvID;references:ProvID"`
}

type District struct {
	DisID uint `gorm:"primarykey"`
	DisName string
	CityID uint
	City City `gorm:"foreignKey:CityID;references:CityID"`
}


type ProvinceResponse struct {
	ProvID int
	ProvName string
}
type CityResponse struct {
	CityID int
	CityName string
	ProvID int
}
type DistrictResponse struct {
	DisID int
	DisName string
	CityID int
}