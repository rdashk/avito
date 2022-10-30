package models

import (
	"database/sql"
	"fmt"
	"log"
)

// DB exported global variable to hold the database connection
var DB *sql.DB

type UserBalance struct {
	UserName string          `json:"user_name"`
	Balance  sql.NullFloat64 `json:"balance"`
}

type Balance struct {
	UserID  int             `json:"user_id"`
	Balance sql.NullFloat64 `json:"balance"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Report struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	OrderID   int     `json:"order_id"`
	ServiceID int     `json:"service_id"`
	Money     float64 `json:"Money"`
}

// AllUsers returns a slice of all users in the users table.
func AllUsers() ([]User, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersList []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usersList, nil
}

// AllBalances returns a slice of all balances in the balance table.
func AllBalances() ([]Balance, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM balance")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balanceList []Balance

	for rows.Next() {
		var balance Balance
		err := rows.Scan(&balance.UserID, &balance.Balance)
		if err != nil {
			return nil, err
		}
		balanceList = append(balanceList, balance)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return balanceList, nil
}

// AllReserves returns a slice of all reserve_money in the report table.
func AllReserves() ([]Report, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM reserve_money")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reserveList []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(&r.ID, &r.UserID, &r.ServiceID, &r.OrderID, &r.Money)
		if err != nil {
			return nil, err
		}
		reserveList = append(reserveList, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reserveList, nil
}

// AllReports returns a slice of all reports in the report table.
func AllReports() ([]Report, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM report")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportList []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(&r.ID, &r.UserID, &r.ServiceID, &r.OrderID, &r.Money)
		if err != nil {
			return nil, err
		}
		reportList = append(reportList, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reportList, nil
}

// BalanceUser returns a balance current user from balance table
func BalanceUser(userId int) (UserBalance, error) {

	var balance UserBalance
	err := DB.QueryRow("SELECT u.name, b.balance FROM users u LEFT OUTER JOIN balance b ON b.id_user = u.id WHERE b.id_user = ?", userId).Scan(&balance.UserName, &balance.Balance)
	//err := DB.QueryRow("SELECT * FROM balance WHERE balance.id_user = ?", userId).Scan(&balance.UserID, &balance.Balance)
	if err == sql.ErrNoRows {
		DB.QueryRow("INSERT INTO balance(id_user) VALUES (?)", userId)
	}
	err = DB.QueryRow("SELECT u.name, b.balance FROM users u LEFT OUTER JOIN balance b ON b.id_user = u.id WHERE b.id_user = ?", userId).Scan(&balance.UserName, &balance.Balance)
	if err != nil {
		log.Fatal("Change user id!")
		panic(err.Error())
	}
	return balance, nil
}

// ChangeMoneyToBalance adding Money in the balance table for current user
func ChangeMoneyToBalance(userId int, money float64) {

	balance, err := BalanceUser(userId)
	if err != nil {
		panic(err.Error())
		return
	}
	newBalance := balance.Balance.Float64 + money
	DB.QueryRow("UPDATE `balance` SET `balance`= ? WHERE id_user = ?", newBalance, userId)
}

// ReserveMoney add new record to reserve_money table
func ReserveMoney(userId int, serviceId int, orderId int, money float64) {

	// checking money amount on balance for reserving
	var amount float64
	DB.QueryRow("SELECT balance FROM balance WHERE id_user = ?", userId).Scan(&amount)
	fmt.Println(amount)
	if amount >= money {
		DB.QueryRow("INSERT INTO reserve_money(`user_id`, `service_id`, `order_id`, `money`) VALUES (?,?,?,?)",
			userId, serviceId, orderId, money)
	} else {
		log.Println("Not enough money on balance!")
	}
}

// DebitingFunds delete record from reserve_money table and add this record to report table
func DebitingFunds(id int) {

	var r Report
	DB.QueryRow("SELECT * FROM reserve_money WHERE id = ?", id).Scan(&r.ID, &r.UserID, &r.ServiceID, &r.OrderID, &r.Money)
	ChangeMoneyToBalance(r.UserID, -1*r.Money)
	DB.QueryRow("DELETE FROM reserve_money WHERE id = ?", id)
	DB.QueryRow("INSERT INTO report(`user_id`, `service_id`, `order_id`, `cash`) VALUES (?,?,?,?)",
		r.UserID, r.ServiceID, r.OrderID, r.Money)
}
