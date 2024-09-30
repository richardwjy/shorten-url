package repo

type DbRepo struct {
	conn string
}

func NewDbRepo(cfg string) *DbRepo {
	return &DbRepo{
		conn: cfg,
	}
}
