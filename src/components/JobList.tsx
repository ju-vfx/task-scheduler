import { useState } from "react";
import { useEffect } from "react";
import "./JobList.css";

const JobList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const [selectedJob, setSelectedJob] = useState("");
  const [selectedTask, setSelectedTask] = useState("");

  const apiUrl = "http://localhost:8080/api/";

  const fetchJobs = async () => {
    try {
      const response = await fetch(apiUrl + "jobs");
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const data = await response.json();
      setDisplayItems(data);
    } catch (error) {}
  };

  useEffect(() => {
    fetchJobs();
    const interval = setInterval(() => {
      fetchJobs();
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col">Priority</th>
          <th scope="col">Status</th>
          <th scope="col">Created At</th>
          <th scope="col">Finished At</th>
        </tr>
      </thead>
      <tbody>
        {displayItems.map((jobItem) => (
          <>
            <tr
              key={jobItem.job_id}
              role="button"
              onClick={() => {
                setSelectedJob(jobItem.job_id);
              }}
              className={selectedJob === jobItem.job_id ? "table-primary" : ""}
            >
              <td>{jobItem.job_name}</td>
              <td>{jobItem.job_priority}</td>
              <td>{jobItem.job_status}</td>
              <td>{jobItem.job_created_at}</td>
              <td>{jobItem.job_finished_at}</td>
            </tr>
            {jobItem.job_tasks.length > 0 && selectedJob === jobItem.job_id && (
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
                      {jobItem.job_tasks.map((taskItem) => (
                        <>
                          <tr
                            role="button"
                            onClick={() => {
                              setSelectedTask(taskItem.task_id);
                            }}
                            className={
                              selectedTask === taskItem.task_id
                                ? "table-primary"
                                : ""
                            }
                          >
                            <td>{taskItem.task_name}</td>
                            <td>{taskItem.task_status}</td>
                            <td>{taskItem.task_created_at}</td>
                            <td>{taskItem.task_finished_at}</td>
                          </tr>
                          {selectedTask === taskItem.task_id && (
                            <>
                              <tr>
                                <td colSpan="5">
                                  Command: "{taskItem.task_command}"
                                </td>
                              </tr>
                              {taskItem.task_output != "" && (
                                <tr>
                                  <td className="terminal-output" colSpan="5">
                                    {taskItem.task_output}
                                  </td>
                                </tr>
                              )}
                            </>
                          )}
                        </>
                      ))}
                    </tbody>
                  </table>
                </td>
              </tr>
            )}
          </>
        ))}
      </tbody>
    </table>
  );
};

export default JobList;
