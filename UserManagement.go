package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

//go:embed user_management_init.sql
var _migrationSQL_ string

// ============================================================
// Models
// ============================================================

type _User_ struct {
	ID_        int64     `json:"_id_" db:"_id_"`
	Username_  string    `json:"_username_" db:"_username_"`
	Email_     string    `json:"_email_" db:"_email_"`
	Password_  string    `json:"_password_" db:"_password_"`
	FirstName_ string    `json:"_first_name_" db:"_first_name_"`
	LastName_  string    `json:"_last_name_" db:"_last_name_"`
	IsActive_  bool      `json:"_is_active_" db:"_is_active_"`
	CreatedAt_ time.Time `json:"_created_at_" db:"_created_at_"`
	UpdatedAt_ time.Time `json:"_updated_at_" db:"_updated_at_"`
}

func (u _User_) MarshalJSON() ([]byte, error) {
	type Alias _User_
	return json.Marshal(&struct {
		Alias
		Password_ string `json:"_password_"`
	}{
		Alias:     (Alias)(u),
		Password_: "[REDACTED]",
	})
}

type _Role_ struct {
	ID_          int64     `json:"_id_" db:"_id_"`
	Name_        string    `json:"_name_" db:"_name_"`
	Description_ string    `json:"_description_" db:"_description_"`
	CreatedAt_   time.Time `json:"_created_at_" db:"_created_at_"`
}

type _UserRole_ struct {
	UserID_     int64     `json:"_user_id_" db:"_user_id_"`
	RoleID_     int64     `json:"_role_id_" db:"_role_id_"`
	AssignedAt_ time.Time `json:"_assigned_at_" db:"_assigned_at_"`
}

// ============================================================
// Repositories
// ============================================================

type _UserRepo_ struct{ db *sql.DB }

func (r *_UserRepo_) Create(u *_User_) error {
	return r.db.QueryRow(
		`INSERT INTO _users_ (_username_, _email_, _password_, _first_name_, _last_name_, _is_active_)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING _id_, _created_at_, _updated_at_`,
		u.Username_, u.Email_, u.Password_, u.FirstName_, u.LastName_, u.IsActive_,
	).Scan(&u.ID_, &u.CreatedAt_, &u.UpdatedAt_)
}

func (r *_UserRepo_) GetByID(id int64) (*_User_, error) {
	u := &_User_{}
	err := r.db.QueryRow(
		`SELECT _id_, _username_, _email_, _password_, _first_name_, _last_name_, _is_active_, _created_at_, _updated_at_
		 FROM _users_ WHERE _id_ = $1`, id,
	).Scan(&u.ID_, &u.Username_, &u.Email_, &u.Password_, &u.FirstName_, &u.LastName_, &u.IsActive_, &u.CreatedAt_, &u.UpdatedAt_)
	return u, err
}

func (r *_UserRepo_) List() ([]_User_, error) {
	rows, err := r.db.Query(
		`SELECT _id_, _username_, _email_, _password_, _first_name_, _last_name_, _is_active_, _created_at_, _updated_at_
		 FROM _users_ ORDER BY _id_`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var _users_ []_User_
	for rows.Next() {
		var u _User_
		if err := rows.Scan(&u.ID_, &u.Username_, &u.Email_, &u.Password_, &u.FirstName_, &u.LastName_, &u.IsActive_, &u.CreatedAt_, &u.UpdatedAt_); err != nil {
			return nil, err
		}
		_users_ = append(_users_, u)
	}
	return _users_, rows.Err()
}

func (r *_UserRepo_) Update(u *_User_) error {
	return r.db.QueryRow(
		`UPDATE _users_ SET _username_ = $1, _email_ = $2, _password_ = $3, _first_name_ = $4, _last_name_ = $5, _is_active_ = $6, _updated_at_ = NOW()
		 WHERE _id_ = $7
		 RETURNING _updated_at_`,
		u.Username_, u.Email_, u.Password_, u.FirstName_, u.LastName_, u.IsActive_, u.ID_,
	).Scan(&u.UpdatedAt_)
}

func (r *_UserRepo_) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM _users_ WHERE _id_ = $1`, id)
	return err
}

// ----

type _RoleRepo_ struct{ db *sql.DB }

func (r *_RoleRepo_) Create(_role_ *_Role_) error {
	return r.db.QueryRow(
		`INSERT INTO _roles_ (_name_, _description_) VALUES ($1, $2) RETURNING _id_, _created_at_`,
		_role_.Name_, _role_.Description_,
	).Scan(&_role_.ID_, &_role_.CreatedAt_)
}

func (r *_RoleRepo_) GetByID(id int64) (*_Role_, error) {
	_role_ := &_Role_{}
	err := r.db.QueryRow(
		`SELECT _id_, _name_, _description_, _created_at_ FROM _roles_ WHERE _id_ = $1`, id,
	).Scan(&_role_.ID_, &_role_.Name_, &_role_.Description_, &_role_.CreatedAt_)
	return _role_, err
}

func (r *_RoleRepo_) List() ([]_Role_, error) {
	rows, err := r.db.Query(`SELECT _id_, _name_, _description_, _created_at_ FROM _roles_ ORDER BY _id_`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var _roles_ []_Role_
	for rows.Next() {
		var _role_ _Role_
		if err := rows.Scan(&_role_.ID_, &_role_.Name_, &_role_.Description_, &_role_.CreatedAt_); err != nil {
			return nil, err
		}
		_roles_ = append(_roles_, _role_)
	}
	return _roles_, rows.Err()
}

func (r *_RoleRepo_) Update(_role_ *_Role_) error {
	_, err := r.db.Exec(
		`UPDATE _roles_ SET _name_ = $1, _description_ = $2 WHERE _id_ = $3`,
		_role_.Name_, _role_.Description_, _role_.ID_,
	)
	return err
}

func (r *_RoleRepo_) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM _roles_ WHERE _id_ = $1`, id)
	return err
}

