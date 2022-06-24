package constants

const (
	PostUser   string = "user"
	PostForum  string = "forum"
	PostThread string = "thread"
)

type SortType string

const (
	SortFlat       SortType = "flat"
	SortTree       SortType = "tree"
	SortParentTree SortType = "parent_tree"
)

var (
	DescSincePostQuery = map[SortType]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND id < $2 ORDER BY id DESC LIMIT $3",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2) ORDER BY path DESC LIMIT $3",
		SortParentTree: `
WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
      AND parent IS NULL
      AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
    ORDER BY path[1] DESC
    LIMIT $3
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts WHERE thread = $1 AND path[1] IN (SELECT * FROM roots) ORDER BY path[1] DESC, path[2:]`,
	}
	AscSincePostQuery = map[SortType]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND id > $2 ORDER BY id LIMIT $3",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2) " +
			"ORDER BY path LIMIT $3",
		SortParentTree: `
WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
      AND parent IS NULL
      AND path[1] > (SELECT path[1] FROM posts WHERE id = $2)
    ORDER BY path[1]
    LIMIT $3
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts WHERE thread = $1 AND path[1] IN (SELECT * FROM roots) ORDER BY path`,
	}
	DescNoSincePostQuery = map[SortType]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY id DESC LIMIT $2",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY path DESC LIMIT $2",
		SortParentTree: `
WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
    ORDER BY path[1] DESC
    LIMIT $2
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts WHERE thread = $1 AND path[1] IN (SELECT * FROM roots) ORDER BY path[1] DESC, path[2:]`,
	}
	AscNoSincePostQuery = map[SortType]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY id LIMIT $2",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY path LIMIT $2\n",
		SortParentTree: `WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
    ORDER BY path[1]
    LIMIT $2
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts
WHERE thread = $1 AND path[1] IN (SELECT * FROM roots)
ORDER BY path`,
	}
)

var (
	ForumQuery = map[SortType]string{
		"Create":                `INSERT INTO forums ("user", slug, title) VALUES ((SELECT nickname FROM users WHERE nickname = $3), $1, $2) RETURNING slug, title, "user", posts, threads`,
		"Get":                   `SELECT id, slug, title, "user", posts, threads FROM forums WHERE slug = $1`,
		"GetThreadsDesc":        ` AND created <= $2 ORDER BY created DESC LIMIT $3`,
		"GetThreadsSinceDesc":   ` ORDER BY created DESC LIMIT $2`,
		"GetThreadsNoDesc":      ` AND created >= $2 ORDER BY created LIMIT $3`,
		"GetThreadsSinceNoDesc": ` ORDER BY created LIMIT $2`,
		"GetThreads":            `SELECT id, COALESCE(slug, ''), author, forum, title, message, created, votes FROM threads WHERE forum = $1`,
		"GetUsers":              `SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users AS fu JOIN users AS u ON fu.nickname = u.nickname WHERE fu.forum = $1 `,
		"GetUsersDesc":          `ORDER BY u.nickname DESC LIMIT $2`,
		"GetUsersSinceDesc":     `AND u.nickname < $2 ORDER BY u.nickname DESC LIMIT $3`,
		"GetUsersNoDesc":        `ORDER BY u.nickname LIMIT $2`,
		"GetUsersSinceNoDesc":   `AND u.nickname > $2 ORDER BY u.nickname LIMIT $3`,
		"CreateThread": `INSERT INTO threads (slug, author, forum, title, message, created) VALUES (NULLIF($1, ''), (SELECT nickname FROM users WHERE nickname = $2), 
		(SELECT slug FROM forums WHERE slug = $3), $4, $5, $6) RETURNING id, $1, author, forum, title, message, created, votes`,
	}
	PostQuery = map[SortType]string{
		"Get": `SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message FROM posts WHERE id = $1`,
		"Update": `UPDATE posts SET message = COALESCE(NULLIF($1, ''), message), 
		isEdited = CASE WHEN (isEdited = TRUE OR (isEdited = FALSE AND NULLIF($1, '') IS NOT NULL AND NULLIF($1, '') <> message)) 
		THEN TRUE ELSE FALSE END WHERE id = $2 
		RETURNING id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message`,
	}
	ServiceQuery = map[SortType]string{
		"Clear":        `TRUNCATE users, forums, threads, votes, posts, forum_users`,
		"queryUsers":   `SELECT COUNT(*) FROM users`,
		"queryForums":  `SELECT COUNT(*) FROM forums`,
		"queryThreads": `SELECT COUNT(*) FROM threads`,
		"queryPosts":   `SELECT COUNT(*) FROM posts`,
	}
	ThreadQuery = map[SortType]string{
		"GetBySlug":        `SELECT id, COALESCE(slug, ''), author, forum, title, message, created, votes FROM threads WHERE slug = $1`,
		"PostsCreate":      `INSERT INTO posts(parent, author, forum, thread, message, created) VALUES `,
		"CreatePostsTwo":   ` RETURNING id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message`,
		"VoteByID":         `INSERT INTO votes (nickname, thread, value) VALUES ($1, $2, $3) ON CONFLICT (nickname, thread) DO UPDATE SET value = $3`,
		"CreatePostsBatch": `INSERT INTO posts(parent, author, forum, thread, message, created) VALUES (NULLIF($1, 0), $2, $3, $4, $5, $6) RETURNING id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message`,
		"VoteBySlug": `INSERT INTO votes (nickname, thread, value) VALUES ($1, (SELECT id FROM threads WHERE slug=$2), $3) 
		ON CONFLICT (nickname, thread) DO UPDATE SET value = $3`,
		"UpdateByID": `UPDATE threads SET title = COALESCE(NULLIF($1, ''), title), 
		message = COALESCE(NULLIF($2, ''), message) WHERE id = $3 
		RETURNING id, COALESCE(slug, ''), author, forum, title, message, created, votes`,
		"GetByID": `SELECT id, COALESCE(slug, ''), author, forum, title, message, created, votes FROM threads WHERE id = $1`,
		"UpdateBySlug": `UPDATE threads SET title = COALESCE(NULLIF($1, ''), title), 
		message = COALESCE(NULLIF($2, ''), message) WHERE slug = $3 
		RETURNING id, slug, author, forum, title, message, created, votes`,
	}
	UserQuery = map[SortType]string{
		"Get":               `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1`,
		"Create":            `INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`,
		"GetUsersByUserNOE": `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2`,
		"Update": `UPDATE users SET fullname = COALESCE(NULLIF($1, ''), fullname), 
		about = COALESCE(NULLIF($2, ''), about), 
		email = COALESCE(NULLIF($3, ''), email) WHERE nickname = $4 
		RETURNING nickname, fullname, about, email`,
	}
)
