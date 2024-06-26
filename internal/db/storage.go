package db

import (
	"fmt"
	"log"

	"Avito-Project/internal/config"
	"Avito-Project/internal/models"
	"database/sql"
	_ "github.com/lib/pq"
)

type DataBase struct {
	DB *sql.DB
}

func (db *DataBase) GetStorage(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	var err error
	db.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}

func (db *DataBase) Stop() error {
	if db.DB != nil {
		err := db.DB.Close()
		{
			if err != nil {
				log.Fatalf("Failed to closed database: %v", err)
				return err
			}
		}
	}
	return nil
}

func (db *DataBase) GetUserByToken(token string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, access_levels, created_at, updated_at, token FROM Users WHERE token = $1"
	row := db.DB.QueryRow(query, token)

	err := row.Scan(&user.Id, &user.Name, &user.AccessLevels, &user.CreatedAt, &user.UpdatedAt, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to scan row: %v", err)
		return nil, err
	}

	return &user, nil
}

func (db *DataBase) GetBanner(id int) (*models.Banner, error) {
	var banner models.Banner
	query := "SELECT id, title, text, url, created_at, updated_at, owner_id, f_id FROM Banners WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	err := row.Scan(&banner.Id, &banner.Title, &banner.Text, &banner.Url, &banner.CreatedAt, &banner.UpdatedAt, &banner.OwnerId, &banner.FId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Printf("Failed to scan row: %v\n", err)
		return nil, err
	}
	return &banner, nil
}

func (db *DataBase) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, access_levels, created_at, updated_at, token FROM Users WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	err := row.Scan(&user.Id, &user.Name, &user.AccessLevels, &user.CreatedAt, &user.UpdatedAt, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to scan row: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func (db *DataBase) GetBannerByTagID(tag int) ([]models.Banner, error) {
	var banners []models.Banner

	query := `
	SELECT b.id, b.title, b.text, b.url, b.created_at, b.updated_at, b.owner_id, b.f_id
   FROM banners b
   JOIN tags t ON b.id = t.banner_id
   WHERE t.tag = $1`

	rows, err := db.DB.Query(query, tag)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(&banner.Id, &banner.Title, &banner.Text, &banner.Url, &banner.CreatedAt, &banner.UpdatedAt, &banner.OwnerId, &banner.FId)
		if err != nil {
			log.Printf("Failed to scan rows: %v\n", err)
			return nil, err
		}

		tags, err := db.GetTagByBanner(banner.Id)
		if err != nil {
			log.Printf("Failed to get tags: %v", err)
		}

		banner.Tags = tags
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %v\n", err)
		return nil, err
	}

	return banners, nil
}

func (db *DataBase) GetTagByBanner(bannerId uint) ([]int, error) {
	tags := make([]int, 0)

	tagsQuery := "SELECT tag FROM Tags WHERE banner_id = $1"
	tagsRows, err := db.DB.Query(tagsQuery, bannerId)
	if err != nil {
		log.Printf("Failed to execute tags query: %v\n", err)
		return nil, err
	}
	defer tagsRows.Close()

	for tagsRows.Next() {
		var tag int
		err := tagsRows.Scan(&tag)
		if err != nil {
			log.Printf("Failed to scan tag row: %v\n", err)
			return nil, err
		}

		tags = append(tags, tag)

	}
	if err := tagsRows.Err(); err != nil {
		log.Printf("Tag row iteration error: %v\n", err)
		return nil, err
	}

	return tags, nil
}

func (db *DataBase) GetBannerByFID(fID int) ([]models.Banner, error) {
	var banners []models.Banner
	query := "SELECT id, title, text, url, created_at, updated_at, owner_id, f_id FROM Banners WHERE f_id = $1"
	rows, err := db.DB.Query(query, fID)
	if err != nil {
		log.Fatalf("Failed to execute query: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(&banner.Id, &banner.Title, &banner.Text, &banner.Url, &banner.CreatedAt, &banner.UpdatedAt, &banner.OwnerId, &banner.FId)
		if err != nil {
			log.Printf("Failed to scan rows")
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed iteration rows: %v\n", err)
		return nil, err
	}

	return banners, nil
}

func (db *DataBase) GetAllBanners() ([]models.Banner, error) {
	var banners []models.Banner

	query := "SELECT id, title, text, url, created_at, updated_at, owner_id, f_id FROM Banners"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(&banner.Id, &banner.Title, &banner.Text, &banner.Url, &banner.CreatedAt, &banner.UpdatedAt, &banner.OwnerId, &banner.FId)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed iteration rows: %v\n", err)
		return nil, err
	}

	return banners, nil
}

func (db *DataBase) GetAllUsers() ([]models.User, error) {
	var users []models.User

	query := "SELECT id, name, access_levels, created_at, updated_at, token FROM Users"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.AccessLevels, &user.CreatedAt, &user.UpdatedAt, &user.Token)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Failed iteration rows: %v\n", err)
		return nil, err
	}

	return users, nil
}

