package repository

import "database/sql"

type ActorRoleRepository struct {
	DB *sql.DB
}

func NewActorRoleRepository(db *sql.DB) *ActorRoleRepository {
	return &ActorRoleRepository{DB: db}
}

func (r *ActorRoleRepository) GetRoleNameByID(id int) (string, error) {
	var name string
	query := "SELECT name FROM actor_role WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&name)
	return name, err
}
