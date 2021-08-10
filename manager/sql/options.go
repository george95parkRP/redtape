package sql

type SqlManagerOptions struct {
	Dialect    string
	ConnString string
}

type SqlManagerOption func(*SqlManagerOptions)

// Returns a new SqlManageOptions object after applying options.
func NewSqlManagerOptions(opts ...SqlManagerOption) SqlManagerOptions {
	options := SqlManagerOptions{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// SetDialect sets the SqlManager dialect option.
func SetDialect(d string) SqlManagerOption {
	return func(o *SqlManagerOptions) {
		o.Dialect = d
	}
}

// SetDialect sets the SqlManager connstring option.
func SetConnString(conn string) SqlManagerOption {
	return func(o *SqlManagerOptions) {
		o.ConnString = conn
	}
}
