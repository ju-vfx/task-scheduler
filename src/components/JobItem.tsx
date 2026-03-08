import { useState } from "react";
import TaskItem from "./TaskItem";

const JobItem = ({ item, onSelectItem, selectedJob }) => {
  const [selectedTask, setSelectedTask] = useState("");
  const handleSelectedTask = (id: string) => {
    if (selectedTask === id) {
      setSelectedTask("");
    } else {
      setSelectedTask(id);
    }
  };

  let className = "";
  if (item.job_status === "Finished") {
    className = "status-finished";
  } else if (item.job_status === "Error") {
    className = "status-error";
  }

  return (
    <>
      <tr
        key={item.job_id}
        role="button"
        onClick={() => {
          onSelectItem(item.job_id);
        }}
        className={selectedJob === item.job_id ? "table-primary" : ""}
      >
        <td>{item.job_name}</td>
        <td>{item.job_priority}</td>
        <td className={className}>{item.job_status}</td>
        <td>{item.job_created_at}</td>
        <td>{item.job_finished_at}</td>
      </tr>
      {item.job_tasks.length > 0 && selectedJob === item.job_id && (
        <tr>
          <td colSpan="5">
            <table className="table table-sm">
              <thead>
                <tr>
                  <th scope="col">Task</th>
                  <th scope="col">Status</th>
                  <th scope="col">Created At</th>
                  <th scope="col">Finished At</th>
                </tr>
              </thead>
              <tbody>
                {item.job_tasks.map((taskItem) => (
                  <TaskItem
                    item={taskItem}
                    onSelectTask={handleSelectedTask}
                    selectedTask={selectedTask}
                  />
                ))}
              </tbody>
            </table>
          </td>
        </tr>
      )}
    </>
  );
};

export default JobItem;
