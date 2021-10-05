package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Member struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ChatRoom struct {
	RoomName  string `json:"room_name"`
	RoomOwner string `json:"room_owner"`
}

type render_data struct {
	Username    string
	ChatRooms   []ChatRoom
	Pages       []int
	CurrentPage int
}

// 建立資料表 member
func CreateUserTable(db *sql.DB) {
	sql := `create table member(
		id int(20) auto_increment primary key not null,
		username char(20) not null,
		password char(20) not null,
		email char(20 not null)
		); `

	if _, err := db.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}

func CreateRoomTable(db *sql.DB) {
	sql := `create table my_chatroom(
		room_name  char(20)  not null,
		room_owner char(20)  not null
		); `

	if _, err := db.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}

// 註冊會員
func CreateMember(db *sql.DB, m Member) bool {
	if CompareUserid(db, m.Username) == true && CompareEmail(db, m.Email) == true {

		stmt, err := db.Prepare("INSERT `member` SET id=?, username=?, email=?, password=?")
		checkErr(err)

		res, err := stmt.Exec(m.Id, m.Username, m.Email, m.Password)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		if err != nil {
			fmt.Println("create Member failed:", err)
			return false
		}
		fmt.Println("新增成功 : ", id)
		return true
	} else {
		fmt.Println("新增會員失敗")
		return false
	}
}

// 刪除會員
func DeleteMember(db *sql.DB, name string) {
	stmt, err := db.Prepare("delete from `member` where Binary username=?")
	checkErr(err)

	res, err := stmt.Exec(name)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("刪除成功 :", affect)
}

func UpdateMember(db *sql.DB, new string, old string) {
	stmt, err := db.Prepare("update `member` set Binary username=? where Binary username=?")
	checkErr(err)

	res, err := stmt.Exec(new, old)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("更新成功 : ", affect)
}

//	取得所有會員資訊
func GetMember(db *sql.DB) ([]Member, error) {
	rows, err := db.Query("select * from `member` ")
	if err != nil {
		fmt.Printf("Query failed,err:%v\n", err)
		return nil, err
	}
	m := Member{}
	members := []Member{}
	//一筆一筆讀取
	for rows.Next() {
		err = rows.Scan(&m.Id, &m.Username, &m.Password, &m.Email)
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return nil, err
		}
		members = append(members, m)
		defer rows.Close()
		if err != nil {
			return nil, err
		}
	}
	// fmt.Println(members)
	return members, nil
}

func CompareUserid(db *sql.DB, username string) bool {
	m := Member{}
	row := db.QueryRow("select username from `member` where Binary username=? limit 1", username)
	if err := row.Scan(&m.Username); err != nil {
		// fmt.Printf("scan failed, err : %v\n", err)
		fmt.Println("err : ", err)
		return true
	} else {
		fmt.Println("Username已存在")
		return false
	}
}
func CompareEmail(db *sql.DB, email string) bool {
	m := Member{}
	row := db.QueryRow("select email from `member` where Binary email=? limit 1", email)
	if err := row.Scan(&m.Email); err != nil {
		// fmt.Printf("scan failed, err : %v\n", err)
		fmt.Println("err : ", err)
		return true
	} else {
		// fmt.Println("E-mail已存在")
		return false
	}
}

func CompareLogin(db *sql.DB, email, password string) bool {
	m := Member{}
	if CompareEmail(db, email) == false {

		row := db.QueryRow("select password from `member` where Binary email=? limit 1", email)
		if err := row.Scan(&m.Password); err != nil {
			// fmt.Printf("scan failed, err : %v\n", err)
			fmt.Println("err   : ", err)
			return false
		} else if password != m.Password {

			fmt.Println("密碼錯誤 : ", password)
			return false
		} else {
			fmt.Println("登入成功")
			return true
		}

	} else {
		fmt.Println("尚未註冊Email")
		return false
	}
}

func From_Email_GetUserName(db *sql.DB, email string) string {
	m := Member{}
	row := db.QueryRow("select username from `member` where Binary email=?", email)
	if err := row.Scan(&m.Username); err != nil {
		// fmt.Printf("scan failed, err : %v\n", err)
		fmt.Println("err : ", err)
	}
	// fmt.Println("username : ", m.Username)
	return m.Username
}

