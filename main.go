package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
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
		room_name int(20) auto_increment primary key not null,
		room_owner char(20) not null
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

		stmt, err := db.Prepare("INSERT member SET id=?, username=?, email=?, password=?")
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
	stmt, err := db.Prepare("delete from member where Binary username=?")
	checkErr(err)

	res, err := stmt.Exec(name)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("刪除成功 :", affect)
}

func UpdateMember(db *sql.DB, new string, old string) {
	stmt, err := db.Prepare("update member set Binary username=? where Binary username=?")
	checkErr(err)

	res, err := stmt.Exec(new, old)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("更新成功 : ", affect)
}

//	取得所有會員資訊
func GetMember(db *sql.DB) ([]Member, error) {
	rows, err := db.Query("select * from member ")
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
	row := db.QueryRow("select username from member where Binary username=? limit 1", username)
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
	row := db.QueryRow("select email from member where Binary email=? limit 1", email)
	if err := row.Scan(&m.Email); err != nil {
		// fmt.Printf("scan failed, err : %v\n", err)
		fmt.Println("err : ", err)
		return true
	} else {
		fmt.Println("E-mail已存在")
		return false
	}
}

func CompareLogin(db *sql.DB, email, password string) bool {
	m := Member{}
	if CompareEmail(db, email) == false {

		row := db.QueryRow("select password from member where Binary email=? limit 1", email)
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
	row := db.QueryRow("select username from member where Binary email=?", email)
	if err := row.Scan(&m.Username); err != nil {
		// fmt.Printf("scan failed, err : %v\n", err)
		fmt.Println("err : ", err)
	}
	fmt.Println("username : ", m.Username)
	return m.Username
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

	// member, err := GetMember(db)
	// js, err := json.MarshalIndent(member, "", "")
	// fmt.Println("json := ", string(js))
	// fmt.Println("----------所有會員----------\n")
	// for _, m := range member {
	// 	fmt.Println(m)
	// }
	// fmt.Println("---------------------------\n")
	// fmt.Println(reflect.TypeOf(member))

	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Print(err)
	}

	// ---------------------------------------------------------------------------------------------------

	e := echo.New()
	// sessionPath := "."         //	設置session資料存放位置
	// sessionKey := "QSDFSfdsfs" // 設置cookies加密秘鑰
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(32))))
	// e.Use(session.Middleware(sessions.NewFilesystemStore("./", securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))))

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer
	//	首頁
	e.Static("/home", "templates/home.html")
	// e.Static("/my_chatroom", "templates/my_chatroom.html")
	e.Static("/singup", "templates/singup.html")

	//	會員登入
	e.POST("/login", func(c echo.Context) error {
		email := strings.TrimSpace(c.FormValue("email"))
		password := strings.TrimSpace(c.FormValue("password"))
		if CompareLogin(db, email, password) == true {
			sess, _ := session.Get("user_session", c)
			sess.Options = &sessions.Options{
				Path:   "/",       //	所有頁面都可以訪問session資料
				MaxAge: 86400 * 7, //	Session有效期(秒)，
			}

			sess.Values["email"] = email
			sess.Values["isLogin"] = true
			sess.Values["username"] = From_Email_GetUserName(db, email)
			sess.Save(c.Request(), c.Response()) //	保存使用者Session

			fmt.Println("email :", email)
			fmt.Println("password :", password)
			// return c.HTML(http.StatusOK, fmt.Sprintf("<p><h2>Login success</h2> <br>email : %s <br> password : %s</p>", email, password))
			// return c.Render(http.StatusOK, "my_chatroom.html", "")
			return c.Redirect(http.StatusFound, "/my_chatroom")

		} else {
			return c.HTML(403, fmt.Sprintf("<p><h2>登入失敗</h2></p>"))
		}
	})
	// 會員登出
	e.POST("/logout", func(c echo.Context) error {
		sess, _ := session.Get("user_session", c)
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
		fmt.Println("")
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

	e.GET("/my_chatroom", func(c echo.Context) error {
		sess, err := session.Get("user_session", c)
		if err != nil {
			return err
		}
		if sess.Values["isLogin"] == true {
			fmt.Println("存取成功 : ", sess.Values["username"])
			return c.Render(http.StatusOK, "my_chatroom", sess.Values)
		} else {
			fmt.Println("存取失敗，請先登入")
			return c.Redirect(http.StatusFound, "/home")

		}
	})
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/home")
	})

	// //	新增會員
	// e.GET("/users/add", func(c echo.Context) error {
	// 	id := c.QueryParam("username")
	// 	var m Member
	// 	m.Id = 0
	// 	m.Username = id
	// 	m.Password = id + "123"
	// 	m.Email = id + "@gmail.com"
	// 	CreateMember(db, m)
	// 	return c.String(http.StatusOK, "新增成功 id : "+id)
	// })
	// // 刪除會員
	// e.GET("/users/del/:id", func(c echo.Context) error {
	// 	id := c.Param("id")
	// 	DeleteMember(db, id)
	// 	return c.String(http.StatusOK, "刪除成功 : "+id)
	// })

	// // 顯示所有會員
	e.GET("/users/show", func(c echo.Context) error {
		member, err := GetMember(db)
		js, err := json.MarshalIndent(member, "", "")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(js))
		return c.String(http.StatusOK, string(js))
	})

	// e.Any("/users/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, c.QueryParam("pass"))
	// })
	e.Logger.Fatal(e.Start(":5000"))

}
