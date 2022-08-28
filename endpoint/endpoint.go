package endpoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"program/service"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Courses struct {
	Course []Course `json:"courses"`
}

type Course struct {
	Name string `json:"name"`
	Days Days   `json:"days"`
}

type Days struct {
	Monday    []int `json:"monday"`
	Tuesday   []int `json:"tuesday"`
	Wednesday []int `json:"wednesday"`
	Thursday  []int `json:"thursday"`
	Friday    []int `json:"friday"`
}

type CourseTable struct {
	Monday    []string
	Tuesday   []string
	Wednesday []string
	Thursday  []string
	Friday    []string
}

func PostCourse(w http.ResponseWriter, r *http.Request) {
	//Post an array of courses
	//WARNING: Overwrites files
	courses := Courses{ //default values for courses array
		Course: []Course{
			{
				Name: "default",
				Days: Days{
					Monday:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Tuesday:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Wednesday: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Thursday:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Friday:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
			},
		},
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &courses)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	fmt.Fprintf(w, "Added to JSON file\n")

	//Display added array and its contents
	for i := 0; i < len(courses.Course); i++ {
		fmt.Fprintf(w, "Name = %s\n", courses.Course[i].Name)
		fmt.Fprintf(w, "Monday = %d\n", courses.Course[i].Days.Monday)
		fmt.Fprintf(w, "Tuesday = %d\n", courses.Course[i].Days.Tuesday)
		fmt.Fprintf(w, "Wednesday = %d\n", courses.Course[i].Days.Wednesday)
		fmt.Fprintf(w, "Thursday = %d\n", courses.Course[i].Days.Thursday)
		fmt.Fprintf(w, "Friday = %d\n", courses.Course[i].Days.Friday)
	}

	file, _ := json.MarshalIndent(courses, "", " 	")
	service.WriteFile(file)
}

func ViewCourses(w http.ResponseWriter, r *http.Request) {
	fileContent := service.ReadFile()

	var courses Courses
	json.Unmarshal(fileContent, &courses)

	//Display added array and its contents
	for i := 0; i < len(courses.Course); i++ {
		fmt.Fprintf(w, "Name = %s\n", courses.Course[i].Name)
		fmt.Fprintf(w, "Monday = %d\n", courses.Course[i].Days.Monday)
		fmt.Fprintf(w, "Tuesday = %d\n", courses.Course[i].Days.Tuesday)
		fmt.Fprintf(w, "Wednesday = %d\n", courses.Course[i].Days.Wednesday)
		fmt.Fprintf(w, "Thursday = %d\n", courses.Course[i].Days.Thursday)
		fmt.Fprintf(w, "Friday = %d\n", courses.Course[i].Days.Friday)
	}
}

func PlanCourses(w http.ResponseWriter, r *http.Request) {
	courseTable := CourseTable{
		Monday:    []string{"Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty"},
		Tuesday:   []string{"Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty"},
		Wednesday: []string{"Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty"},
		Thursday:  []string{"Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty"},
		Friday:    []string{"Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty", "Empty"},
	}

	fileContent := service.ReadFile()

	//Assign shortened versions of table parts
	mon := courseTable.Monday
	tue := courseTable.Tuesday
	wed := courseTable.Wednesday
	thu := courseTable.Thursday
	fri := courseTable.Friday

	var courses Courses
	json.Unmarshal(fileContent, &courses)

	//Checks for conflicts
	for i := 0; i < len(courses.Course); i++ {
		days := courses.Course[i].Days

		mon = checkTimeslots(mon, days.Monday, "Monday", courses.Course[i].Name)
		tue = checkTimeslots(tue, days.Tuesday, "Tuesday", courses.Course[i].Name)
		wed = checkTimeslots(wed, days.Wednesday, "Wednesday", courses.Course[i].Name)
		thu = checkTimeslots(thu, days.Thursday, "Thursday", courses.Course[i].Name)
		fri = checkTimeslots(fri, days.Friday, "Friday", courses.Course[i].Name)
	}

	//Makes a better looking table
	prettyTable := table.NewWriter()
	prettyTable.SetOutputMirror(os.Stdout)
	prettyTable.AppendHeader(table.Row{"#", "9:30-10:20", "10:30-11:20", "11:30-12:20", "12:30-13:20", "13:30-14:20", "14:30-15:20",
		"15:30-16:20", "16:30-17:20", "17:30-18:20", "18:30-19:20"})
	prettyTable.AppendRows([]table.Row{
		{"Monday", mon[0], mon[1], mon[2], mon[3], mon[4], mon[5], mon[6], mon[7], mon[8], mon[9]},
		{"Tuesday", tue[0], tue[1], tue[2], tue[3], tue[4], tue[5], tue[6], tue[7], tue[8], tue[9]},
		{"Wednesday", wed[0], wed[1], wed[2], wed[3], wed[4], wed[5], wed[6], wed[7], wed[8], wed[9]},
		{"Thursday", thu[0], thu[1], thu[2], thu[3], thu[4], thu[5], thu[6], thu[7], thu[8], thu[9]},
		{"Friday", fri[0], fri[1], fri[2], fri[3], fri[4], fri[5], fri[6], fri[7], fri[8], fri[9]},
	},
	)
	prettyTable.Style().Options.SeparateRows = true

	fmt.Fprintf(w, "Course Table:\n%s", prettyTable.Render())
}

func checkTimeslots(timeslots []string, days []int, dayName string, courseName string) []string {
	for i := 0; i < 10; i++ {
		if timeslots[i] != "Empty" && days[i] == 1 {
			fmt.Print("Timeslot on " + dayName + " already full by " + courseName)
		} else if timeslots[i] == "Empty" && days[i] == 1 {
			timeslots[i] = courseName
		}
	}
	return timeslots
}
