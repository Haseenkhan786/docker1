package main

import ( 
	"github.com/gin-contrib/cors"
	"database/sql"
	"fmt"
	"net/http"
    "mime/multipart"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
    "os"
)
var err = godotenv.Load()

var (

    databaseHost     = os.Getenv("host")

    databaseUser     = os.Getenv("user")

    databasePassword = os.Getenv("password")

    databasename     = os.Getenv("dbname")
	x=os.Getenv("localHost")
	resumeLOcation=os.Getenv("fileLocation")
)
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "haseen"
// 	password = "1302001"
// 	dbname   = "postgres"
// )

type EmployeeDetail struct {
	FullName      string `form:"name" `
	Gender        string `form:"gender"`
	From_date     string `form:"fromDate" `
	To_date       string `form:"toDate"`
	Phone         int `form:"phone"`
	Resume *multipart.FileHeader `form:"selectedFile" `
	Email         string `form:"email"`
}
type CompanyEmployeeDetails struct {
	FullName      string `form:"name"`
	Gender        string `form:"gender"`
	From_date     string `form:"fromDate"`
	To_date       string `form:"toDate"`
	Phone         int `form:"phone"`
	Resume	string `form:"selectedFile"				`
	Email         string `form:"email"`
}
const port=5432

var employees = []EmployeeDetail{}




func main() {
	
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"} // Replace with the URL of your Angular application
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	router.Use(cors.New(config))
	router.GET("/getemployees", getEmployees)
	router.Static("/static", "./static")
	

	router.POST("/postemployees", postEmployees)

	router.Run(x)


	
}

// POST
func postEmployees(c *gin.Context) {
	var newEmployee EmployeeDetail

	if err := c.ShouldBind(&newEmployee); err != nil {
		c.String(http.StatusBadRequest, "Failed to parse form data")
		return
	}
	
	// employees = append(employees, newEmployee)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseHost, port, databaseUser, databasePassword, databasename)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()
	
	
	
	fileLink := ""

    if newEmployee.Resume != nil {

        // Handle file upload and save the file to the desired folder (myfolder)

        filePath := resumeLOcation + newEmployee.Resume.Filename

        if err := c.SaveUploadedFile(newEmployee.Resume, filePath); err != nil {

            c.String(http.StatusInternalServerError, "Failed to save file on the server")

            return

        }

        fileLink = filePath

    }

	insertDynStmt := `INSERT INTO "employee_detail"  VALUES ($1, $2, $3, $4, $5, $6,$7)`

	_, err = db.Exec(insertDynStmt, newEmployee.FullName, newEmployee.Gender, newEmployee.From_date, newEmployee.To_date, newEmployee.Phone, fileLink, newEmployee.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to insert form data into the database")
		return
	}

	c.String(http.StatusOK, "Leave application submitted successfully!")
}

// GET
func getEmployees(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseHost, port, databaseUser, databasePassword, databasename)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM employee_detail")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee data from the database"})
		return
	}
	defer rows.Close()

	var employees []CompanyEmployeeDetails

	for rows.Next() {
		var emp CompanyEmployeeDetails
		err := rows.Scan(&emp.FullName, &emp.Gender, &emp.From_date, &emp.To_date, &emp.Phone,&emp.Resume , &emp.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan employee data from the database"})
			return
		}
		employees = append(employees, emp)
	}

	c.IndentedJSON(http.StatusOK, employees)
}


// Add is our function that sums two integers
func Add(x, y int) int{
	return x + y
}

func oddeven(a int)int {
	if a%2==0{
		return 1 ;
	}
	return -1;
}

