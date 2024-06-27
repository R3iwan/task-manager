package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Task represents a task with an ID, name, and description.
type Task struct {
	Id          int64
	Name        string
	Description string
}

var tasks []Task     // Slice to store all tasks.
var nextId int64 = 1 // Counter for generating unique IDs.

// main function is the entry point of the program.
func main() {
	err := loadTasks("tasks.json")
	if err != nil {
		fmt.Println("Error loading tasks", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		// Display menu options.
		var choice int
		fmt.Println("Options:")
		fmt.Println("1.Add a task")
		fmt.Println("2.Update a task")
		fmt.Println("3.Show tasks")
		fmt.Println("4.Delete a task")
		fmt.Println("5.Exit")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Please enter a valid option")
			continue
		}
		var id int64
		switch choice {
		case 1:
			task := Task{}
			addTask(&task, reader)
			tasks = append(tasks, task)
			saveTasks("tasks.json")
		case 2:
			fmt.Println("Please enter a task ID to update name and description")
			fmt.Scanln(&id)
			updateTask(id)
			saveTasks("tasks.json")
		case 3:
			showTasks()
		case 4:
			fmt.Println("Please enter a task ID to delete task")
			fmt.Scanln(&id)
			deleteTask(id)
			saveTasks("tasks.json")
		case 5:
			return
		default:
			fmt.Println("Please enter a valid option")
		}
	}
}

// addTask adds a new task to the tasks slice.
func addTask(t *Task, reader *bufio.Reader) {
	t.Id = nextId
	nextId++

	fmt.Println("Please enter your name:")
	t.Name, _ = reader.ReadString('\n')
	t.Name = strings.TrimSpace(t.Name)

	fmt.Println("Please enter your description:")
	t.Description, _ = reader.ReadString('\n')
	t.Description = strings.TrimSpace(t.Description)
}

// updateTask updates the name and description of a task with the given ID.
func updateTask(id int64) {
	for i, t := range tasks {
		if t.Id == id {
			fmt.Println("Please enter your new name:")
			fmt.Scanln(&tasks[i].Name)
			fmt.Println("Please enter your new description:")
			fmt.Scanln(&tasks[i].Description)
			fmt.Println("Task updated successfully")
			return
		}
	}
	fmt.Println("Task not found")
}

// showTasks prints all tasks in the tasks slice.
func showTasks() {
	for _, t := range tasks {
		fmt.Printf("ID: %d, Name: %s, Description: %s\n", t.Id, t.Name, t.Description)
	}
}

// deleteTask deletes a task with the given ID from the tasks slice.
func deleteTask(id int64) {
	for i, t := range tasks {
		if t.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Println("Task deleted successfully")
			return
		}
	}
	fmt.Println("Task not found")
}

// saveTasks saves all tasks in the tasks slice to a JSON file.
func saveTasks(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding tasks", err)
		return
	}
}

// loadTasks loads tasks from a JSON file into the tasks slice.
func loadTasks(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		if err == io.EOF {
			tasks = []Task{}
			return nil
		}
		return err
	}

	for _, t := range tasks {
		if t.Id >= nextId {
			nextId = t.Id + 1
		}
	}

	return nil
}
