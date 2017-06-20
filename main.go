package main

import (
    "os"
    "log"
    "time"
    "database/sql"

    "github.com/blankbook/shared/web"
)

const updateRanksInterval = 1000
const databaseUsernameEnvVar = "BB_CONTENT_DB_USERNAME"
const databasePasswordEnvVar = "BB_CONTENT_DB_PASSWORD"
const databaseServerEnvVar = "BB_CONTENT_DB_SERVER"
const dbName = "blankbookcontent"

// SetupRoutes configures the service API endpoints
func main() {
    dbUsername := os.Getenv(databaseUsernameEnvVar)
    dbPassword := os.Getenv(databasePasswordEnvVar)
    dbServer := os.Getenv(databaseServerEnvVar)
    db, err := web.GetMSSqlDatabase(dbUsername, dbPassword, dbServer, dbName)
    if err != nil {
        log.Panic(err.Error())
    }
    UpdateRanks(db)
    for range time.Tick(time.Second * updateRanksInterval) {
        UpdateRanks(db)
    }
    <-make(chan bool) // prevent exiting
}

func UpdateRanks(db *sql.DB) {
    query :=`
        UPDATE a 
        SET Rank = b.NewRank
        FROM Posts a 
        INNER JOIN 
            (SELECT Score, RANK() OVER (ORDER BY Score DESC) AS NewRank FROM Posts) b
        ON a.Score = b.Score`
    _, err := db.Exec(query)
    if err != nil {
        log.Printf(err.Error())
    }
}
