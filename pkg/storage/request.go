// //package storage
// //
// //import (
// //	"context"
// //	"cryptotracker/models"
// //	//"cryptotracker/pkg/globals"
// //	"github.com/fatih/color"
// //	"go.mongodb.org/mongo-driver/bson"
// //	"go.mongodb.org/mongo-driver/mongo"
// //
// //	//"go.mongodb.org/mongo-driver/mongo"
// //	"log"
// //	//"time"
// //)
// //
// //const (
// //	userCollection    = "users"
// //	requestCollection = "requests"
// //	//database         = "cryptotracker" // Change to globals.Database if globally defined
// //)
// //
// //// GetAllUsers retrieves all users from MongoDB
// //func GetAllUsers(client *mongo.Client) ([]*models.User, error) {
// //	coll := client.Database(database).Collection(userCollection)
// //
// //	// Retrieve all users
// //	cursor, err := coll.Find(context.TODO(), bson.M{})
// //	if err != nil {
// //		color.New(color.FgRed).Printf("Error fetching users: %v\n", err)
// //		log.Println("Error fetching users:", err)
// //		return nil, err
// //	}
// //	defer cursor.Close(context.TODO())
// //
// //	var users []*models.User
// //	for cursor.Next(context.TODO()) {
// //		var user models.User
// //		if err := cursor.Decode(&user); err != nil {
// //			color.New(color.FgRed).Printf("Error decoding user data: %v\n", err)
// //			log.Println("Error decoding user data:", err)
// //			return nil, err
// //		}
// //		users = append(users, &user)
// //	}
// //
// //	color.New(color.FgGreen).Println("Users retrieved successfully from MongoDB.")
// //	return users, nil
// //}
package storage

//
//import (
//	"context"
//	"cryptotracker/models"
//	"fmt"
//	"github.com/fatih/color"
//	"github.com/jackc/pgx/v4"
//)
//
////const usersTable = "users"
//
//// GetAllUsers retrieves all users from PostgreSQL
//func GetAllUsers(conn *pgx.Conn) ([]*models.User, error) {
//	query := `
//		SELECT Username, Password, Email, Mobile, NotificationPreference, IsAdmin, Role
//		FROM users
//	`
//
//	rows, err := conn.Query(context.Background(), query)
//	if err != nil {
//		color.New(color.FgRed).Printf("Error fetching users: %v\n", err)
//		return nil, fmt.Errorf("error fetching users: %v", err)
//	}
//	defer rows.Close()
//
//	var users []*models.User
//	for rows.Next() {
//		var user models.User
//		err := rows.Scan(
//			&user.Username,
//			&user.Password,
//			&user.Email,
//			&user.Mobile,
//			&user.NotificationPreference,
//			&user.IsAdmin,
//			&user.Role,
//		)
//		if err != nil {
//			color.New(color.FgRed).Printf("Error decoding user data: %v\n", err)
//			return nil, fmt.Errorf("failed to scan user: %v", err)
//		}
//		users = append(users, &user)
//	}
//
//	if err := rows.Err(); err != nil {
//		color.New(color.FgRed).Printf("Rows error while fetching users: %v\n", err)
//		return nil, fmt.Errorf("rows error while fetching users: %v", err)
//	}
//
//	color.New(color.FgGreen).Println("Users retrieved successfully from PostgreSQL.")
//	return users, nil
//}
