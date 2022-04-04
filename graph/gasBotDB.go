package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

func initDB() (*sql.DB, error) {
	var connectName, connectString string

	if viper.GetBool("live") {
		connectName = "DBCONNECT"

	} else {
		connectName = "DBCONNECT_STAGING"

	}
	connectString = viper.GetString(connectName)
	if len(connectString) == 0 {
		log.Fatal("NO DATABASE SET IN CONFIG UNDER", connectName)
	}
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Println("gasbot/db", err)
		return nil, err
	}
	return db, nil
}

func rowExists(db *sql.DB, subquery string, args ...interface{}) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT exists (%s)", subquery)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func saveUserToDatabase(fName string, lastName string, fullName string, ID int) (exists bool, err error) {
	db, err := sql.Open("postgres", viper.GetString("DBCONNECT"))
	if err != nil {
		fmt.Println("saveToDatabase0: ", err)
		return false, err
	}
	defer db.Close()

	ex, err := rowExists(db, "select id from users where user_id=$1", ID)
	if err != nil {
		return false, err
	}
	if ex {
		return true, nil
	}
	version := viper.GetFloat64("Version")
	stmtStr := `insert into users (firstname,lastname,fullname,user_id,dateadded,show_level,show_info,version)
  values
  ($1,$2,$3,$4,now()::timestamp AT TIME ZONE 'Singapore',false,false,$5)
  returning id`
	var id int
	err = db.QueryRow(stmtStr,
		fName,
		lastName,
		fullName,
		ID, version).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	log.Printf("User %s added %d\n", fullName, id)
	return false, nil
}

func saveLevelToDatabase(userID int, level decimal.Decimal) (err error) {
	//query := "insert into gas_levels (user_id,gas_level) values ($1,$2) on conflict on constraint gas_levels_user_id_key do update set gas_level=$2 where gas_levels.user_id=$1 returning id"
	query := "update users set show_level=true, gas_level=$2 where user_id=$1"
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, query, userID, level)
	return
}

func setOne(userID int64, field string, value interface{}) (err error) {
	//query := "insert into gas_levels (user_id,gas_level) values ($1,$2) on conflict on constraint gas_levels_user_id_key do update set gas_level=$2 where gas_levels.user_id=$1 returning id"
	query := fmt.Sprintf("update users set %s=$2 where user_id=$1", field)
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, query, userID, value)
	return
}
func turnOffLevelAlerts(userID int64) (err error) {
	return setOne(userID, "show_level", false)
}

func turnOnInfo(userID int64) (err error) {
	return setOne(userID, "show_info", true)
}

func turnOffInfo(userID int64) (err error) {
	return setOne(userID, "show_info", false)
}

func getAllBelowOrEqual(currentLevel decimal.Decimal) (users []int64, err error) {
	query := "select user_id from users where gas_level >= $1 and show_level"
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, currentLevel)
	if err != nil {
		return
	}
	for rows.Next() {
		var user int64
		rows.Scan(&user)
		users = append(users, user)
	}
	return
}

func getAllWhoWantInfo() (users []int64, err error) {
	query := "select user_id from users where show_info"
	return getUserList(query)
}

func getUserList(query string) (users []int64, err error) {
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	for rows.Next() {
		var user int64
		rows.Scan(&user)
		users = append(users, user)
	}
	return
}

func getUserInfo(userID int) (levelState bool, infoState bool, levelValue decimal.Decimal, err error) {
	query := "select show_level,show_info,gas_level  from users where user_id=$1"
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query, userID).Scan(&levelState, &infoState, &levelValue)
	if err != nil {
		return
	}
	return
}

func getAllUsers() (users []int64, err error) {
	query := "select user_id from users"
	return getUserList(query)
}

func getAllUsersBelowVersion(version float64) (users []int64, err error) {
	query := fmt.Sprintf("select user_id from users where version < %v", version)
	fmt.Println(query)
	return getUserList(query)
}

func updateUserVersion(userID int64, version float64) (err error) {
	return setOne(userID, "version", version)
}

func loadLatestGasData() (gds *GasData, err error) {
	query := "select fastest, fast, safelow,average,fastestwait, fastwait, safelowwait,averagewait,blocknum,blocktime from gasrates order by id desc limit 1"
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query).Scan(&gds.Fastest, &gds.Fast, &gds.SafeLow, &gds.Average, &gds.FastestWait, &gds.FastWait, &gds.SafeLowWait, &gds.AverageWait, &gds.BlockNum, &gds.BlockTime)
	return
}

func loadGasDataBetween(start, end string) (gds []GasData, err error) {
	log.Println(start, stop)
	query := "select fastest, fast, safelow,average,fastestwait, fastwait, safelowwait,averagewait,blocknum,blocktime, dateadded from gasrates where dateadded between $1 and $2 order by blocknum"
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, start, stop)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var gd GasData
		err = rows.Scan(&gd.Fastest, &gd.Fast, &gd.SafeLow, &gd.Average, &gd.FastestWait, &gd.FastWait, &gd.SafeLowWait, &gd.AverageWait, &gd.BlockNum, &gd.BlockTime, &gd.DateAdded)
		if err != nil {
			return
		}
		gds = append(gds, gd)
	}
	return
}

func saveGasData(gd *GasData) (err error) {
	query := "insert into gasrates (fastest, fast, safelow,average,fastestwait, fastwait, safelowwait,averagewait,blocknum,blocktime,dateadded) values  ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,now())"
	db, err := initDB()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, query, gd.Fastest, gd.Fast, gd.SafeLow, gd.Average, gd.FastestWait, gd.FastWait, gd.SafeLowWait, gd.AverageWait, gd.BlockNum, gd.BlockTime)
	return
}
