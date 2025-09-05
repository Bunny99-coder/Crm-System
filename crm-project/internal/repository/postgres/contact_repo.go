// File: internal/repository/postgres/contact_repo.go
package postgres

import (
	"crm-project/internal/models"
    "database/sql"
	"github.com/jmoiron/sqlx"
	"context"
	
)

// ContactRepo is a repository for the contacts table.
type ContactRepo struct {
	db *sqlx.DB
}

// NewContactRepo creates a new ContactRepo.
func NewContactRepo(db *sqlx.DB) *ContactRepo {
	return &ContactRepo{db: db}
}

// GetAll retrieves all contacts from the database.
func (r *ContactRepo) GetAll(ctx context.Context) ([]models.Contact, error) {
	

	var contacts []models.Contact
	
	query := "SELECT contact_id, first_name, last_name, email, primary_phone FROM contacts ORDER BY created_at DESC"
	
	err := r.db.SelectContext(ctx, &contacts, query)
	if err != nil {

    return nil, err
	
	}

	return contacts, nil
}


// Add this method to your contact_repo.go file

// Create inserts a new contact into the database.
// It returns the ID of the newly created contact.
func (r *ContactRepo) Create(ctx context.Context, contact models.Contact) (int, error) {
	var newContactID int

	// Note the 'RETURNING contact_id' clause. This is a PostgreSQL feature
	// that allows us to get the ID of the new row immediately.
	query := `INSERT INTO contacts (first_name, last_name, email, primary_phone, secondary_phone, address, city, sub_city, contact_source)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING contact_id`

	// db.QueryRowx is used for queries that are expected to return a single row.
	// We then use .Scan() to assign the returned value to our variable.
	err := r.db.QueryRowxContext(
		ctx,
		query,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.PrimaryPhone,
		contact.SecondaryPhone,
		contact.Address,
		contact.City,
		contact.SubCity,
		contact.ContactSource,
	).Scan(&newContactID)

	if err != nil {
		return 0, err
	}

	return newContactID, nil
}


// GetByID retrieves a single contact from the database by its ID.
func (r *ContactRepo) GetByID(ctx context.Context, id int) (*models.Contact, error) {
	var contact models.Contact

	// Use all columns from your model to ensure everything is populated
	query := `SELECT 
				contact_id, first_name, last_name, email, primary_phone, 
				secondary_phone, address, city, sub_city, contact_source, created_at
			  FROM contacts 
			  WHERE contact_id = $1`

	// db.Get is a convenient sqlx method for fetching a single row
	// and scanning it into a struct.
	err := r.db.GetContext(ctx, &contact, query, id)
	if err != nil {
		// It's idiomatic in Go for repository Get methods to return a specific
		// error when a row isn't found, so the service layer can handle it.
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil to indicate "not found"
		}
		return nil, err // Return other errors (e.g., connection issues)
	}

	return &contact, nil
}



// Update modifies an existing contact in the database.
// It returns an error if the update fails or if the contact does not exist.
func (r *ContactRepo) Update(ctx context.Context, contact models.Contact) error {
	query := `UPDATE contacts SET
				first_name = $1,
				last_name = $2,
				email = $3,
				primary_phone = $4,
				secondary_phone = $5,
				address = $6,
				city = $7,
				sub_city = $8,
				contact_source = $9
			  WHERE contact_id = $10`

	result, err := r.db.ExecContext(
		ctx,
		query,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.PrimaryPhone,
		contact.SecondaryPhone,
		contact.Address,
		contact.City,
		contact.SubCity,
		contact.ContactSource,
		contact.ID, // The ID for the WHERE clause
	)

	if err != nil {
		return err
	}

	// It's good practice to check if any row was actually affected.
	// If RowsAffected is 0, it means the WHERE clause didn't match any rows.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// We can return sql.ErrNoRows to signal "not found" to the service layer.
		return sql.ErrNoRows
	}

	return nil
}


// Delete removes a contact from the database by its ID.
func (r *ContactRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM contacts WHERE contact_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Just like in Update, we check if any row was actually deleted.
	// If not, it means the contact was not found.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Signal "not found"
	}

	return nil
}