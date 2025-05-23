// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package pg_db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addFriendStmt, err = db.PrepareContext(ctx, addFriend); err != nil {
		return nil, fmt.Errorf("error preparing query AddFriend: %w", err)
	}
	if q.createNewPlayerStmt, err = db.PrepareContext(ctx, createNewPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNewPlayer: %w", err)
	}
	if q.createUserSessionStmt, err = db.PrepareContext(ctx, createUserSession); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserSession: %w", err)
	}
	if q.getAllUserSessionByIDStmt, err = db.PrepareContext(ctx, getAllUserSessionByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllUserSessionByID: %w", err)
	}
	if q.getDestinationByIDStmt, err = db.PrepareContext(ctx, getDestinationByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetDestinationByID: %w", err)
	}
	if q.getFriendsIdListOfPlayerByIDStmt, err = db.PrepareContext(ctx, getFriendsIdListOfPlayerByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetFriendsIdListOfPlayerByID: %w", err)
	}
	if q.getLeaderboardDetailsStmt, err = db.PrepareContext(ctx, getLeaderboardDetails); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeaderboardDetails: %w", err)
	}
	if q.getLeaderboardForFriendsStmt, err = db.PrepareContext(ctx, getLeaderboardForFriends); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeaderboardForFriends: %w", err)
	}
	if q.getPlayerByIdStmt, err = db.PrepareContext(ctx, getPlayerById); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerById: %w", err)
	}
	if q.getRandomDestinationForSessionsStmt, err = db.PrepareContext(ctx, getRandomDestinationForSessions); err != nil {
		return nil, fmt.Errorf("error preparing query GetRandomDestinationForSessions: %w", err)
	}
	if q.getRandomDestinationsStmt, err = db.PrepareContext(ctx, getRandomDestinations); err != nil {
		return nil, fmt.Errorf("error preparing query GetRandomDestinations: %w", err)
	}
	if q.getRandomDestinationsForQuestionsStmt, err = db.PrepareContext(ctx, getRandomDestinationsForQuestions); err != nil {
		return nil, fmt.Errorf("error preparing query GetRandomDestinationsForQuestions: %w", err)
	}
	if q.getRandomDestinationsForSessionQuestionsStmt, err = db.PrepareContext(ctx, getRandomDestinationsForSessionQuestions); err != nil {
		return nil, fmt.Errorf("error preparing query GetRandomDestinationsForSessionQuestions: %w", err)
	}
	if q.getUserSessionByIDStmt, err = db.PrepareContext(ctx, getUserSessionByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserSessionByID: %w", err)
	}
	if q.updatePlayerScoreStmt, err = db.PrepareContext(ctx, updatePlayerScore); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayerScore: %w", err)
	}
	if q.updateUserSessionStmt, err = db.PrepareContext(ctx, updateUserSession); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserSession: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addFriendStmt != nil {
		if cerr := q.addFriendStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addFriendStmt: %w", cerr)
		}
	}
	if q.createNewPlayerStmt != nil {
		if cerr := q.createNewPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createNewPlayerStmt: %w", cerr)
		}
	}
	if q.createUserSessionStmt != nil {
		if cerr := q.createUserSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserSessionStmt: %w", cerr)
		}
	}
	if q.getAllUserSessionByIDStmt != nil {
		if cerr := q.getAllUserSessionByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllUserSessionByIDStmt: %w", cerr)
		}
	}
	if q.getDestinationByIDStmt != nil {
		if cerr := q.getDestinationByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDestinationByIDStmt: %w", cerr)
		}
	}
	if q.getFriendsIdListOfPlayerByIDStmt != nil {
		if cerr := q.getFriendsIdListOfPlayerByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFriendsIdListOfPlayerByIDStmt: %w", cerr)
		}
	}
	if q.getLeaderboardDetailsStmt != nil {
		if cerr := q.getLeaderboardDetailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeaderboardDetailsStmt: %w", cerr)
		}
	}
	if q.getLeaderboardForFriendsStmt != nil {
		if cerr := q.getLeaderboardForFriendsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeaderboardForFriendsStmt: %w", cerr)
		}
	}
	if q.getPlayerByIdStmt != nil {
		if cerr := q.getPlayerByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerByIdStmt: %w", cerr)
		}
	}
	if q.getRandomDestinationForSessionsStmt != nil {
		if cerr := q.getRandomDestinationForSessionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRandomDestinationForSessionsStmt: %w", cerr)
		}
	}
	if q.getRandomDestinationsStmt != nil {
		if cerr := q.getRandomDestinationsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRandomDestinationsStmt: %w", cerr)
		}
	}
	if q.getRandomDestinationsForQuestionsStmt != nil {
		if cerr := q.getRandomDestinationsForQuestionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRandomDestinationsForQuestionsStmt: %w", cerr)
		}
	}
	if q.getRandomDestinationsForSessionQuestionsStmt != nil {
		if cerr := q.getRandomDestinationsForSessionQuestionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRandomDestinationsForSessionQuestionsStmt: %w", cerr)
		}
	}
	if q.getUserSessionByIDStmt != nil {
		if cerr := q.getUserSessionByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserSessionByIDStmt: %w", cerr)
		}
	}
	if q.updatePlayerScoreStmt != nil {
		if cerr := q.updatePlayerScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerScoreStmt: %w", cerr)
		}
	}
	if q.updateUserSessionStmt != nil {
		if cerr := q.updateUserSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserSessionStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                           DBTX
	tx                                           *sql.Tx
	addFriendStmt                                *sql.Stmt
	createNewPlayerStmt                          *sql.Stmt
	createUserSessionStmt                        *sql.Stmt
	getAllUserSessionByIDStmt                    *sql.Stmt
	getDestinationByIDStmt                       *sql.Stmt
	getFriendsIdListOfPlayerByIDStmt             *sql.Stmt
	getLeaderboardDetailsStmt                    *sql.Stmt
	getLeaderboardForFriendsStmt                 *sql.Stmt
	getPlayerByIdStmt                            *sql.Stmt
	getRandomDestinationForSessionsStmt          *sql.Stmt
	getRandomDestinationsStmt                    *sql.Stmt
	getRandomDestinationsForQuestionsStmt        *sql.Stmt
	getRandomDestinationsForSessionQuestionsStmt *sql.Stmt
	getUserSessionByIDStmt                       *sql.Stmt
	updatePlayerScoreStmt                        *sql.Stmt
	updateUserSessionStmt                        *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                    tx,
		tx:                                    tx,
		addFriendStmt:                         q.addFriendStmt,
		createNewPlayerStmt:                   q.createNewPlayerStmt,
		createUserSessionStmt:                 q.createUserSessionStmt,
		getAllUserSessionByIDStmt:             q.getAllUserSessionByIDStmt,
		getDestinationByIDStmt:                q.getDestinationByIDStmt,
		getFriendsIdListOfPlayerByIDStmt:      q.getFriendsIdListOfPlayerByIDStmt,
		getLeaderboardDetailsStmt:             q.getLeaderboardDetailsStmt,
		getLeaderboardForFriendsStmt:          q.getLeaderboardForFriendsStmt,
		getPlayerByIdStmt:                     q.getPlayerByIdStmt,
		getRandomDestinationForSessionsStmt:   q.getRandomDestinationForSessionsStmt,
		getRandomDestinationsStmt:             q.getRandomDestinationsStmt,
		getRandomDestinationsForQuestionsStmt: q.getRandomDestinationsForQuestionsStmt,
		getRandomDestinationsForSessionQuestionsStmt: q.getRandomDestinationsForSessionQuestionsStmt,
		getUserSessionByIDStmt:                       q.getUserSessionByIDStmt,
		updatePlayerScoreStmt:                        q.updatePlayerScoreStmt,
		updateUserSessionStmt:                        q.updateUserSessionStmt,
	}
}