//	取得所有聊天室
func GetMyChatroom(db *sql.DB, email string) []ChatRoom {
	rows, err := db.Query("select * from my_chatroom")
	if err != nil {
		fmt.Printf("Query failed,err:%v\n", err)
		return nil
	}
	r := ChatRoom{}
	rooms := []ChatRoom{}
	//一筆一筆讀取
	for rows.Next() {
		err = rows.Scan(&r.RoomName, &r.RoomOwner)
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return nil
		}
		rooms = append(rooms, r)
		defer rows.Close()
		if err != nil {
			return nil
		}
	}
	// fmt.Println(members)
	return rooms
}

//	新增聊天室
func CreateRoom(db *sql.DB, r ChatRoom) bool {

	// c := ChatRoom{}
	row := db.QueryRow("select room_name from my_chatroom where Binary room_name=? limit 1", r.RoomName)
	if err := row.Scan(&r.RoomName); err != nil {
		// fmt.Printf("scan failed, err : %v\n", err)
		stmt, err := db.Prepare("INSERT my_chatroom SET room_name=?, room_owner=?")
		checkErr(err)
		res, err := stmt.Exec(r.RoomName, r.RoomOwner)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)

		if err != nil {
			fmt.Println("create ChatRoom failed:", err)
			return false
		} else {
			fmt.Println("新增成功 : ", id)
			return true
		}
	} else {
		fmt.Println("聊天室名稱已被使用")
		return false
	}

}

//	取得所有聊天室 名稱與室長
func GetALLChatroom(db *sql.DB) ([]ChatRoom, error) {
	rows, err := db.Query("select * from my_chatroom ")
	if err != nil {
		fmt.Printf("Query failed,err : %v \n", err)
		return nil, err
	}
	c := ChatRoom{}
	all_chatroom := []ChatRoom{}
	//一筆一筆讀取
	for rows.Next() {
		err = rows.Scan(&c.RoomName, &c.RoomOwner)
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return nil, err
		}
		all_chatroom = append(all_chatroom, c)
		defer rows.Close()
		if err != nil {
			return nil, err
		}
	}
	// fmt.Println(members)
	return all_chatroom, nil
}

// 顯示下拉式選單的頁面 0~10 10~20 20~30 30~40
func GetSelectPage(db *sql.DB, page int) ([]ChatRoom, error) {
	rows, err := db.Query("select * from my_chatroom limit ?,10", (page*10)-10)
	if err != nil {
		fmt.Printf("Query failed,err : %v \n", err)
		return nil, err
	}
	c := ChatRoom{}
	all_chatroom := []ChatRoom{}
	//一筆一筆讀取
	for rows.Next() {
		err = rows.Scan(&c.RoomName, &c.RoomOwner)
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return nil, err
		}
		all_chatroom = append(all_chatroom, c)
		defer rows.Close()
		if err != nil {
			return nil, err
		}
	}
	// fmt.Println(members)
	return all_chatroom, nil
}

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

//----------------------------------------------------------------------
//    *****  Server *****

// var addr = flag.String("addr", ":5000", "http service address")

//----------------------------------------------------------------------
//    *****  CLient *****
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	data *Data
}

type Data struct {
	User     string   `json:"user"`
	UserList []string `json:"user_list"`
	RoomName string   `json:"room_name"`
}

