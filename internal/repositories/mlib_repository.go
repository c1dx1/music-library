package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"music-library/internal/models"
)

type MLibRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewMLibRepository(db *pgxpool.Pool, log *logrus.Logger) *MLibRepository {
	return &MLibRepository{db: db, log: log}
}

func (r *MLibRepository) GetLibrary(ctx context.Context, filters map[string]interface{}, page, limit int) ([]models.Song, error) {
	r.log.Info("Entering GetLibrary function")
	r.log.Debugf("Filters: %+v, Page: %d, Limit: %d", filters, page, limit)

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		r.log.Error("Failed to acquire DB connection:", err)
		return nil, err
	}
	defer conn.Release()

	offset := (page - 1) * limit
	query := `SELECT s.id, g.group_name, s.song_name, s.release_date, s.text, s.link
		FROM songs AS s
		JOIN groups AS g ON s.group_id = g.id`
	args := []interface{}{}
	index := 1

	for k, v := range filters {
		if k == "g.group_name" || k == "s.song_name" || k == "s.text" {
			query += fmt.Sprintf(" AND %s ILIKE $%d", k, index)
		} else {
			query += fmt.Sprintf(" AND %s = $%d", k, index)
		}
		args = append(args, v)
		index++
		r.log.Debugf("Adding filter: %s = %v", k, v)
	}

	query += fmt.Sprintf("\n\tLIMIT $%d OFFSET $%d", index, index+1)
	args = append(args, limit, offset)
	r.log.Debugf("Query: %s", query)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		r.log.Error("Query execution failed:", err)
		return nil, fmt.Errorf("mlib repo: getLib: query: %w", err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			r.log.Error("Row scanning failed:", err)
			return nil, fmt.Errorf("mlib repo: getLib: rows scan: %w", err)
		}
		r.log.Debugf("Scanned song: %+v", song)
		songs = append(songs, song)
	}

	r.log.Infof("Successfully fetched %d songs", len(songs))
	return songs, nil
}

func (r *MLibRepository) GetText(ctx context.Context, id int) (string, error) {
	r.log.Infof("Entering GetText function for song ID: %d", id)

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		r.log.Error("Failed to acquire DB connection:", err)
		return "", fmt.Errorf("mlib_repo: getText: db acquire: %w", err)
	}
	defer conn.Release()

	var text string
	err = conn.QueryRow(ctx, "SELECT text FROM songs WHERE id=$1", id).Scan(&text)
	if err != nil {
		r.log.Error("QueryRow failed:", err)
		return "", fmt.Errorf("mlib_repo: getText: queryRow: %w", err)
	}

	r.log.Infof("Successfully fetched text for song ID: %d", id)
	return text, nil
}

