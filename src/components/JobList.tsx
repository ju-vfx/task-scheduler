import { useState } from "react";
import { useEffect } from "react";

const JobList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const apiUrl = "http://localhost:8080/api/";

  const fetchJobs = async () => {
    const response = await fetch(apiUrl + "jobs");
    const data = await response.json();
    setDisplayItems(data);
  };

  useEffect(() => {
    fetchJobs();
  }, []);

  return (
    <div className="list-group">
      {displayItems.map((jobItem) => (
        <div key={jobItem.job_id}>
          <button
            type="button"
            className="list-group-item list-group-item-action list-group-item-primary"
          >
            {jobItem.job_name} {jobItem.job_priority} {jobItem.job_status}
          </button>
          <div className="list-group">
            {jobItem.job_tasks.map((taskItem) => (
              <a
                href="#"
                className="list-group-item list-group-item-action list-group-item-secondary"
                key={taskItem.task_id}
              >
                {taskItem.task_name} {taskItem.task_status}
              </a>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

export default JobList;
