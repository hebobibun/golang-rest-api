package controller

import (
	"api/config"
	"api/model"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	Mdl model.UserModel
}

func CreateToken(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.InitConfig().JWTKey))
}

func ExtractToken(e echo.Context) int {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userID"].(float64)
		return int(userId)
	}
	return -1
}

func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		tmp := model.User{}
		if err := c.Bind(&tmp); err != nil {
			return c.JSON(http.StatusBadRequest, "Wrong input format")
		}

		res, err := uc.Mdl.Login(tmp.Email, tmp.Password)

		if err != nil {
			if strings.Contains(err.Error(), "matched") {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": err.Error(),
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": err.Error(),
			})
		}
		jwtToken, err := CreateToken(int(res.ID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"token":   jwtToken,
			"message": "Logged in Successfully",
		})
	}
}

func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		tmp := model.User{}
		err := c.Bind(&tmp)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Wrong input format")
		}

		res, err := uc.Mdl.Register(tmp)
		if err != nil {
			log.Println("Query error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, "Unable to process")
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "Registered a new account successfully"})
	}
}

func (uc *UserController) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := uc.Mdl.GetAll()
		if err != nil {
			log.Println("Query error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, "Unable to process")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "Displayed all users data successfully"})
	}
}

func (uc *UserController) GetID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := ExtractToken(c)

		res, err := uc.Mdl.GetByID(id)
		if err != nil {
			log.Println("Query Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, "Unable to process")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "Displayed user data successfully"})
	}
}

func (uc *UserController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := ExtractToken(c)

		body := model.User{}
		err := c.Bind(&body)
		if err != nil {
			log.Println("Bind body error : ", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Please input correctly",
			})
		}
		body.ID = uint(id)
		res, err := uc.Mdl.Update(body)

		if err != nil {
			log.Println("query error ", err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "Updated data successfully",
		})

	}
}

func (uc *UserController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := ExtractToken(c)

		err := uc.Mdl.Delete(id)

		if err != nil {
			log.Println("Delete error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Deleted a user successfully",
		})
	}
}