func (r *MLibRepository) DeleteSong(ctx context.Context, id int) error {
	r.log.Infof("DeleteSong called with ID: %d", id)

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		r.log.Errorf("Failed to acquire database connection: %v", err)
		return fmt.Errorf("mlib_repo: deleteSong: db acquire: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		r.log.Errorf("Failed to begin transaction: %v", err)
		return fmt.Errorf("mlib_repo: deleteSong: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var groupId int
	err = tx.QueryRow(ctx, "SELECT group_id FROM songs WHERE id = $1", id).Scan(&groupId)
	if err != nil {
		r.log.Errorf("Error retrieving group_id for song ID %d: %v", id, err)
		return fmt.Errorf("mlib_repo: deleteSong: queryRow: %w", err)
	}
	r.log.Debugf("Retrieved group_id: %d for song ID: %d", groupId, id)

	_, err = tx.Exec(ctx, "DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		r.log.Errorf("Error deleting song with ID %d: %v", id, err)
		return fmt.Errorf("mlib_repo: deleteSong: delete song: %w", err)
	}
	r.log.Debugf("Successfully deleted song with ID: %d", id)

	var countGroupSongs int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM songs WHERE group_id = $1", groupId).Scan(&countGroupSongs)
	if err != nil {
		r.log.Errorf("Error counting songs for group_id %d: %v", groupId, err)
		return fmt.Errorf("mlib_repo: deleteSong: queryRow: %w", err)
	}

	if countGroupSongs == 0 {
		_, err = tx.Exec(ctx, "DELETE FROM groups WHERE id = $1", groupId)
		if err != nil {
			r.log.Errorf("Error deleting group with ID %d: %v", groupId, err)
			return fmt.Errorf("mlib_repo: deleteSong: delete group: %w", err)
		}
		r.log.Debugf("Successfully deleted group with ID: %d", groupId)
	}

	if err = tx.Commit(ctx); err != nil {
		r.log.Errorf("Transaction commit failed: %v", err)
		return fmt.Errorf("mlib_repo: deleteSong: commit tx: %w", err)
	}

	r.log.Infof("Successfully deleted song with ID: %d", id)
	return nil
}

func (r *MLibRepository) editGroup(ctx context.Context, tx pgx.Tx, id int, update interface{}) error {
	r.log.Infof("editGroup called with song ID: %d and new group: %v", id, update)

	var groupExists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM groups WHERE group_name ILIKE $1)", update).Scan(&groupExists)
	if err != nil {
		r.log.Errorf("Failed to check if group exists for %v: %v", update, err)
		return fmt.Errorf("mlib_repo: editGroup: select exists new group: %w", err)
	}
	r.log.Debugf("Group existence for %v: %t", update, groupExists)

	var oldGroupId int
	err = tx.QueryRow(ctx, "SELECT group_id FROM songs WHERE id = $1", id).Scan(&oldGroupId)
	if err != nil {
		r.log.Errorf("Failed to retrieve old group ID for song ID %d: %v", id, err)
		return fmt.Errorf("mlib_repo: editGroup: queryRow group_id: %w", err)
	}
	r.log.Debugf("Retrieved old group ID: %d for song ID: %d", oldGroupId, id)

	var groupSongsCount int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM songs WHERE group_id = $1", oldGroupId).Scan(&groupSongsCount)
	if err != nil {
		r.log.Errorf("Failed to count songs in group ID %d: %v", oldGroupId, err)
		return fmt.Errorf("mlib_repo: editGroup: queryRow count: %w", err)
	}
	r.log.Debugf("Number of songs in group ID %d: %d", oldGroupId, groupSongsCount)

	var updatedGroupId int
	if groupExists {
		err = tx.QueryRow(ctx, "SELECT id FROM groups WHERE group_name ILIKE $1", update).Scan(&updatedGroupId)
		if err != nil {
			r.log.Errorf("Failed to retrieve ID for existing group %v: %v", update, err)
			return fmt.Errorf("mlib_repo: editGroup: select id group if exists: %w", err)
		}
		r.log.Debugf("Retrieved group ID %d for existing group %v", updatedGroupId, update)

		_, err = tx.Exec(ctx, "UPDATE songs SET group_id = $1 WHERE id = $2", updatedGroupId, id)
		if err != nil {
			r.log.Errorf("Failed to update group ID in songs: %v", err)
			return fmt.Errorf("mlib_repo: editGroup: edit group_id in songs if exists: %w", err)
		}
		r.log.Infof("Updated group ID to %d for song ID %d", updatedGroupId, id)

		if groupSongsCount == 1 {
			_, err = tx.Exec(ctx, "DELETE FROM groups WHERE id = $1", oldGroupId)
			if err != nil {
				r.log.Errorf("Failed to delete old group ID %d: %v", oldGroupId, err)
				return fmt.Errorf("mlib_repo: editGroup: delete old group if exists and count = 1: %w", err)
			}
			r.log.Infof("Deleted old group ID %d", oldGroupId)
		}
	} else {
		if groupSongsCount == 1 {
			_, err = tx.Exec(ctx, "UPDATE groups SET group_name = $1 WHERE id = $2", update, oldGroupId)
			if err != nil {
				r.log.Errorf("Failed to update group name for ID %d: %v", oldGroupId, err)
				return fmt.Errorf("mlib_repo: editGroup: update group if !exists and count = 1: %w", err)
			}
			r.log.Infof("Updated group name to %v for group ID %d", update, oldGroupId)
			updatedGroupId = oldGroupId
		} else {
			err = tx.QueryRow(ctx, "INSERT INTO groups (group_name) VALUES ($1) RETURNING id", update).Scan(&updatedGroupId)
			if err != nil {
				r.log.Errorf("Failed to insert new group %v: %v", update, err)
				return fmt.Errorf("mlib_repo: editGroup: insert new group if !exists and count > 1: %w", err)
			}
			r.log.Infof("Inserted new group %v with ID %d", update, updatedGroupId)

			_, err = tx.Exec(ctx, "UPDATE songs SET group_id = $1 WHERE id = $2", updatedGroupId, id)
			if err != nil {
				r.log.Errorf("Failed to update song with new group ID %d: %v", updatedGroupId, err)
				return fmt.Errorf("mlib_repo: editGroup: edit group_id in songs if !exists and count > 1: %w", err)
			}
			r.log.Infof("Updated song ID %d with new group ID %d", id, updatedGroupId)
		}
	}
	return nil
}

func (r *MLibRepository) EditSong(ctx context.Context, id int, updates map[string]interface{}) error {
	r.log.Infof("EditSong called with song ID: %d and updates: %+v", id, updates)

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		r.log.Errorf("Failed to acquire database connection: %v", err)
		return fmt.Errorf("mlib_repo: editSong: db acquire: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		r.log.Errorf("Failed to begin transaction: %v", err)
		return fmt.Errorf("mlib_repo: editSong: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if val, ok := updates["group_name"]; ok {
		r.log.Debugf("Editing group for song ID: %d with new group name: %v", id, val)
		err := r.editGroup(ctx, tx, id, val)
		if err != nil {
			r.log.Errorf("Failed to edit group for song ID %d: %v", id, err)
			return fmt.Errorf("mlib_repo: editSong: update group: %w", err)
		}
		delete(updates, "group_name")
	}

	if len(updates) >= 1 {
		r.log.Debugf("Updating song fields for song ID: %d with updates: %+v", id, updates)

		query := `UPDATE songs SET `
		args := []interface{}{}
		index := 1

		for k, v := range updates {
			if index > 1 {
				query += ", "
			}
			query += fmt.Sprintf("%s = $%d", k, index)
			args = append(args, v)
			index++
		}

		query += fmt.Sprintf(" WHERE id = $%d", index)
		args = append(args, id)

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			r.log.Errorf("Failed to update song fields for song ID %d: %v", id, err)
			return fmt.Errorf("mlib_repo: editSong: update song: %w", err)
		}
		r.log.Infof("Successfully updated song fields for song ID %d", id)
	}

	if err = tx.Commit(ctx); err != nil {
		r.log.Errorf("Transaction commit failed for song ID %d: %v", id, err)
		return fmt.Errorf("mlib_repo: editSong: commit tx: %w", err)
	}
	r.log.Infof("Successfully edited song with ID: %d", id)
	return nil
}

func (r *MLibRepository) AddSong(ctx context.Context, song models.Song) error {
	r.log.Infof("AddSong called with song: %+v", song)

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		r.log.Errorf("Failed to acquire database connection: %v", err)
		return fmt.Errorf("mlib_repo: AddSong: db acquire: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		r.log.Errorf("Failed to begin transaction: %v", err)
		return fmt.Errorf("mlib_repo: AddSong: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO groups (group_name) VALUES($1) ON CONFLICT (group_name) DO NOTHING", *song.Group)
	if err != nil {
		r.log.Errorf("Error inserting or updating group: %v", err)
		return fmt.Errorf("mlib_repo: AddSong: query insert groups on conflict: %w", err)
	}
	r.log.Debugf("Group inserted or exists: %s", *song.Group)

	var groupId int
	err = tx.QueryRow(ctx, "SELECT id FROM groups WHERE group_name = $1", *song.Group).Scan(&groupId)
	if err != nil {
		r.log.Errorf("Error retrieving group_id for group_name %s: %v", *song.Group, err)
		return fmt.Errorf("mlib_repo: AddSong: query group id from groups: %w", err)
	}
	r.log.Debugf("Retrieved group_id: %d for group_name: %s", groupId, *song.Group)

	_, err = tx.Exec(ctx, "INSERT INTO songs (group_id, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		groupId, *song.Song, *song.ReleaseDate, *song.Text, *song.Link)
	if err != nil {
		r.log.Errorf("Error inserting song: %v", err)
		return fmt.Errorf("mlib_repo: AddSong: insert into songs: %w", err)
	}
	r.log.Debugf("Successfully inserted song: %s", *song.Song)

	if err = tx.Commit(ctx); err != nil {
		r.log.Errorf("Transaction commit failed: %v", err)
		return fmt.Errorf("mlib_repo: AddSong: commit transaction: %w", err)
	}

	r.log.Infof("Successfully added song: %s", *song.Song)
	return nil
}
