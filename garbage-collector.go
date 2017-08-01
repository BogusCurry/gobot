// Garbage Collector.
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// garbageCollector goes through the database every few hours or so, pings the objects and sees if they're alive, checks their timestamps,
//  and if they are too old
func garbageCollector() {
	log.Println("Garbage Collector called.")
	
	// use a ticker
	ticker := time.NewTicker(time.Hour * 4)
    //ticker := time.NewTicker(time.Minute * 1)
    go func() {
		var agent AgentType
		var position PositionType
		var toDeleteUUIDs []string // collect list of UUIDs to delete
		var callResult string

        for t := range ticker.C {
            log.Println("\n\n\n\n\n" + "running at", t)
            
            // see which Agents are dead
            db, err := sql.Open(PDO_Prefix, GoBotDSN)
			checkErr(err)
			
			rows, err := db.Query("SELECT `UUID`, `PermURL` FROM `Agents`")
			checkErr(err)

			toDeleteUUIDs = nil // reset our array

			for rows.Next() {
				err = rows.Scan(
					&agent.UUID,
					&agent.PermURL,
				)
				checkErr(err) // cannot hurt
				
				// now call the ping on this Agent
				callResult, err = callURL(*agent.PermURL.Ptr(), "command=ping")
				log.Println("Agent", *agent.UUID.Ptr(), "Result of calling", *agent.PermURL.Ptr(), ":", callResult, "Error:", err)
				if err != nil || callResult != "pong" {
					// either dead, zombie, or misbehaving, kill this agent by placing it in the list
					toDeleteUUIDs = append(toDeleteUUIDs, "'" + *agent.UUID.Ptr() + "'")
				}
			}
			rows.Close()
			
			if len(toDeleteUUIDs) > 0 {
				killAgents := strings.Join(toDeleteUUIDs, ",")
				log.Println(funcName(), "The following agents did not reply to the ping command: ", killAgents)
				result, err := db.Exec("DELETE FROM `Agents` WHERE `UUID` IN (" + killAgents + ")")
				checkErr(err)
				rowsAffected, err := result.RowsAffected()
				if err == nil {
					log.Println("deleted", rowsAffected, "zombie agents with UUIDs: ", killAgents) // probably not needed
				} else {
					log.Println("deleted an unknown number of zombie agents with UUIDs: ", killAgents)
				}
			} else {
				log.Println("no agents to delete")
			}			

			// see which Cubes (positions) are dead			
			rows, err = db.Query("SELECT `UUID`, `PermURL` FROM `Positions`")
			checkErr(err)

			toDeleteUUIDs = nil // reset our array again
			
			for rows.Next() {
				err = rows.Scan(
					&position.UUID,
					&position.PermURL,
				)
				checkErr(err)
				
				// now call the ping on this cube
				callResult, err = callURL(*position.PermURL.Ptr(), "command=ping")
				log.Println("Cube", *position.UUID.Ptr(), "Result of calling", *position.PermURL.Ptr(), ":", callResult, "Error:", err)
				if err != nil || callResult != "pong" {
					// either dead, zombie, or misbehaving, kill this cube by placing it in the list
					toDeleteUUIDs = append(toDeleteUUIDs, "'" + *position.UUID.Ptr() + "'")
				}
			}
			rows.Close()

			if len(toDeleteUUIDs) > 0 {
				killPositions := strings.Join(toDeleteUUIDs, ",")
				result, err := db.Exec("DELETE FROM `Positions` WHERE `UUID` IN (" + killPositions + ")")
				checkErr(err)
				rowsAffected, err := result.RowsAffected()
				if err == nil {
					log.Println("deleted", rowsAffected, "zombie cubes with UUIDs: ", killPositions)
				} else {
					log.Println("deleted an unknown number of cubes with UUIDs: ", killPositions)
				}
			} else {
				fmt.Println("no cubes to delete")
			}			

			// see which Objects (positions) are dead
			// This is trickier, because objects are passive!
			// The best we can do is simply to delete old entries — 4 hours is an arbitrary number — and expect that agents will
			// start checking in again new entries
			// TODO(gwyneth): maybe cubes ought also run the Sensorama.lsl script to aid in detecting (and refreshing) 'permanent' objects?
			// Actually, that's a pretty good idea :-) (20170731)	
			result, err := db.Exec("DELETE FROM `Obstacles` WHERE `LastUpdate` < ADDDATE(NOW(), INTERVAL -4 HOUR)")
			checkErr(err)
			rowsAffected, err := result.RowsAffected()
			if err == nil {
				fmt.Println("deleted", rowsAffected, "obstacles which haven't been seen in the past 4 hours")
			} else {
				fmt.Println("Couldn't delete any obstacles; reason:", err)
			}
						
			db.Close()
        }
    }()
}