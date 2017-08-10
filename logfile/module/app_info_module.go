package module

import (
	"time"
)

type User struct {
	Id          int `grom:"primary_key"`
	Name        string `grom:"not null;unique"`
	Age         int `grom:"not null"`
	Address     string
	Status      int	`grom:"not null"`
	Created     time.Time	`gorm:"column:created"`
	Modified    time.Time	`gorm:"column:modified"`
}

func (user *User) Save(trackId string, user User) error {
	res := Db_instance.Create(&user)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return SAVE_AFFECTED_ZERO_ERROR
	}

	return nil
}
func (user *User) Update(trackId string, user User) error {
	record := User{}
	updates := map[string]interface{}{
		"name": user.Name,
		"age": user.Age,
		"address": user.Address,
		"status": user.Status,
	}

	res := Db_instance.Model(&record).Where("id=?", user.Id).Update(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return UPDATE_AFFECTED_ZERO_ERROR
	}

	return nil
}
func (user *User) Delete(trackId string, id int) error {
	record := User{}

	res := Db_instance.Where("id=?", id).Delete(record)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return DELETE_AFFECTED_ZERO_ERROR
	}

	return nil
}

func (user *User) Query(trackId string, queryParams map[string]interface{}) ([]User, error) {
	records := [] User{}

	res := Db_instance.Where(queryParams).Find(&records)

	if res.Error != nil {
		return []User{}, res.Error
	}

	return records, nil
}
