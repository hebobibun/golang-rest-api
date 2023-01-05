package model

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name string `json:"name"`
	Desc string `json:"desc"`
	Qty int `json:"qty"`
	UserID uint `json:"userID"`
}

type ItemModel struct {
	DB *gorm.DB
}

func (im *ItemModel) Insert(newItem Item) (Item, error) {
	err := im.DB.Create(&newItem).Error
	if err != nil {
		return Item{}, err
	}
	return newItem, nil
}

func (im *ItemModel) GetAll(id int) ([]Item, error) {
	res := []Item{}
	if err := im.DB.Where("user_id = ?", id).Find(&res).Error; err != nil {
		log.Println("Get All item query error : ", err.Error())
		return nil, err
	}

	return res, nil
}

func (im *ItemModel) GetByID(id, idLogin int) (Item, error) {
	res := Item{}
	if err := im.DB.Where("id = ?", id).First(&res).Error; err != nil {
		log.Println("Get item By ID query error : ", err.Error())
		return Item{}, err
	}

	if int(res.UserID) != idLogin {
		log.Println("Unauthorized request")
		return Item{}, errors.New("Unauthorized request")
	} else {
		return res, nil
	}
}

func (im *ItemModel) Update(updatedItem Item, idUser int) (Item, error) {
	res := Item{}
	qry := im.DB.Where("id = ?", updatedItem.ID).First(&res)

	err := qry.Error
	if err != nil {
		log.Println("Select item query error : ", err.Error())
		return Item{}, nil
	}

	if int(res.UserID) != idUser {
		log.Println("Unauthorized request")
		return Item{}, errors.New("Unauthorized request")
	}

	qryUpdate := im.DB.Model(&Item{}).Where("id = ?", updatedItem.ID).Updates(&updatedItem)

	err = qryUpdate.Error

	if err != nil {
		log.Println("Update item query error : ", err.Error())
		return Item{}, errors.New("Unable to update the item")
	}

	return updatedItem, nil
}

func (im *ItemModel) Delete(itemID, id int) error {
	res := Item{}
	qry := im.DB.Where("id = ?", itemID).First(&res)

	err := qry.Error

	if err != nil {
		log.Println("Delete query error : ", err.Error())
		return errors.New("Unable to delete data")
	}

	if int(res.UserID) != id {
		log.Println("Unauthorized request")
		return errors.New("Unauthorized request")
	}

	qryDelete := im.DB.Delete(&Item{}, itemID)
	affRow := qryDelete.RowsAffected

	if affRow <= 0 {
		log.Println("No rows affected")
		return errors.New("Failed to delete, data not found")
	}

	return nil
}

