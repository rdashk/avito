// запуск: iris-cli run
package main

import (
	"database/sql"
	"firstProject/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"log"
)

var db *sql.DB // Database connection pool.

func main() {

	var err error
	models.DB, err = sql.Open("mysql", "root:@/test")
	if err != nil {
		log.Fatal(err)
	}

	app := iris.New()

	// Method:    GET
	// Resource:  http://localhost:800
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("message", "My first service on Golang!")
	})

	// Method:    POST
	// Resource:  http://localhost:800/2/add/12.2
	// app.Get("/{id:string regexp(^[0-9]+$)}")
	app.Get("/{id:uint64}/add/{money}", addMoney)

	// Method:    GET
	// Resource:  http://localhost:800/user/2
	// app.Get("/{id:string regexp(^[0-9]+$)}")
	app.Get("/user/{id:uint64}", getBalance)

	// Method:    POST
	// Resource:  http://localhost:800/user/1/reserve/1/1/12.2
	// app.Get("/{id:string regexp(^[0-9]+$)}")
	app.Get("/user/{id:uint64}/reserve/{service_id:uint}/{order_id:uint}/{money}", reserve)

	// Method:    POST
	// Resource:  http://localhost:800/confirm/1
	// app.Get("/{id:string regexp(^[0-9]+$)}")
	app.Get("/confirm/{id:uint}", confirmDebit)

	app.Listen(":800")
}

// Метод начисления средств на баланс.
// Принимает id пользователя и сколько средств зачислить.
func addMoney(ctx iris.Context) {
	userID, _ := ctx.Params().GetInt("id")
	money, _ := ctx.Params().GetFloat64("money")
	//ctx.Writef("User UserID: %d %f\n\n", userID, money)

	ctx.WriteString("Old balances\n")
	b, err := models.AllBalances()
	if err != nil {
		panic(err.Error())
		return
	}
	for _, bal := range b {
		ctx.Writef("%d, %.2f\n", bal.UserID, bal.Balance.Float64)
	}

	models.ChangeMoneyToBalance(userID, money)

	ctx.WriteString("New balances\n")
	b, err = models.AllBalances()
	if err != nil {
		panic(err.Error())
		return
	}
	for _, bal := range b {
		ctx.Writef("%d, %.2f\n", bal.UserID, bal.Balance.Float64)
	}
}

// Метод резервирования средств с основного баланса на отдельном счете.
// Принимает id пользователя, ИД услуги, ИД заказа, стоимость.
func reserve(ctx iris.Context) {

	userID, _ := ctx.Params().GetInt("id")
	serviceID, _ := ctx.Params().GetInt("service_id")
	orderID, _ := ctx.Params().GetInt("order_id")
	money, _ := ctx.Params().GetFloat64("money")

	models.ReserveMoney(userID, serviceID, orderID, money)

	res, err := models.AllReserves()
	if err != nil {
		log.Print(err)
		panic(err.Error())
		return
	}
	ctx.JSON(res)
}

// Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии.
// Принимает id резервированного счета.
func confirmDebit(ctx iris.Context) {

	id, _ := ctx.Params().GetInt("id")

	res, err := models.AllReserves()
	if err != nil {
		log.Print(err)
		panic(err.Error())
		return
	}
	ctx.JSON(res)

	models.DebitingFunds(id)

	ctx.WriteString("\n\nAllReserves after confirm\n")
	res, err = models.AllReserves()
	if err != nil {
		panic(err.Error())
		return
	}
	ctx.JSON(res)

	ctx.WriteString("\n\nAll reports\n")
	res, err = models.AllReports()
	if err != nil {
		panic(err.Error())
		return
	}
	ctx.JSON(res)
}

// Метод получения баланса пользователя.
// Принимает id пользователя.
func getBalance(ctx iris.Context) {

	userID, _ := ctx.Params().GetInt("id")
	//ctx.Writef("User UserID: %d", userID)

	// get current balance
	balance, err := models.BalanceUser(userID)
	if err != nil {
		panic(err.Error())
		return
	}

	ctx.JSON(balance)
}
