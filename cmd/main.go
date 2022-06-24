package main

import (
	"context"
	"db_project/app/handlers"
	"db_project/app/repositories"
	"db_project/app/usecases"
	"db_project/utils/queryCheck"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Urls struct {
	Root    string
	User    string
	Forum   string
	Thread  string
	Service string
	Post    string
}

func GetUrls() Urls {
	return Urls{
		Root:    "/api",
		User:    "/user",
		Forum:   "/forum",
		Thread:  "/thread",
		Service: "/service",
		Post:    "/post",
	}
}

type Repositories struct {
	User    repositories.IUserRepository
	Forum   repositories.IForumRepository
	Thread  repositories.IThreadRepository
	Service repositories.IServiceRepository
	Post    repositories.IPostRepository
}

type UseCases struct {
	User    usecases.IUserUseCase
	Forum   usecases.IForumUseCase
	Thread  usecases.IThreadUseCase
	Service usecases.IServiceUseCase
	Post    usecases.IPostUseCase
}

func main() {
	APIPort := "5000"
	DSN := "host=localhost port=5432 user=forum_user password=forum_user_password dbname=forum sslmode=disable"
	Urls := GetUrls()
	APIAddr := fmt.Sprintf("0.0.0.0:%v", APIPort)
	Repositories := Repositories{}
	UseCases := UseCases{}

	_, err := queryCheck.GetInstance()
	if err != nil {
		fmt.Printf("Can't create queryCheck instance: %v", err)
		return
	}

	config, err := pgxpool.ParseConfig(DSN)
	config.MaxConns = 2000
	if err != nil {
		fmt.Printf("Can't parese DSN: %v\n", err)
		return
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Printf("Can't create DB connection pool: %v", err)
		return
	}

	gin.SetMode("release")
	router := gin.New()
	router.Use(gin.Recovery())
	apiGroup := router.Group(Urls.Root)

	Repositories.User = repositories.CreateUserRepository(db)
	Repositories.Thread = repositories.CreateThreadRepository(db)
	Repositories.Forum = repositories.CreateForumRepository(db)
	Repositories.Service = repositories.CreateServiceRepository(db)
	Repositories.Post = repositories.CreatePostRepository(db)

	UseCases.User = usecases.CreateUserUseCase(Repositories.User)
	UseCases.Thread = usecases.CreateThreadUseCase(Repositories.Thread)
	UseCases.Forum = usecases.CreateForumUseCase(Repositories.Forum, Repositories.Thread)
	UseCases.Service = usecases.CreateServiceUseCase(Repositories.Service)
	UseCases.Post = usecases.CreatePostUseCase(Repositories.Post, UseCases.Forum, UseCases.User, UseCases.Thread)

	userHandler := handlers.MakeUsersHandler(UseCases.User)
	userRouter := apiGroup.Group(Urls.User)
	userRouter.GET("/:nickname/profile", userHandler.Get)
	userRouter.POST("/:nickname/profile", userHandler.Update)
	userRouter.POST("/:nickname/create", userHandler.Create)

	forumHandler := handlers.MakeForumsHandler(UseCases.Forum)
	forumRouter := apiGroup.Group(Urls.Forum)
	forumRouter.GET("/:slug/details", forumHandler.Get)
	forumRouter.POST("/create", forumHandler.Create)
	forumRouter.GET("/:slug/users", forumHandler.GetUsers)
	forumRouter.GET("/:slug/threads", forumHandler.GetThreads)
	forumRouter.POST("/:slug/create", forumHandler.CreateThread)

	threadHandler := handlers.MakeThreadsHandler(UseCases.Thread)
	threadRouter := apiGroup.Group(Urls.Thread)
	threadRouter.GET("/:slug_or_id/details", threadHandler.Get)
	threadRouter.POST("/:slug_or_id/details", threadHandler.Update)
	threadRouter.POST("/:slug_or_id/vote", threadHandler.Vote)
	threadRouter.POST("/:slug_or_id/create", threadHandler.PostsCreate)
	threadRouter.GET("/:slug_or_id/posts", threadHandler.GetPosts)

	serviceHandler := handlers.MakeServicesHandler(UseCases.Service)
	serviceRouter := apiGroup.Group(Urls.Service)
	serviceRouter.POST("/clear", serviceHandler.Clear)
	serviceRouter.GET("/status", serviceHandler.Status)

	postHandler := handlers.MakePostsHandler(UseCases.Post)
	postRouter := apiGroup.Group(Urls.Post)
	postRouter.GET("/:id/details", postHandler.Get)
	postRouter.POST("/:id/details", postHandler.Update)

	err = router.Run(APIAddr)
	if err != nil {
		fmt.Printf("Can't start server: %v\n", err)
		return
	}
}
