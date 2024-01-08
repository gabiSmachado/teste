package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gabiSmachado/intents/datamodel"
	_ "github.com/go-sql-driver/mysql"
)

func CurrentId(db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("SELECT MAX(id) FROM intents").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}


func Insert(db *sql.DB, intent datamodel.Intent) (int, error) {
	_, err := db.Exec(`INSERT INTO intents (name,label,day_of_the_week,start_tiime,end_time,
		minimum_cell_offset,maximum_cell_offset) VALUES (?,?,?,?,?,?,?)`, intent.Name,
		intent.Condition.Labels, intent.Condition.When.DayOfWeek, intent.Condition.When.TimeSpan.StartTime,
		intent.Condition.When.TimeSpan.EndTime, intent.Objective.MinimumCellOffset, intent.Objective.MaximumCellOffset)
	if err != nil {
		log.Printf("Error %s when inserting in table", err)
		return 0, err
	}
	id, _ := CurrentId(db)
	return id, nil
}

func ListIntents(db *sql.DB) ([]datamodel.Intent, error) {
	rows, err := db.Query("SELECT id,name FROM intents")
	if err != nil {
		log.Printf("Error %s when listing intents", err)
		return nil, err
	}
	defer rows.Close()

	var intents []datamodel.Intent
	var name string
	var id int
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			log.Printf("Error %s when listing intents", err)
		}
		intent := datamodel.Intent{Idx: id, Name: name}
		intents = append(intents, intent)
	}
	if err = rows.Err(); err != nil {
		//log.Printf("Error %s when listing intents", err)
		return nil, err
	}

	return intents, nil

}

func DeleteIntent(db *sql.DB, id int) error {
	_, err := db.Query("DELETE FROM intents WHERE id = ?", id)
	return err
}

func DBconnect() (*sql.DB, error) {
	log.Printf("teste")
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/intent")
	if err != nil {
		panic(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
		return nil, err
	} else {
		fmt.Println("Connection successful!")
		return db, nil
	}

}

func IntentShow(db *sql.DB, id int) (datamodel.Intent, error) {
	var name, day, start, end, label string
	var min, max int
	err := db.QueryRow("SELECT * FROM intents WHERE id = ?", id).Scan(&id, &name,
		&label, &day, &start, &end, &min, &max)
	intent := datamodel.Intent{
		Idx:  id,
		Name: name,
		Condition: datamodel.Condition{
			When: datamodel.When{
				DayOfWeek: day,
				TimeSpan: datamodel.TimeSpan{
					StartTime: start,
					EndTime:   end,
				},
			},
			Labels: label,
		},
		Objective: datamodel.Objective{
			MinimumCellOffset: min,
			MaximumCellOffset: max,
		}}

	if err != nil {
		log.Printf("Error %s when selecting intent", err)
		return intent, err
	}
	return intent, nil
}

