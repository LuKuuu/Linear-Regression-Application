package main

import (
	"fmt"
	"github.com/LuKuuu/Kun/LKmath"
	"strings"
	"text/template"
	"time"

	"net/http"
	"strconv"
)


type User struct {
	inputMatrix      LKmath.Matrix
	YMatrix          LKmath.Matrix
	modelMatrix      LKmath.Matrix
	PredictionMatrix LKmath.Matrix
}

func (m *User) Index(w http.ResponseWriter, r *http.Request) { //展示界面

	t, err := template.ParseFiles("./index.html")
	if err != nil {
		panic(err)
		return
	}

	t.Execute(w, nil) //使用http.ResponseWrite输出网页

}

func (m *User) Input(w http.ResponseWriter, r *http.Request) {

	row :=5
	col :=2

	r.ParseForm()
	for i, v := range r.Form {

		if i == "RowNum" {
			row, _ = strconv.Atoi(v[0])

		} else if i == "ColNum" {
			col, _ = strconv.Atoi(v[0])

		}
	}

	if row != 0 && col != 0 {
		m.inputMatrix = LKmath.NewEmptyMatrix(row, col+1)
		m.YMatrix = LKmath.NewEmptyMatrix(row, 1)
	}

	t, err := template.ParseFiles("./input.html")
	if err != nil {
		panic(err)
		return
	}

	table := "<table>"

	for i := 0; i < m.inputMatrix.Row+1; i++ {


		table += "<tr>"
		for j := 0; j < m.inputMatrix.Column+1; j++ {


			table += "<td>"

			if i == 0 && j == 0 {

			} else if i == 0 && j == m.inputMatrix.Column {
				table += "result"
			} else if i == 0 {
				table += fmt.Sprintf("feature %d", j)

			} else if j == 0 {
				table += fmt.Sprintf("example %d", i)
				m.inputMatrix.Cell[i-1][0]=1


			} else {
				table += fmt.Sprintf("<input type='text' name='%d,%d'>", i-1, j)
			}

			table += "</td>"

		}

		table += "</tr>"

	}

	table += "</table><input type=submit value='Submit'>"

	tableData := map[string]interface{}{"table": table}

	t.Execute(w, tableData)

}

func (m *User) Model(w http.ResponseWriter, r *http.Request) { //展示界面

	row := 0
	col := 0
	value := 0.0

	r.ParseForm()
	for i, v := range r.Form {

		row, col = CellLocation(i)
		value, _ = strconv.ParseFloat(v[0], 64)

		//fmt.Printf("\nrow:%d, col:%d, value:%f\n",row,col,value)

		if col == m.inputMatrix.Column {
			m.YMatrix.Cell[row][0] = value
		} else {
			m.inputMatrix.Cell[row][col] = value

		}

	}

	m.inputMatrix.Hprint("\ninput matrix")

	m.YMatrix.Hprint("\ny matrix")

	m.modelMatrix = LKmath.RegularizedNormalEquation(m.inputMatrix, m.YMatrix,0.0001)

	m.modelMatrix.Hprint("model")

	t, err := template.ParseFiles("./model.html") //使用text/template（html/template会导致输出部分被引号包围无法使用）
	if err != nil {
		panic(err)
		return
	}

	table := "<table>"

	for i := 0; i < 2; i++ {
		table += "<tr>"
		for j := 1; j < m.inputMatrix.Column; j++ {

			table += "<td>"

			if i == 0 {
				table += fmt.Sprintf("feature %d", j)
			} else {
				table += fmt.Sprintf("<input type='text' name='%d,%d'>", 0, j)
			}

			table += "</td>"

		}

		table += "</tr>"

	}

	table += "</table><input type=submit value='Submit'>"

	tableData := map[string]interface{}{"table": table}

	t.Execute(w, tableData)

}

func CellLocation(info string) (int, int) {
	comma := strings.Index(info, ",")

	rowNum, _ := strconv.Atoi(info[0:comma])
	colNum, err := strconv.Atoi(info[comma+1 : len(info)])
	if err != nil {
		fmt.Printf("%v", err)
	}

	return rowNum, colNum
}

func (m *User) Result(w http.ResponseWriter, r *http.Request) { //展示界面

	m.PredictionMatrix = LKmath.NewEmptyMatrix( 1,m.inputMatrix.Column)
	m.PredictionMatrix.Cell[0][1] = 1

	col := 0
	row := 0
	value := 0.0

	r.ParseForm()
	for i, v := range r.Form {

		row, col = CellLocation(i)
		value, _ = strconv.ParseFloat(v[0], 64)

		fmt.Printf("\nrow:%d, col:%d, value:%f\n",row,col,value)

		m.PredictionMatrix.Cell[0][col] = value

	}



	resultMatrix :=LKmath.MatrixMultiplication(m.PredictionMatrix,m.modelMatrix)

	t, err := template.ParseFiles("./result.html") //使用text/template（html/template会导致输出部分被引号包围无法使用）
	if err != nil {
		panic(err)
		return
	}


	tableData := map[string]interface{}{"result": fmt.Sprintf("%f",resultMatrix.Cell[0][0])}

	t.Execute(w, tableData)

}

func main() {

	fmt.Printf("starting...\n")

	fmt.Printf("you should go to 127.0.0.1:2000\n")

	ur := User{}
	http.HandleFunc("/", ur.Index)
	http.HandleFunc("/input", ur.Input)
	http.HandleFunc("/model", ur.Model)
	http.HandleFunc("/result", ur.Result)


	err := http.ListenAndServe(":800", nil) //创建服务器
	if err != nil {
		fmt.Printf("%v", err)
		time.Sleep(time.Second*10)
	}
}
