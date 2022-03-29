package routes

import (
	"awesomeProject3/adapters/database"
	models "awesomeProject3/app/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"time"
)

//{
//	id:1,
//	user:{
//		id:23,
//		first_name:"mel"
//		last_name:"sarı"
//}
//	product:{
//		id:24,
//		name:"ıphone"
//		serial_number:"1235"
//}
//}

type Order struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {

	return Order{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}

}
func CreateOrder(c *fiber.Ctx) error {

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var product models.Product
	if err := findProduct(order.UserRefer, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&order)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)
	return c.Status(200).JSON(responseOrder)
}
func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}
	database.Database.Db.Find(&orders)
	responseOrders := []Order{}

	for _, order := range orders {
		var user models.User
		var product models.Product
		database.Database.Db.Find(&user, "id=?", order.UserRefer)
		database.Database.Db.Find(&product, "id=?", order.ProductRefer)
		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)
	}

	return c.Status(200).JSON(responseOrders)
}
func findOrder(id int, order *models.Order) error {
	database.Database.Db.Find(&order, "id=?", id)
	if order.ID == 0 {
		return errors.New("order does not exist")
	}
	return nil

}
func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var order models.Order

	if err != nil {
		return c.Status(400).JSON("please ensure that : id is an integer")

	}
	if err := findOrder(id, &order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User
	var product models.Product

	database.Database.Db.First(&user, order.UserRefer)
	database.Database.Db.First(&product, order.ProductRefer)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)

}
