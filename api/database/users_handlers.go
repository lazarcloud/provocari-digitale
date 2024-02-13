package database

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/auth"
	"github.com/lazarcloud/provocari-digitale/api/auth/jwt"
	"github.com/lazarcloud/provocari-digitale/api/globals"
	"github.com/lazarcloud/provocari-digitale/api/utils"
)

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		if err.Error() == "EOF" {
			utils.RespondWithError(w, "Please send a request body")
			return
		}
		utils.RespondWithError(w, err.Error())
		return
	}

	email := registerRequest.Email

	if email == "" {
		utils.RespondWithError(w, "Email cannot be empty")
		return
	}

	_, err = IsValidEmail(email)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	password := registerRequest.Password

	if password == "" {
		utils.RespondWithError(w, "Password cannot be empty")
		return
	}

	confirmPassword := registerRequest.ConfirmPassword

	if confirmPassword == "" {
		utils.RespondWithError(w, "Password2 cannot be empty")
		return
	}

	if password != confirmPassword {
		utils.RespondWithError(w, "Passwords do not match")
		return
	}

	userExists, err := CheckUserExistsByEmail(email)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	if userExists {
		utils.RespondWithError(w, "User already exists")
		return
	}

	userId, err := CreateNewUser(email, password)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"userId": userId,
	})
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		if err.Error() == "EOF" {
			utils.RespondWithError(w, "Please send a request body")
			return
		}
		utils.RespondWithError(w, err.Error())
		return
	}

	email := loginRequest.Email

	if email == "" {
		utils.RespondWithError(w, "Email cannot be empty")
		return
	}

	_, err = IsValidEmail(email)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	password := loginRequest.Password

	if password == "" {
		utils.RespondWithError(w, "Password cannot be empty")
		return
	}

	var user map[string]interface{}

	user, err = GetUserByEmail(email)

	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	//user doesnt exist
	if user["email"] == nil {
		utils.RespondWithError(w, "User does not exist")
		return
	}

	//check if password is correct with hash
	if !auth.CompareHash(password, user["password"].(string)) {
		utils.RespondWithError(w, "Incorrect password")
		return
	}

	//create jwt token
	tokenStr, err := jwt.CreateJWTWithClaims(globals.AuthAccessType, globals.AuthAccessTypeDuration, user["id"].(string), globals.AuthRoleLoggedIn)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	//create refresh token
	refreshTokenStr, err := jwt.CreateJWTWithClaims(globals.AuthRefreshType, globals.AuthRefreshTypeDuration, user["id"].(string), globals.AuthRoleLoggedIn)

	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"status":       "ok",
		"accessToken":  tokenStr,
		"refreshToken": refreshTokenStr,
		"expiresIn":    globals.AuthRefreshTypeDuration,
	})
}
func GetEmailHandler(w http.ResponseWriter, r *http.Request) {

	userId := auth.GetUserId(r)

	if userId == "" {
		utils.RespondWithError(w, "Not logged in")
		return
	}

	user, err := GetUserById(userId)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status": "ok",
		"email":  user["email"],
	})

}
func PrepareAuthRouter(authRouter *mux.Router) {
	authRouter.HandleFunc("/login", LoginHandler).Methods("POST")
	authRouter.HandleFunc("/register", RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/email", GetEmailHandler).Methods("GET")
}
