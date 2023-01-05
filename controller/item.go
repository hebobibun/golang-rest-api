package controller

import (
	"api/model"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type ItemController struct {
	Mdl model.ItemModel
}

func (ic *ItemController) Insert() echo.HandlerFunc {
	return func(c echo.Context) error {
		tmp := model.Item{}
		err := c.Bind(&tmp)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Wrong input format")
		}
		
		id := ExtractToken(c)
		tmp.UserID = uint(id)

		res, err := ic.Mdl.Insert(tmp)
		if err != nil {
			log.Println("Query error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, "Unable to process")
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "Inserted a new item successfully"})
	}
}

func (ic *ItemController) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := ExtractToken(c)
		res, err := ic.Mdl.GetAll(id)
		if err != nil {
			log.Println("query error", err.Error())
			return c.JSON(http.StatusInternalServerError, "Unable to process")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "Displayed all the items successfully"})
	}
}

func (ic *ItemController) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		idLogin := ExtractToken(c)
		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("Convert id error : ", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Please input number only",
			})
		}
		res, err := ic.Mdl.GetByID(id, idLogin)
		if err != nil {
			if strings.Contains(err.Error(), "Unauthorized") {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": err.Error(),
				})
			}
			log.Println("query error", err.Error())
			return c.JSON(http.StatusInternalServerError, "Unable to process")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "Displayed an item successfully"})
	}
}

func (ic *ItemController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := ExtractToken(c)

		paramID := c.Param("id")
		itemID, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("Convert id error : ", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Please input number only",
			})
		}

		body := model.Item{}
		err = c.Bind(&body)
		if err != nil {
			log.Println("Bind body error : ", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Please input correctly",
			})
		}
		
		body.ID = uint(itemID)
		res, err := ic.Mdl.Update(body, id)

		if err != nil {
			if strings.Contains(err.Error(), "Unauthorized") {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": err.Error(),
				})
			}
			log.Println("query error ", err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "Updated data successfully",
		})

	}
}

func (ic *ItemController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := ExtractToken(c)

		paramID := c.Param("id")
		itemID, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("Convert id error : ", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Please input number only",
			})
		}

		err = ic.Mdl.Delete(itemID, id)

		if err != nil {
			log.Println("Delete error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Deleted an item successfully",
		})
	}
}