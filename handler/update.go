package handler

import (
	"context"
	"encoding/json"
	"io"
	"time"

	sqrl "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var psql = sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)

// Update is a mission update.
type Update struct {
	Drift   float64 `json:"drift"`
	Elapsed float64 `json:"elapsed"`
	Mission struct {
		Action string `json:"action"`
		Clock  struct {
			Launched   string  `json:"launched"`
			Observer   string  `json:"observer"`
			ObserverEt float64 `json:"observer_et"`
			Relative   string  `json:"relative"`
			RelativeEt float64 `json:"relative_et"`
		} `json:"clock"`
		State struct {
			V string `json:"v"`
			X string `json:"x"`
		} `json:"state"`
	} `json:"mission"`
}

func (u Update) Args() []interface{} {
	args := []interface{}{
		time.Now().Unix(),
		u.Drift,
		u.Elapsed,
		u.Mission.Action,
		u.Mission.Clock.Launched,
		u.Mission.Clock.Observer,
		u.Mission.Clock.ObserverEt,
		u.Mission.Clock.Relative,
		u.Mission.Clock.RelativeEt,
		u.Mission.State.V,
		u.Mission.State.X,
	}

	return args
}

func (u *Update) Store(ctx context.Context, db *pgxpool.Pool) error {
	query := psql.Insert("updates").Columns(columns...).Values(u.Args()...)
	stmt, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, stmt, args...)
	return err
}

var columns = []string{
	"created",
	"drift",
	"elapsed",
	"action",
	"launched",
	"observer",
	"observer_et",
	"relative",
	"relative_et",
	"velocity",
	"distance",
}

func FetchLastUpdate(ctx context.Context, db *pgxpool.Pool) (int64, *Update, error) {
	query := sqrl.Select(columns...).
		From("updates").
		OrderBy("created desc").
		Limit(1)
	stmt, args, err := query.ToSql()
	if err != nil {
		return 0, nil, err
	}

	u := &Update{}
	var created int64
	row := db.QueryRow(ctx, stmt, args...)
	err = row.Scan(
		&created,
		&u.Drift,
		&u.Elapsed,
		&u.Mission.Action,
		&u.Mission.Clock.Launched,
		&u.Mission.Clock.Observer,
		&u.Mission.Clock.ObserverEt,
		&u.Mission.Clock.Relative,
		&u.Mission.Clock.RelativeEt,
		&u.Mission.State.V,
		&u.Mission.State.X,
	)

	if err != nil {
		return 0, nil, err
	}

	return created, u, nil
}

func UpdateFromReader(r io.Reader) (*Update, error) {
	u := &Update{}
	dec := json.NewDecoder(r)
	err := dec.Decode(u)
	return u, err
}

func ClearAllUpdates(ctx context.Context, db *pgxpool.Pool) error {
	query := psql.Delete("updates")
	stmt, _, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, stmt)
	return err
}
