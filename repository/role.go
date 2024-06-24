package repository

import (
	"github.com/ulbithebest/todolist-be/model"
	"gorm.io/gorm"
)

func GetAllRole(db *gorm.DB) ([]model.Roles, error) { // Mengambil semua data Role dari database
	var role []model.Roles
	if err := db.Find(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func GetRoleById(db *gorm.DB, id string) (model.Roles, error) { // Mengambil data Role berdasarkan ID dari database
	var role model.Roles
	if err := db.First(&role, id).Error; err != nil {
		return role, err
	}
	return role, nil
}

func InsertRole(db *gorm.DB, role model.Roles) error { // Insert data role ke dalam database
	if err := db.Create(&role).Error; err != nil {
		return err
	}
	return nil
}

func UpdateRole(db *gorm.DB, id string, updatedRole model.Roles) error { // Memperbarui data role dalam database berdasarkan ID
	if err := db.Model(&model.Roles{}).Where("id_role = ?", id).Updates(updatedRole).Error; err != nil {
		return err
	}
	return nil
}

func DeleteRole(db *gorm.DB, id string) error { // Menghapus data role dari database berdasarkan ID
	if err := db.Delete(&model.Roles{}, id).Error; err != nil {
		return err
	}
	return nil
}
