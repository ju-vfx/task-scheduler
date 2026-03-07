import { useState } from "react";
import { useEffect } from "react";

const JobList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const [selectedJob, setSelectedJob] = useState("");

  const apiUrl = "http://localhost:8080/api/";

  const fetchJobs = async () => {
    const response = await fetch(apiUrl + "jobs");
    const data = await response.json();
    setDisplayItems(data);
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
                      </tr>
                    </thead>
                    <tbody>
                      {jobItem.job_tasks.map((taskItem) => (
                        <tr>
                          <td>{taskItem.task_name}</td>
                          <td>{taskItem.task_status}</td>
                        </tr>
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