func (db *DataBase) AddUser(user *models.User) error {
	query := "INSERT INTO Users (name, access_levels,token) VALUES ($1, $2, $3) RETURNING id"
	err := db.DB.QueryRow(query, user.Name, user.AccessLevels, user.Token).Scan(&user.Id)
	if err != nil {
		log.Printf("Failed to add user: %v", err)
		return err
	}
	return nil
}

func (db *DataBase) DeleteUser(userId int) error {
	query := "DELETE FROM Users WHERE id = $1"
	_, err := db.DB.Exec(query, userId)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return err
	}
	return nil
}

func (db *DataBase) AddBanner(banner *models.Banner) error {
	query := "INSERT INTO Banners (title, text, url, owner_id, f_id) VALUES ($1, $2,$3, $4, $5) RETURNING id"
	err := db.DB.QueryRow(query, banner.Title, banner.Text, banner.Url, banner.OwnerId, banner.FId).Scan(&banner.Id)
	if err != nil {
		log.Printf("Failed to add banner: %v", err)
		return err
	}
	return nil
}

func (db *DataBase) DeleteBanner(bannerId int) error {
	query := "DELETE FROM Banners WHERE id = $1"
	_, err := db.DB.Exec(query, bannerId)
	if err != nil {
		log.Printf("Failed to delete banner: %v", err)
		return err
	}
	return nil
}

func (db *DataBase) AddAccessLevel(level *models.AccessLevel) error {
	query := "INSERT INTO Access_levels (level, job_title) VALUES ($1, $2)"
	err := db.DB.QueryRow(query, level.Level, level.JobTitle)
	if err != nil {
		log.Printf("Failed to add access level: %v", err)
	}
	return nil
}

func (db *DataBase) GetUsersPaginated(page, size int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * size

	query := `
	SELECT id, name, access_levels, created_at, updated_at, token 
	FROM Users 
	LIMIT $1
	OFFSET $2`

	rows, err := db.DB.Query(query, size, offset)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.AccessLevels, &user.CreatedAt, &user.UpdatedAt, &user.Token)
		if err != nil {
			log.Printf("Failed to scan rows: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed iteration rows: %v\n", err)
		return nil, err
	}

	return users, nil
}

func (db *DataBase) GetBannersPaginated(page, size int) ([]models.Banner, error) {
	var banners []models.Banner
	offset := (page - 1) * size

	query := `
	SELECT id, title, text, url, created_at, updated_at, owner_id, f_id 
	FROM Banners 
	LIMIT $1 
	OFFSET $2`

	rows, err := db.DB.Query(query, size, offset)
	if err != nil {
		log.Fatalf("Failed to execute query: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var banner models.Banner

		err := rows.Scan(&banner.Id, &banner.Title, &banner.Text, &banner.Url, &banner.CreatedAt, &banner.UpdatedAt, &banner.OwnerId, &banner.FId)
		if err != nil {
			log.Fatalf("Failed to rows scan: %v\n", err)
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iteration rows: %v\n", err)
		return nil, err
	}

	return banners, nil
}

func (db *DataBase) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, name, access_levels, created_at, updated_at, token, password 
	FROM Users
	WHERE name = $1 
	AND
	password = $2`

	err := db.DB.QueryRow(query, username, password).Scan(&user.Id, &user.Name, &user.AccessLevels, &user.CreatedAt, &user.UpdatedAt, &user.Token, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to authenticate user: %v\n ", err)
		return nil, err
	}

	return &user, err
}

func (db *DataBase) UpdateUser(user *models.User) error {
	query := `
	UPDATE Users 
	SET name=$1, access_levels=$2, updated_at=$3, token=$4, password=$5 
	WHERE id=$6
	`

	_, err := db.DB.Exec(query, user.Name, user.AccessLevels, user.CreatedAt, user.UpdatedAt, user.Token, user.Password, user.Id)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return err
	}

	return nil
}