type send_msg struct {
	Msg      string `json:"msg"`
	RoomName string `json:"room_name"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	// 用戶離開後註銷用戶
	defer func() {
		// c.hub.broadcast <- []byte("******* " + c.user + " 離開聊天室 *******") //
		// fmt.Println("所有用戶: ", user_list)
		// fmt.Println("離開用戶  :", c.data.User)
		// if len(user_list) == 1 && user_list[0] == c.data.User {
		// 	user_list = []string{}
		// } else if user_list[len(user_list)-1] == c.data.User {
		// 	user_list = user_list[:len(user_list)-1]
		// } else if user_list[0] == c.data.User {
		// 	user_list = user_list[1:]
		// } else {
		// 	// fmt.Println("ALL user: ", user_list)
		// 	for i := 0; i < len(user_list); i++ {
		// 		// fmt.Printf("user_list[%d] = %s \n", i, user_list[i])
		// 		if user_list[i] == c.data.User {
		// 			user_list = append(user_list[:i], user_list[i+1:]...)
		// 		}
		// 	}
		// }
		// fmt.Println("在線人員: ", user_list)
		// fmt.Println("*--------------------------------------------- ")
		// js_data, _ := json.Marshal(c.data)
		// fmt.Println("js_data : ", string(js_data))

		// c.hub.broadcast <- js_data
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// fmt.Println("登入人員 :", user_list)
		// c.data.UserList = user_list
		// js_data, _ := json.Marshal(c.data)
		// c.hub.broadcast <- js_data

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// fmt.Println("message : ", string(message))
		new_msg := send_msg{
			Msg:      string(message),
			RoomName: c.data.RoomName,
		}
		js_msg, _ := json.Marshal(new_msg)

		c.hub.broadcast <- js_msg

		// c.hub.broadcast <- message

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
// 	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, c echo.Context) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	sess, _ := session.Get("User", c)
	var username string
	var roomname string

	for k, v := range sess.Values {
		if k == "username" {
			username = v.(string)
		}
		if k == "RoomName" {
			roomname = v.(string)
		}
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), data: &Data{RoomName: roomname, User: username}}
	client.hub.register <- client
	// client.userlist = append(client.userlist,[]byte(username))

	// All_RoomUser := append(client.userlist, username)
	// var All_RoomUser string
	// fmt.Println("client.userlist", All_RoomUser)
	// for i, v := range client.user {
	// 	fmt.Println("i", i)
	// 	fmt.Println("v", v)

	// }
	// client.hub.broadcast <- []byte("聊天室成員有 :" + All_RoomUser)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

//-----------------------------------------------------------------------
//    *****  HUB *****
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			// user_list = append(user_list, client.data.User)
			// client.data.UserList = user_list
			// js_data, _ := json.Marshal(client.data)
			// fmt.Println(client.data)
			// client.send <- js_data
			var user_list = []string{}
			type user_list_json struct {
				User_list []byte `json:"user_list"`
				Roomname  string `json:"roomname"`
			}
			// js := user_list_json{}
			js_data, _ := json.Marshal(user_list)

			for c := range h.clients {
				fmt.Println("現在房間 :", client.data.RoomName)
				fmt.Println("所有的房間 :", c.data.RoomName)
				if client.data.RoomName == c.data.RoomName {
					user_list = append(user_list, c.data.User)
					js_data, _ = json.Marshal(user_list)
					fmt.Printf("進入js_data : %T  \n", js_data)
					fmt.Println("進入js_data value :", string(js_data))

					// js = user_list_json{User_list: js_data, Roomname: client.data.RoomName}
				} else {
					js_data, _ = json.Marshal(string(client.data.User))
					// js = user_list_json{User_list: js_data, Roomname: client.data.RoomName}
					fmt.Printf("進入js_data : %T  \n", js_data)
					fmt.Println("進入js_data value :", string(js_data))

					c.send <- js_data
				}
				// c.send <- js_data
			}

			for c := range h.clients {
				if client.data.RoomName == c.data.RoomName {
					c.send <- js_data
				}
			}

			fmt.Println("-----------------------------------------")

			// client.send <- js_data
			// fmt.Println("js_data : ", string(js_data))

		// 判斷用戶列表中是否存在此用戶 ， 是 -> 註銷
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

			var user_list = []string{}
			type user_list_json struct {
				User_list []byte `json:"user_list"`
				Roomname  string `json:"roomname"`
				AAA       []string
			}

			// js := user_list_json{}
			js_data, _ := json.Marshal(user_list)

			for c := range h.clients {
				fmt.Println("現在房間 :", client.data.RoomName)
				fmt.Println("所有的房間 :", c.data.RoomName)
				if client.data.RoomName == c.data.RoomName {
					user_list = append(user_list, c.data.User)
					js_data, _ = json.Marshal(user_list)
					fmt.Printf("離開 js_data : %T  \n", js_data)
					fmt.Println("離開 js_data value :", string(js_data))

					// js = user_list_json{User_list: js_data, Roomname: client.data.RoomName}
				} else {
					js_data, _ = json.Marshal(client.data.User)
					// js = user_list_json{User_list: js_data, Roomname: client.data.RoomName}
					fmt.Printf("進入js_data : %T  \n", js_data)
					fmt.Println("進入js_data value :", string(js_data))

					c.send <- js_data
				}
				// c.send <- js_data
			}

			for c := range h.clients {
				if client.data.RoomName == c.data.RoomName {
					c.send <- js_data
				}
			}
		// 取得message 並發發送給所有 Client
		case message := <-h.broadcast:
			// fmt.Println("broadcast  message :", string(message))
			var send_msg_json send_msg
			json.Unmarshal(message, &send_msg_json)
			// var user_list_json Data
			// json.Unmarshal(message, &user_list_json)
			// fmt.Println("js_message :", user_list_json.RoomName)

			for client := range h.clients {
				if send_msg_json.RoomName == client.data.RoomName {
					select {
					// 發送訊息
					case client.send <- []byte(send_msg_json.Msg):
					// 發送訊息失敗則刪除connection資訊
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}

			}
		}
	}
}

//-----------------------------------------------------------------------

func main() {
	//	連線DB
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chatroom?charset=utf8")
	checkErr(err)

	// CreateRoomTable(db)
	// CreateTable(db)
	// m := Member{
	// 	Id:       0,
	// 	Username: "Bbbcc",
	// 	Password: "AAaa1234567",
	// 	Email:    "Aaa@aaa",
	// }
	// CompareUserid(db, mike.Username)
	// CompareLogin(db, m.Email, m.Password)
	// CreateMember(db, mike)
	// DeleteMember(db, "mike")
	// UpdateMember(db, "Mike", "mike") // new,old

	// c := ChatRoom{
	// 	RoomName:  "jayroom",
	// 	RoomOwner: "jay",
	// }
	// CreateRoom(db, c)

	// var user_list = []string{}

	// user_list = append(user_list, "123")
	// fmt.Println("user_list : ", user_list)
	// fmt.Printf("user_list_TYPE %T: \n", user_list)

	// fmt.Println("----------------------------------------------------------------")
	// js, _ := json.Marshal(user_list)
	// fmt.Println("js : ", string(js))
	// fmt.Printf("js_TYPE %T: \n", js)

	// fmt.Println("----------------------------------------------------------------")
	// jss, _ := json.Marshal("123")
	// fmt.Println("jss : ", string(jss))
	// fmt.Printf("jss_TYPE %T: \n", jss)

	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Print(err)

	}

	// ---------------------------------------------------------------------------------------------------
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(32))))
	// e.Use(session.Middleware(sessions.NewFilesystemStore("./", securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))))

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer
	//	首頁
	e.Static("/home", "templates/home.html")
	e.Static("/chatroom", "templates/chatroom.html")
	e.Static("/singup", "templates/singup.html")
	e.Static("/create_chatroom", "templates/create_chatroom.html")

	//	會員登入
	e.POST("/login", func(c echo.Context) error {
		email := strings.TrimSpace(c.FormValue("email"))
		password := strings.TrimSpace(c.FormValue("password"))
		if CompareLogin(db, email, password) == true {
			sess, _ := session.Get("User", c)
			sess.Options = &sessions.Options{
				Path:   "/",       //	所有頁面都可以訪問session資料
				MaxAge: 86400 * 7, //	Session有效期(秒)，
			}

			sess.Values["email"] = email
			sess.Values["isLogin"] = true
			sess.Values["username"] = From_Email_GetUserName(db, email)
			sess.Values["current_page"] = 1
			sess.Save(c.Request(), c.Response()) //	保存使用者Session
			// fmt.Println("email :", email)
			// fmt.Println("password :", password)
			// return c.HTML(http.StatusOK, fmt.Sprintf("<p><h2>Login success</h2> <br>email : %s <br> password : %s</p>", email, password))
			// return c.Render(http.StatusOK, "my_chatroom.html", "")
			return c.Redirect(http.StatusFound, "/my_chatroom")

		} else {
			return c.HTML(403, fmt.Sprintf("<p><h2>登入失敗</h2></p>"))
		}
	})
	// 會員登出
	e.POST("/logout", func(c echo.Context) error {
		sess, _ := session.Get("User", c)
		sess.Options = &sessions.Options{
			Path:   "/", //	所有頁面都可以訪問session資料
			MaxAge: -1,  //	Session有效期(秒)，
		}
		sess.Values["isLogin"] = nil
		sess.Save(c.Request(), c.Response()) //	保存使用者Session
		fmt.Println("登出 : ", sess.Values["email"])
		return c.Redirect(http.StatusFound, "/home")
	})

	// 註冊會員
	e.POST("/Singup", func(c echo.Context) error {
		email := strings.TrimSpace(c.FormValue("email"))
		password := strings.TrimSpace(c.FormValue("password"))
		re_password := strings.TrimSpace(c.FormValue("re_password"))
		username := strings.TrimSpace(c.FormValue("username"))
		if password != re_password {
			return c.HTML(403, fmt.Sprintf("<p><h2>二次密碼輸入錯誤</h2> password : %s<br> re_password : %s</p>", password, re_password))
		}

		new_member := Member{
			Id:       0,
			Username: username,
			Password: password,
			Email:    email,
		}
		fmt.Println("------------------------------------")
		fmt.Println("email :", new_member.Email)
		fmt.Println("password :", new_member.Password)
		fmt.Println("id :", new_member.Id)
		fmt.Println("Nickname :", new_member.Username)
		if CreateMember(db, new_member) == true {
			// return c.HTML(http.StatusOK, fmt.Sprintf("<p><h2>註冊成功</h2> <br>email : %s <br> password : %s<br> re_password : %s<br> Nickname : %s</p>", email, password, re_password, username))
			// return c.Render(http.StatusOK, "home.html", "")
			return c.Redirect(http.StatusFound, "/home")
		} else {
			return c.HTML(http.StatusFound, fmt.Sprintf("<p><h2>註冊失敗</h2></p>"))
		}

		// return c.String(http.StatusOK, "email %s& password : "+email+"\n"+password)

	})

	// 聊天室列表
	e.GET("/my_chatroom", func(c echo.Context) error {
		sess, err := session.Get("User", c)
		all_chatroom, _ := GetALLChatroom(db)
		p := strings.TrimSpace(c.FormValue("page"))
		// fmt.Println("p : ", p)
		var selected_page int
		selected_page, _ = strconv.Atoi(p)
		// fmt.Println("p : ", selected_page)
		selected_page_chatroom := []ChatRoom{}
		if selected_page == 0 {
			selected_page = 1
		}
		selected_page_chatroom, _ = GetSelectPage(db, selected_page)

		var username string
		for k, v := range sess.Values {
			if k == "username" {
				// fmt.Println("k : ", k)
				username = v.(string)
				// fmt.Println("room_owner : ", username)
			}

		}

		var page = 1
		new_datas := len(all_chatroom)
		var pages []int
		if new_datas/10 > 0 {
			page += new_datas / 10
			if new_datas%10 == 0 {
				page--
			}
		}

		for i := 1; i <= page; i++ {
			pages = append(pages, i)
		}
		// fmt.Println("pages :", pages)

		new_data := render_data{
			Username:    username,
			ChatRooms:   selected_page_chatroom,
			Pages:       pages,
			CurrentPage: selected_page,
		}
		// fmt.Println("CurrentPage :", sess.Values["current_page"])
		// fmt.Println("len of new_data.ChatRooms : ", new_data.ChatRooms)
		// fmt.Println("pages : ", pages)

		if err != nil {
			return err
		}
		if sess.Values["isLogin"] == true {
			// fmt.Println("存取成功 : ", sess.Values["username"])
			return c.Render(http.StatusOK, "my_chatroom", new_data)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")

		}
	})

	//	首頁
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/home")
	})

	// 新增聊天室
	e.POST("/create_chatroom", func(c echo.Context) error {
		sess, err := session.Get("User", c)
		if err != nil {
			return err
		}
		// all_chatroom, _ := GetALLChatroom(db)

		if sess.Values["isLogin"] == true {
			// fmt.Println("存取成功 : ", sess.Values["username"])
			room_name := strings.TrimSpace(c.FormValue("room_name"))
			var room_owner string
			for k, v := range sess.Values {
				if k == "username" {
					fmt.Println("k : ", k)
					room_owner = v.(string)
					fmt.Println("room_owner : ", room_owner)
				}

			}
			// fmt.Println("room_owner : ", room_owner)
			// room_owner = sess.Values["room_owner"]
			new_room := ChatRoom{
				RoomName:  room_name,
				RoomOwner: room_owner,
			}
			if CreateRoom(db, new_room) == true {
				return c.Redirect(http.StatusFound, "/my_chatroom")

			} else {
				// fmt.Println("新增失敗")
				return c.HTML(http.StatusOK, fmt.Sprintf("<p><h2>新增失敗，聊天室名稱已被使用，請回上一頁重新新增</h2></p>"))
			}

		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")
		}
		// return nil
	})

	//點選下拉式選單更換頁面
	e.POST("/selected_page", func(c echo.Context) error {
		p := strings.TrimSpace(c.FormValue("page")) // p = 下拉式選單所選擇的頁數
		// fmt.Println("p : ", p)
		selected_page, _ := strconv.Atoi(p)
		all_chatroom, _ := GetALLChatroom(db)
		selected_page_chatroom, _ := GetSelectPage(db, selected_page)

		// fmt.Println("all_chatroom : ", all_chatroom)

		sess, err := session.Get("User", c)
		if err != nil {
			return err
		}

		sess.Values["current_page"] = selected_page
		sess.Save(c.Request(), c.Response()) //	保存使用者Session
		// fmt.Println("current_page : ", sess.Values["current_page"])
		var username string
		for k, v := range sess.Values {
			if k == "username" {
				// fmt.Println("k : ", k)
				username = v.(string)
				// fmt.Println("room_owner : ", username)
			}

		}

		var page = 1
		new_datas := len(all_chatroom)
		var pages []int
		if new_datas/10 > 0 {
			page += new_datas / 10
			if new_datas%10 == 0 {
				page--
			}

		}

		for i := 1; i <= page; i++ {
			pages = append(pages, i)
		}
		// fmt.Println("pages :", pages)
		new_data := render_data{
			Username:    username,
			ChatRooms:   selected_page_chatroom,
			Pages:       pages,
			CurrentPage: selected_page,
		}
		if sess.Values["isLogin"] == true {
			// fmt.Println("存取成功 : ", sess.Values["username"])
			return c.Render(http.StatusOK, "my_chatroom", new_data)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")

		}
	})

	// 點選prev按鈕，回到上一頁聊天室列表
	e.GET("/prev_page", func(c echo.Context) error {
		sess, err := session.Get("User", c)
		var current_page int
		var prev_page int
		for k, v := range sess.Values {
			if k == "current_page" {
				// fmt.Println("cp ,v : ", k, v)
				current_page = v.(int)
			}

		}
		if current_page > 1 {
			prev_page = current_page - 1
		} else {
			prev_page = current_page
		}
		// fmt.Println("prev_page : ", prev_page)

		// current_page, _ := strconv.Atoi(cp)
		// fmt.Println("current_page : ", current_page)
		// prev_page := current_page - 1
		// fmt.Println("prev_page : ", prev_page)
		sess.Values["current_page"] = prev_page
		sess.Save(c.Request(), c.Response()) //	保存使用者Session
		all_chatroom, _ := GetALLChatroom(db)
		selected_page_chatroom, _ := GetSelectPage(db, prev_page)

		// // fmt.Println("all_chatroom : ", all_chatroom)

		if err != nil {
			return err
		}
		// // var chatroom []Chatroom
		var username string
		for k, v := range sess.Values {
			if k == "username" {
				username = v.(string)
			}
		}

		var page = 1
		new_datas := len(all_chatroom)
		var pages []int
		if new_datas/10 > 0 {
			page += new_datas / 10
			if new_datas%10 == 0 {
				page--
			}

		}

		for i := 1; i <= page; i++ {
			pages = append(pages, i)
		}
		// fmt.Println("pages :", pages)
		new_data := render_data{
			Username:    username,
			ChatRooms:   selected_page_chatroom,
			Pages:       pages,
			CurrentPage: prev_page,
		}
		if sess.Values["isLogin"] == true {
			// fmt.Println("存取成功 : ", sess.Values["username"])
			return c.Render(http.StatusOK, "my_chatroom", new_data)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")

		}
	})

	// 點選next按鈕，進入下一頁聊天室列表
	e.GET("/next_page", func(c echo.Context) error {
		sess, err := session.Get("User", c)
		if err != nil {
			return err
		}

		all_chatroom, _ := GetALLChatroom(db)
		var username string
		for k, v := range sess.Values {
			if k == "username" {
				username = v.(string)
			}
		}

		var page = 1
		new_datas := len(all_chatroom)
		var pages []int
		if new_datas/10 > 0 {
			page += new_datas / 10
			if new_datas%10 == 0 {
				page--
			}

		}

		for i := 1; i <= page; i++ {
			pages = append(pages, i)
		}
		var current_page int
		var next_page int
		for k, v := range sess.Values {
			if k == "current_page" {
				// fmt.Println("cp ,v : ", k, v)
				current_page = v.(int)
			}

		}
		if current_page < page {
			next_page = current_page + 1
		} else {
			next_page = current_page
		}
		// fmt.Println("next_page : ", next_page)

		// current_page, _ := strconv.Atoi(cp)
		// fmt.Println("current_page : ", current_page)
		// prev_page := current_page - 1
		// fmt.Println("prev_page : ", prev_page)
		sess.Values["current_page"] = next_page
		sess.Save(c.Request(), c.Response()) //	保存使用者Session
		selected_page_chatroom, _ := GetSelectPage(db, next_page)

		// // fmt.Println("all_chatroom : ", all_chatroom)

		// // var chatroom []Chatroom

		// fmt.Println("pages :", pages)
		new_data := render_data{
			Username:    username,
			ChatRooms:   selected_page_chatroom,
			Pages:       pages,
			CurrentPage: next_page,
		}
		if sess.Values["isLogin"] == true {
			// fmt.Println("存取成功 : ", sess.Values["username"])
			return c.Render(http.StatusOK, "my_chatroom", new_data)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")
		}
	})

	// // 刪除會員
	// e.GET("/users/del/:id", func(c echo.Context) error {
	// 	id := c.Param("id")
	// 	DeleteMember(db, id)
	// 	return c.String(http.StatusOK, "刪除成功 : "+id)
	// })

	// // 顯示所有會員
	// e.GET("/users/show", func(c echo.Context) error {
	// 	member, err := GetMember(db)
	// 	js, err := json.MarshalIndent(member, "", "")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(string(js))
	// 	return c.String(http.StatusOK, string(js))
	// })

	//----------------------------------------------------------------
	flag.Parse()
	hub := newHub()
	go hub.run()

	e.GET("/chatroom", func(c echo.Context) error {
		println("connection successs")
		sess, err := session.Get("User", c)
		if err != nil {
			panic(err)
		}

		var username string
		for k, v := range sess.Values {
			if k == "username" {
				username = v.(string)
			}
		}

		type ChatRoom struct {
			RoomName  string `json:"room_name"`
			RoomOwner string `json:"room_owner"`
			LoginUser string `json:"loom_owner"`
		}

		send_data := ChatRoom{
			RoomName:  "Room1",
			RoomOwner: "Miles",
			LoginUser: username,
		}

		if sess.Values["isLogin"] == true {
			// serveWs(hub, c)
			return c.Render(http.StatusFound, "chatroom", send_data)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")
		}

	})

	// 透過Websocket連線
	e.GET("/chatroom/ws", func(c echo.Context) error {
		println("ws connection")
		serveWs(hub, c)
		return nil
	})

	// 點選聊天室，傳送被點選房間之名稱
	e.POST("/chatroom", func(c echo.Context) error {
		// println("POST　ＯＫ")
		sess, err := session.Get("User", c)
		if err != nil {
			panic(err)
		}
		selected_room_name := strings.TrimSpace(c.FormValue("selected_room_name"))   // p = 下拉式選單所選擇的頁數
		selected_room_owner := strings.TrimSpace(c.FormValue("selected_room_owner")) // p = 下拉式選單所選擇的頁數

		var username string
		for k, v := range sess.Values {
			if k == "username" {
				username = v.(string)
			}
		}

		type ChatRoom struct {
			RoomName  string `json:"room_name"`
			RoomOwner string `json:"room_owner"`
			LoginUser string `json:"loom_owner"`
		}

		send_data := ChatRoom{
			RoomName:  selected_room_name,
			RoomOwner: selected_room_owner,
			LoginUser: username,
		}

		sess.Values["RoomName"] = selected_room_name
		sess.Save(c.Request(), c.Response()) //	保存使用者Session

		if sess.Values["isLogin"] == true {
			// serveWs(hub, c)
			return c.Render(http.StatusFound, "chatroom", send_data)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")
		}
	})

	e.Logger.Fatal(e.Start(":5000"))

}
