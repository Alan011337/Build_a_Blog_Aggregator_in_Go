package main


import _ "github.com/lib/pq"
import (
	"fmt"
	"github.com/Alan011337/internal/config"
	"github.com/Alan011337/internal/database"
	"github.com/google/uuid"
	"os"
	"database/sql"
	"context"
	"time"
	"net/http"
	"io"
	"encoding/xml"
	"html"
	"strconv"
)


type state struct {
    db  *database.Queries
    cfg *config.Config
}


type command struct {
	name string
	args []string
}


type commands struct {
	registeredCommands map[string]func(*state, command) error
}


type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}


type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func (c *commands) run(s *state, cmd command) error {
	_, ok := c.registeredCommands[cmd.name]
	if ok {
		err := c.registeredCommands[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("We can't find this command")
	}
	return nil
}


func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}


func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the command's arg's slice is empty")
	}

	if len(cmd.args) < 1 {
		return fmt.Errorf("the command's arg's slice is enough")
	}

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name: cmd.args[0],
		},
	)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("the user has been set.")
	fmt.Printf("User created successfully: %+v\n", user)

	return nil
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the command's arg's slice is empty")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])

	if err != nil {
		return err
	}

	fmt.Println("the user has been set.")

	return nil
}


func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("DeleteUsers Failed: %w", err)
	}

	fmt.Println("the users has been deleted.")

	return nil
}


func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("GetUsers Failed: %w", err)
	}

	fmt.Println("User has been get.")

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current) \n", user.Name)
		} else {
			fmt.Printf("* %v \n", user.Name)
		}
	}

	return nil

}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	feed := RSSFeed{}

	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, _:= range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}


func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the command's arg's slice is empty")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	
	return nil
}


func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("the command's arg's slice is not enough")
	}

	params := database.CreateFeedParams {
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)

	CreateFeedFollowParams := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), CreateFeedFollowParams)
	if err != nil {
		return err
	}

	return nil
}


func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* %s (%s) created by %s\n", feed.FeedName, feed.FeedUrl, feed.UserName)
	}

	return nil
}


func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the command's arg's slice is not enough")
	}

	Feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	CreateFeedFollowParams := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: Feed.ID,
	}

	CreateFeedFollowRow, err := s.db.CreateFeedFollow(context.Background(), CreateFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println(CreateFeedFollowRow.FeedName)
	fmt.Println(CreateFeedFollowRow.UserName)

	return nil
}


func handlerFollowing(s *state, cmd command, user database.User) error {
	GetFeedFollowsForUserRow, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range GetFeedFollowsForUserRow {
		fmt.Println(feed.FeedName)
	}

	return nil
}


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, cmd command) error {
        user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
        if err != nil {
            return err
        }
        return handler(s, cmd, user)
    }
}


func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the command's arg's slice is not enough")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	UnfollowFeedParams := database.UnfollowFeedParams {
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.UnfollowFeed(context.Background(), UnfollowFeedParams)
	if err != nil {
		return err
	}

	return nil
}


func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return nil
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return nil
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		publishAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}

		CreatePostParams := database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishAt,
			FeedID: feed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), CreatePostParams)
		if err != nil {
			return err
		}
	}


	return nil
}


func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) >= 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = parsedLimit
	}

	GetPostsParams := database.GetPostsParams {
		UserID: user.ID,
		Limit: int32(limit),
	}

	GetPostsRow, err := s.db.GetPosts(context.Background(), GetPostsParams)
	if err != nil {
		return err
	}

	for _, post := range GetPostsRow {
		fmt.Printf("%+v", post)
	}

	return nil
}


func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return 
	}

	db, err := sql.Open("postgres", cfg.DBURL) // 讀取連線網址並開啟連線
	if err != nil {
		return
	}
	dbQueries := database.New(db)             // 初始化 SQLC 查詢物件

	s := state {
		db: dbQueries,
		cfg: &cfg,
	}

	cmds := commands {
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
    cmds.register("follow", middlewareLoggedIn(handlerFollow))
    cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))


	if len(os.Args) < 2 {
		fmt.Println("the args are less than 2.")
		os.Exit(1) 
	} else {
		command_name := os.Args[1]
		command_parameter := os.Args[2:]
		cmd := command {
			name: command_name,
			args: command_parameter,
		}
		err := cmds.run(&s, cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}