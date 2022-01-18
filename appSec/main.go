package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

type student struct {
	ID       string `json:"id"`
	Lastname string `json:"Lastname"`
	Name     string `json:"Name"`
	Filiere  string `json:"filiere"`
}

type teacher struct {
	ID       string `json:"id"`
	Lastname string `json:"Lastname"`
	Name     string `json:"Name"`
	Class    string `json:"class"`
}

var students = []student{
	{ID: "1", Lastname: "Montandon", Name: "Philippe", Filiere: "ISC"},
	{ID: "2", Lastname: "Chatillon", Name: "Thibault", Filiere: "ISC"},
}

var teachers = []teacher{
	{ID: "1", Lastname: "Pfeiffer", Name: "Ludovic", Class: "A401"},
	{ID: "2", Lastname: "Ouafi", Name: "Khaled", Class: "A402"},
}

// getStudents responds with the list of all students as JSON.
func getStudents(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, students)
}

// postStudents adds an student from JSON received in the request body.
func postStudents(c *gin.Context) {

	var newStudent student

	// Call BindJSON to bind the received JSON to
	// newStudent.
	if err := c.BindJSON(&newStudent); err != nil {
		fmt.Println("error post student json ?")
		return
	}

	// Add the new student to the slice.
	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
}

// getStudentByID locates the student whose ID value matches the id
// parameter sent by the client, then returns that student as a response.
func getStudentByID(c *gin.Context) {

	id := c.Param("id")

	// Loop through the list of students, looking for
	// an student whose ID value matches the parameter.
	for _, a := range students {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

func deleteStudentByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of students, looking for
	// an student whose ID value matches the parameter.
	for i := 0; i < len(students); i++ {
		if students[i].ID == id {
			c.IndentedJSON(http.StatusOK, students[i])
			copy(students[i:], students[i+1:])    // Shift a[i+1:] left one index.
			students = students[:len(students)-1] // Truncate slice.

			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

// getTeachers responds with the list of all teachers as JSON.
func getTeachers(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, teachers)
}

// postTeachers adds an teacher from JSON received in the request body.
func postTeachers(c *gin.Context) {

	var newTeacher teacher

	// Call BindJSON to bind the received JSON to
	// newTeacher.
	if err := c.BindJSON(&newTeacher); err != nil {
		return
	}

	// Add the new teacher to the slice.
	teachers = append(teachers, newTeacher)
	c.IndentedJSON(http.StatusCreated, newTeacher)
}

// getTeacherByID locates the teacher whose ID value matches the id
// parameter sent by the client, then returns that teacher as a response.
func getTeacherByID(c *gin.Context) {

	id := c.Param("id")

	// Loop through the list of teachers, looking for
	// an teacher whose ID value matches the parameter.
	for _, a := range teachers {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "teacher not found"})
}

func deleteTeacherByID(c *gin.Context) {

	id := c.Param("id")

	// Loop through the list of teachers, looking for
	// an teacher whose ID value matches the parameter.
	for i := 0; i < len(teachers); i++ {
		if teachers[i].ID == id {
			c.IndentedJSON(http.StatusOK, teachers[i])
			copy(teachers[i:], teachers[i+1:])    // Shift a[i+1:] left one index.
			teachers = teachers[:len(teachers)-1] // Truncate slice.

			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "teacher not found"})
}

func check_auth(c *gin.Context) bool {
	user, _, _ := c.Request.BasicAuth()
	if c.Request.Method != "GET" && user != "aristote" {
		return false
	}
	return true
}

var toValidate = map[string]string{
	"aud": "api://default",
	"cid": os.Getenv("0oa3lp6i6zXA2yMvp5d7"),
}

func check_teacher_authentication(c *gin.Context) {
	status := true
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		verifierSetup := jwtverifier.JwtVerifier{
			Issuer:           "https://" + os.Getenv("dev-58917141.okta.com") + "/oauth2/default",
			ClaimsToValidate: toValidate,
		}
		verifier := verifierSetup.New()
		_, err := verifier.VerifyAccessToken(token)
		if err != nil {
			c.String(http.StatusForbidden, err.Error())
			print(err.Error())
			status = false
		}
	} else {
		//c.String(http.StatusUnauthorized, "Unauthorized")
		status = false
	}

	if status {
		c.Next() //continue routing
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func check_student_authorization(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string) // get username

	m := make(map[string][]string) //dico
	m["Thibault"] = append(m["Thibault"], "GET", "POST")
	m["Philippe"] = append(m["Philippe"], "GET", "DELETE")

	if val, ok := m[user]; ok { //if dico contains key (user)
		for _, v := range val {
			if v == c.Request.Method { //if they have correct acces right for query
				c.Next() //continue routing
			}
		}
	} else {
		//if get ok else forbidden
		if c.Request.Method == "GET" {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func main() {
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{"foo": "bar", "aristote": "Eucl1de"}))

	authorized.GET("/students", check_student_authorization, getStudents)
	authorized.GET("/students/:id", check_student_authorization, getStudentByID)
	authorized.POST("/students", check_student_authorization, postStudents)
	authorized.DELETE("/students/:id", check_student_authorization, deleteStudentByID)

	authorized.GET("/teachers", check_teacher_authentication, getTeachers)
	authorized.GET("/teachers/:id", check_teacher_authentication, getTeacherByID)
	authorized.POST("/teachers", check_teacher_authentication, postTeachers)
	authorized.DELETE("/teachers/:id", check_teacher_authentication, deleteTeacherByID)

	router.Run()

}
