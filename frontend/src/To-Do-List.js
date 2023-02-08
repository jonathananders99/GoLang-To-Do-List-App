import { Component } from "react";
import React from "react"
import axios from "axios";

let endpoint = "http://localhost:9000";

class ToDoList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            title: "",
            details: "",
            priority: "",
            allTasks: []
        };
    }


    componentDidMount() {
        this.getTasks();
    }


    onChange = (event) => {
        this.setState({
            [event.target.name]: event.target.value
        });
    };


    //when a new task is added
    onSubmit = (event) => {
        event.preventDefault();
        if (this.state.title  &&  this.state.priority  &&  this.state.details) {
            axios
            .post(
                endpoint + "/api/task",
                {
                    "title": this.state.title,
                    "details": this.state.details,
                    "priority": this.state.priority
                }
            )
            .then((res) => {
                this.getTasks();
                this.setState({title: "", details: "", priority: ""});
            });
        }
    };    


    //updates all tasks
    getTasks = () => {
        axios.get(endpoint + "/api/task").then((res) => {
            if (res.data) {
                this.setState({
                    items: res.data.map((item) => {
                        let fontStyle = {};
                        let textDecoration = "none";
                        if (item.completed) {
                            fontStyle = { fontStyle: "italic" };
                            textDecoration = "line-through";
                        }
                        return (
                            <div className="task_box" key={item._id}>
                                <div className="task_title" style={{ ...fontStyle, textDecoration }}>
                                    {item.title}
                                </div>

                                <div className="task_details" style={{ ...fontStyle, textDecoration }}>
                                    {item.details}
                                </div>

                                <div className="task_priority" style={{ ...fontStyle, textDecoration }}>
                                    {item.priority}
                                </div>

                                <div className="task_actions">

                                    <div className="completed" onClick={() => this.completeTask(item._id)}>
                                        Complete
                                    </div>
                                
                                    <div className="not_complete" onClick={() => this.undoTask(item._id)}>
                                        Undo
                                    </div>
                                
                                    <div className="delete" onClick={() => this.deleteTask(item._id)}>
                                        Delete
                                    </div>

                                </div>
                            </div>
                        );
                    }),
                });
            } 
            //No Tasks
            else {
                this.setState({items: []});
            }
        });
    };


    //complete a task
    completeTask = (id) => {
        axios
        .put(endpoint + "/api/task/" + id, {
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            }
        })
        .then((res) => {
            this.getTasks();
        });
    };


    //undo a completed task
    undoTask = (id) => {
        axios
        .put(endpoint + "/api/undoTask/" + id, {
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        })
        .then((res) => {
            this.getTasks();
        });
    };


    // delete a task
    deleteTask = (id) => {
        axios
        .delete(endpoint + "/api/deleteTask/" + id, 
        {
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        })
        .then((res) => {
            this.getTasks();
        });
    };

    
    //renders the page
    render() {
        return (
        <div>
            <div className="header">
                TO DO LIST
            </div>
            <div className="task_create">
                <form onSubmit={this.onSubmit}>
                    <div className="new_task_button" onClick={this.onSubmit}>Create Task</div>
                    <input 
                        className="title_box"
                        type="text"
                        name="title"
                        onChange={this.onChange}
                        value={this.state.title}
                        placeholder="New Title"
                        required
                    />
                    <input 
                        className="details_box"
                        type="text"
                        name="details"
                        onChange={this.onChange}
                        value={this.state.details}
                        placeholder="New Details"
                        required
                    />
                    <input 
                        className="priority_box"
                        type="text"
                        name="priority"
                        onChange={this.onChange}
                        value={this.state.priority}
                        placeholder="New Priority"
                        required
                    />
                </form>
            </div>
            <div className="all_tasks">
                {this.state.items}
            </div>
        </div>
        );
    }
}

export default ToDoList;
