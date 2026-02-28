package clause

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	OFFSET
	ORDERBY
	WHERE
)
