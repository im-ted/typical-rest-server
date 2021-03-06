package typrepo

const mysqlTmpl = `package {{.Package}}_repo

import({{range $pkg, $alias := .Imports}}
	{{$alias}} "{{$pkg}}"{{end}}
)

var (
	// {{.Name}}TableName is table name for {{.Table}} entity
	{{.Name}}TableName = "{{.Table}}"
	// {{.Name}}Table is columns for {{.Table}} entity
	{{.Name}}Table = struct {
		{{range .Fields}}{{.Name}} string
		{{end}}
	}{
		{{range .Fields}}{{.Name}}: "{{.Column}}",
		{{end}}
	}
)

type (
	// {{.Name}}Repo to get {{.Table}} data from database
	// @mock
	{{.Name}}Repo interface {
		Count(context.Context, ...sqkit.SelectOption) (int64, error)
		Find(context.Context, ...sqkit.SelectOption) ([]*{{.Package}}.{{.Name}}, error)
		Create(context.Context, *{{.Package}}.{{.Name}}) (int64, error)
		Delete(context.Context, sqkit.DeleteOption) (int64, error)
		Update(context.Context, *{{.Package}}.{{.Name}}, sqkit.UpdateOption) (int64, error)
		Patch(context.Context, *{{.Package}}.{{.Name}}, sqkit.UpdateOption) (int64, error)
	}
	// {{.Name}}RepoImpl is implementation {{.Table}} repository
	{{.Name}}RepoImpl struct {
		dig.In
		*sql.DB {{.CtorDB}}
	}
)

func init() {
	typapp.Provide("",New{{.Name}}Repo)
}

// New{{.Name}}Repo return new instance of {{.Name}}Repo
func New{{.Name}}Repo(impl {{.Name}}RepoImpl) {{.Name}}Repo {
	return &impl
}

// Count {{.Table}}
func (r *{{.Name}}RepoImpl) Count(ctx context.Context, opts ...sqkit.SelectOption) (int64, error) {
	builder := sq.
		Select("count(*)").
		From({{.Name}}TableName).
		RunWith(r)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	row := builder.QueryRowContext(ctx)

	var cnt int64
	if err := row.Scan(&cnt); err != nil {
		return -1, err
	}
	return cnt, nil
}


// Find {{.Table}}
func (r *{{.Name}}RepoImpl) Find(ctx context.Context, opts ...sqkit.SelectOption) (list []*{{.Package}}.{{.Name}}, err error) {
	builder := sq.
		Select(
			{{range .Fields}}{{$.Name}}Table.{{.Name}},
			{{end}}
		).
		From({{.Name}}TableName).
		RunWith(r)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*{{.Package}}.{{.Name}}, 0)
	for rows.Next() {
		ent := new({{.Package}}.{{.Name}})
		if err = rows.Scan({{range .Fields}}
			&ent.{{.Name}},{{end}}
		); err != nil {
			return
		}
		list = append(list, ent)
	}
	return
}

// Create {{.Table}}
func (r *{{.Name}}RepoImpl) Create(ctx context.Context, ent *{{.Package}}.{{.Name}}) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	res, err := sq.
		Insert({{$.Name}}TableName).
		Columns({{range .Fields}}{{if not .PrimaryKey}}	{{$.Name}}Table.{{.Name}},{{end}}	
		{{end}}).
		Values({{range .Fields}}{{if .DefaultValue}}	{{.DefaultValue}},{{else if not .PrimaryKey}}	ent.{{.Name}},{{end}}
		{{end}}).
		RunWith(txn.DB).
		ExecContext(ctx)

	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	lastInsertID, err := res.LastInsertId()
	txn.SetError(err)
	return lastInsertID, err
}


// Update {{.Table}}
func (r *{{.Name}}RepoImpl) Update(ctx context.Context, ent *{{.Package}}.{{.Name}}, opt sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update({{.Name}}TableName).{{range .Fields}}{{if and (not .PrimaryKey) (not .SkipUpdate)}}
		Set({{$.Name}}Table.{{.Name}},{{if .DefaultValue}}{{.DefaultValue}}{{else}}ent.{{.Name}},{{end}}).{{end}}{{end}}
		RunWith(txn.DB)

	if opt != nil {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Patch {{.Table}}
func (r *{{.Name}}RepoImpl) Patch(ctx context.Context, ent *{{.Package}}.{{.Name}}, opt sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.Update({{.Name}}TableName).RunWith(txn.DB)
	{{range .Fields}}{{if and (not .PrimaryKey) (not .SkipUpdate)}}{{if .DefaultValue}}
	builder = builder.Set({{$.Name}}Table.{{.Name}}, {{.DefaultValue}}){{else}}
	if !reflectkit.IsZero(ent.{{.Name}}) {
		builder = builder.Set({{$.Name}}Table.{{.Name}}, ent.{{.Name}})
	}{{end}}{{end}}{{end}}

	if opt != nil {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}


// Delete {{.Table}}
func (r *{{.Name}}RepoImpl) Delete(ctx context.Context, opt sqkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.Delete({{.Name}}TableName).RunWith(txn.DB)
	if opt != nil {
		builder = opt.CompileDelete(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}
`