// ----

type _UserRoleRepo_ struct{ db *sql.DB }

func (r *_UserRoleRepo_) Assign(_userID_, _roleID_ int64) (*_UserRole_, error) {
	_ur_ := &_UserRole_{}
	err := r.db.QueryRow(
		`INSERT INTO _user_roles_ (_user_id_, _role_id_) VALUES ($1, $2)
		 ON CONFLICT (_user_id_, _role_id_) DO NOTHING
		 RETURNING _user_id_, _role_id_, _assigned_at_`,
		_userID_, _roleID_,
	).Scan(&_ur_.UserID_, &_ur_.RoleID_, &_ur_.AssignedAt_)
	return _ur_, err
}

func (r *_UserRoleRepo_) Remove(_userID_, _roleID_ int64) error {
	_, err := r.db.Exec(`DELETE FROM _user_roles_ WHERE _user_id_ = $1 AND _role_id_ = $2`, _userID_, _roleID_)
	return err
}

func (r *_UserRoleRepo_) GetRolesByUser(_userID_ int64) ([]_Role_, error) {
	rows, err := r.db.Query(
		`SELECT r._id_, r._name_, r._description_, r._created_at_
		 FROM _roles_ r
		 INNER JOIN _user_roles_ ur ON ur._role_id_ = r._id_
		 WHERE ur._user_id_ = $1
		 ORDER BY r._id_`, _userID_,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var _roles_ []_Role_
	for rows.Next() {
		var _role_ _Role_
		if err := rows.Scan(&_role_.ID_, &_role_.Name_, &_role_.Description_, &_role_.CreatedAt_); err != nil {
			return nil, err
		}
		_roles_ = append(_roles_, _role_)
	}
	return _roles_, rows.Err()
}

func (r *_UserRoleRepo_) GetUsersByRole(_roleID_ int64) ([]_User_, error) {
	rows, err := r.db.Query(
		`SELECT u._id_, u._username_, u._email_, u._password_, u._first_name_, u._last_name_, u._is_active_, u._created_at_, u._updated_at_
		 FROM _users_ u
		 INNER JOIN _user_roles_ ur ON ur._user_id_ = u._id_
		 WHERE ur._role_id_ = $1
		 ORDER BY u._id_`, _roleID_,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var _users_ []_User_
	for rows.Next() {
		var u _User_
		if err := rows.Scan(&u.ID_, &u.Username_, &u.Email_, &u.Password_, &u.FirstName_, &u.LastName_, &u.IsActive_, &u.CreatedAt_, &u.UpdatedAt_); err != nil {
			return nil, err
		}
		_users_ = append(_users_, u)
	}
	return _users_, rows.Err()
}

// ============================================================
// HTTP Handlers — _Users_
// ============================================================

func _handleCreateUser_(_repo_ *_UserRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u _User_
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := _repo_.Create(&u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func _handleGetUser_(_repo_ *_UserRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, err := _repo_.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func _handleListUsers_(_repo_ *_UserRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_users_, err := _repo_.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(_users_)
	}
}

func _handleUpdateUser_(_repo_ *_UserRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var u _User_
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.ID_ = id
		if err := _repo_.Update(&u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func _handleDeleteUser_(_repo_ *_UserRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := _repo_.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ============================================================
// HTTP Handlers — _Roles_
// ============================================================

func _handleCreateRole_(_repo_ *_RoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _role_ _Role_
		if err := json.NewDecoder(r.Body).Decode(&_role_); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := _repo_.Create(&_role_); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(_role_)
	}
}

func _handleGetRole_(_repo_ *_RoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		_role_, err := _repo_.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(_role_)
	}
}

func _handleListRoles_(_repo_ *_RoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_roles_, err := _repo_.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(_roles_)
	}
}

func _handleUpdateRole_(_repo_ *_RoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var _role_ _Role_
		if err := json.NewDecoder(r.Body).Decode(&_role_); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_role_.ID_ = id
		if err := _repo_.Update(&_role_); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(_role_)
	}
}

func _handleDeleteRole_(_repo_ *_RoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := _repo_.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ============================================================
// HTTP Handlers — _UserRole_ (assignments)
// ============================================================

func _handleAssignRole_(_urRepo_ *_UserRoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_userID_, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		var body struct {
			RoleID_ int64 `json:"_role_id_"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ur_, err := _urRepo_.Assign(_userID_, body.RoleID_)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(_ur_)
	}
}

func _handleRemoveRole_(_urRepo_ *_UserRoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_userID_, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		_roleID_, err := strconv.ParseInt(r.PathValue("_role_id_"), 10, 64)
		if err != nil {
			http.Error(w, "invalid role id", http.StatusBadRequest)
			return
		}
		if err := _urRepo_.Remove(_userID_, _roleID_); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func _handleGetUserRoles_(_urRepo_ *_UserRoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_userID_, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		_roles_, err := _urRepo_.GetRolesByUser(_userID_)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(_roles_)
	}
}

func _handleGetRoleUsers_(_urRepo_ *_UserRoleRepo_) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_roleID_, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid role id", http.StatusBadRequest)
			return
		}
		_users_, err := _urRepo_.GetUsersByRole(_roleID_)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(_users_)
	}
}

// ============================================================
// Main
// ============================================================

func main() {
	_dbURL_ := os.Getenv("DATABASE_URL")
	if _dbURL_ == "" {
		_dbURL_ = "postgres://_db_user_:_db_pass_@localhost:5432/_db_name_?sslmode=disable"
	}

	db, err := sql.Open("postgres", _dbURL_)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := range 5 {
		if err = db.Ping(); err == nil {
			break
		}
		log.Printf("DB not ready, retrying... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	if _, err := db.Exec(_migrationSQL_); err != nil {
		log.Fatal(err)
	}
	log.Println("Migration complete")

	_userRepo_ := &_UserRepo_{db: db}
	_roleRepo_ := &_RoleRepo_{db: db}
	_urRepo_ := &_UserRoleRepo_{db: db}

	mux := http.NewServeMux()

	// _Users_
	mux.HandleFunc("POST /_users_", _handleCreateUser_(_userRepo_))
	mux.HandleFunc("GET /_users_", _handleListUsers_(_userRepo_))
	mux.HandleFunc("GET /_users_/{id}", _handleGetUser_(_userRepo_))
	mux.HandleFunc("PUT /_users_/{id}", _handleUpdateUser_(_userRepo_))
	mux.HandleFunc("DELETE /_users_/{id}", _handleDeleteUser_(_userRepo_))

	// _Roles_
	mux.HandleFunc("POST /_roles_", _handleCreateRole_(_roleRepo_))
	mux.HandleFunc("GET /_roles_", _handleListRoles_(_roleRepo_))
	mux.HandleFunc("GET /_roles_/{id}", _handleGetRole_(_roleRepo_))
	mux.HandleFunc("PUT /_roles_/{id}", _handleUpdateRole_(_roleRepo_))
	mux.HandleFunc("DELETE /_roles_/{id}", _handleDeleteRole_(_roleRepo_))

	// _UserRole_ assignments
	mux.HandleFunc("POST /_users_/{id}/_roles_", _handleAssignRole_(_urRepo_))
	mux.HandleFunc("DELETE /_users_/{id}/_roles_/{_role_id_}", _handleRemoveRole_(_urRepo_))
	mux.HandleFunc("GET /_users_/{id}/_roles_", _handleGetUserRoles_(_urRepo_))
	mux.HandleFunc("GET /_roles_/{id}/_users_", _handleGetRoleUsers_(_urRepo_))

	fmt.Println("Serving on :100")
	log.Fatal(http.ListenAndServe(":100", mux))
}
