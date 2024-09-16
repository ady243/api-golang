package models

import "time"

type Users struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    Username      string    `json:"username" gorm:"unique"`
	Email         string    `json:"email" gorm:"unique"`
    PasswordHash  string    `json:"password_hash"`
    CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
    

    ProfilePhoto  string    `json:"profile_photo"`  
    FavoriteSport string    `json:"favorite_sport"`
    Location      string    `json:"location"`      
    SkillLevel    string    `json:"skill_level"`    
    Bio           string    `json:"bio"`            
    
    UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
}
