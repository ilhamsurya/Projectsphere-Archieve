package auth

import (
	"context"

	"github.com/JesseNicholas00/HaloSuster/utils/errorutil"
	"github.com/lib/pq"
)

func (repo *authRepositoryImpl) CreateUser(
	ctx context.Context,
	user User,
) (res User, err error) {
	if err = ctx.Err(); err != nil {
		return
	}

	ctx, sess, err := repo.dbRizzer.GetOrNoTx(ctx)
	if err != nil {
		err = errorutil.AddCurrentContext(err)
		return
	}

	rows, err := sess.
		NamedStmt(ctx, repo.statements.createUser).
		QueryxContext(ctx, user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			err = ErrDuplicateUser
		} else {
			err = errorutil.AddCurrentContext(err)
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&res)
		if err != nil {
			err = errorutil.AddCurrentContext(err)
			return
		}
	}

	return
}
