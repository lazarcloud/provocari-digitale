package database

import (
	"errors"
	"net/mail"
	"regexp"

	"github.com/google/uuid"
	"github.com/lazarcloud/provocari-digitale/api/auth"
)

func IsValidEmail(email string) (bool, error) {
	// Check if email field is empty
	if email == "" {
		return false, errors.New("email field is empty")
	}

	// Check if email is in valid format
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, errors.New("invalid email format")
	}

	// Check if email contains illegal characters or patterns
	illegalPattern := regexp.MustCompile(`[\(\)\<\>\,\;\:\\\"\[\]]`)
	if illegalPattern.MatchString(email) {
		return false, errors.New("email contains illegal characters")
	}

	return true, nil
}
func GenerateRandomUsername() string {
	return "user-" + GenerateUUID()
}
func CreateNewUser(email string, password string) (userId string, err error) {
	// Check if email or password is empty
	if email == "" || password == "" {
		return "", errors.New("email or password cannot be empty")
	}

	// Hash password
	hashedPassword := auth.HashString(password)

	// Insert into db
	myUUID := GenerateUUID()
	statement, err := DB.Prepare("INSERT INTO users (id, email, password, username) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = statement.Exec(myUUID, email, hashedPassword, GenerateRandomUsername())
	if err != nil {
		return "", err
	}

	return myUUID, nil
}

func GetUserById(id string) (map[string]interface{}, error) {
	if id == "" {
		return nil, errors.New("userId cannot be empty")
	}

	// Get user from db
	rows, err := DB.Query("SELECT id, created_at, email, password FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user map[string]interface{} = make(map[string]interface{}, 0)
	for rows.Next() {
		var id []byte
		var createdAt int64
		var email string
		var password string
		err := rows.Scan(&id, &createdAt, &email, &password)
		if err != nil {
			return nil, err
		}

		readableId := uuid.UUID(id).String()
		user = map[string]interface{}{
			"id":         readableId,
			"created_at": createdAt,
			"email":      email,
			"password":   password,
		}
	}

	return user, nil
}

func GetUserByEmail(email string) (map[string]interface{}, error) {
	// Check if email or password is empty
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	// Get user from db
	rows, err := DB.Query("SELECT id, created_at, email, password FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user map[string]interface{} = make(map[string]interface{}, 0)
	for rows.Next() {
		var id []byte
		var createdAt int64
		var email string
		var password string
		err := rows.Scan(&id, &createdAt, &email, &password)
		if err != nil {
			return nil, err
		}

		readableId := uuid.UUID(id).String()
		user = map[string]interface{}{
			"id":         readableId,
			"created_at": createdAt,
			"email":      email,
			"password":   password,
		}
	}

	return user, nil
}

func CheckUserExistsByEmail(email string) (bool, error) {
	// Check if email or password is empty
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	user, err := GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	if user["email"] == nil {
		return false, nil
	}

	return true, nil
}
