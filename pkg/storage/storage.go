package storage

const (
	usersTable              = "users"
	priceNotificationsTable = "price_notifications"
)

//// SaveUser saves a single user to PostgreSQL
//func SaveUser(conn *pgx.Conn, user *models.User) error {
//	query := `
//		INSERT INTO users (Username, Password, Email, Mobile, NotificationPreference, IsAdmin, Role)
//		VALUES ($1, $2, $3, $4, $5, $6, $7)
//	`
//
//	_, err := conn.Exec(context.Background(), query,
//		user.Username,
//		user.Password,
//		user.Email,
//		user.Mobile,
//		user.NotificationPreference,
//		user.IsAdmin,
//		user.Role,
//	)
//	if err != nil {
//		return fmt.Errorf("failed to save user: %v", err)
//	}
//
//	return nil
//}

//// SaveUsers saves multiple users to PostgreSQL
//func SaveUsers(conn *pgx.Conn, users []*models.User) error {
//	query := `
//		INSERT INTO users (Username, Password, Email, Mobile, NotificationPreference, IsAdmin, Role)
//		VALUES ($1, $2, $3, $4, $5, $6, $7)
//	`
//
//	for _, user := range users {
//		_, err := conn.Exec(context.Background(), query,
//			user.Username,
//			user.Password,
//			user.Email,
//			user.Mobile,
//			user.NotificationPreference,
//			user.IsAdmin,
//			user.Role,
//		)
//		if err != nil {
//			return fmt.Errorf("failed to save users: %v", err)
//		}
//	}
//
//	return nil
//}

//// GetUserByUsername retrieves a user by username from PostgreSQL
//func GetUserByUsername(conn *pgx.Conn, username string) (*models.User, error) {
//	query := `
//		SELECT Username, Password, Email, Mobile, NotificationPreference, IsAdmin, Role
//		FROM users
//		WHERE Username = $1
//	`
//
//	row := conn.QueryRow(context.Background(), query, username)
//	var user models.User
//	if err := row.Scan(
//		&user.Username,
//		&user.Password,
//		&user.Email,
//		&user.Mobile,
//		&user.NotificationPreference,
//		&user.IsAdmin,
//		&user.Role,
//	); err != nil {
//		if err == pgx.ErrNoRows {
//			return nil, fmt.Errorf("user %s not found", username)
//		}
//		return nil, fmt.Errorf("failed to get user by username: %v", err)
//	}
//
//	return &user, nil
//}

//// LoadUsers retrieves all users from PostgreSQL
//func LoadUsers(conn *pgx.Conn) ([]*models.User, error) {
//	query := `
//		SELECT Username, Password, Email, Mobile, NotificationPreference, IsAdmin, Role
//		FROM users
//	`
//
//	rows, err := conn.Query(context.Background(), query)
//	if err != nil {
//		return nil, fmt.Errorf("error fetching users: %v", err)
//	}
//	defer rows.Close() // Ensure rows are closed after processing
//
//	var users []*models.User
//	for rows.Next() {
//		var user models.User
//		if err := rows.Scan(
//			&user.Username,
//			&user.Password,
//			&user.Email,
//			&user.Mobile,
//			&user.NotificationPreference,
//			&user.IsAdmin,
//			&user.Role,
//		); err != nil {
//			return nil, fmt.Errorf("failed to scan user: %v", err)
//		}
//		users = append(users, &user)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("rows error while fetching users: %v", err)
//	}
//
//	return users, nil
//}
