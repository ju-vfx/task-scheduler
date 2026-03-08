import { useState } from "react";
import { useEffect } from "react";
import JobItem from "./JobItem";
import "./JobList.css";

const JobList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const [selectedJob, setSelectedJob] = useState("");

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

  const handleSelectedJob = (id: string) => {
    if (selectedJob === id) {
      setSelectedJob("");
    } else {
      setSelectedJob(id);
    }
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
          <JobItem
            item={jobItem}
            onSelectItem={handleSelectedJob}
            selectedJob={selectedJob}
          />
        ))}
      </tbody>
    </table>
  );
};

export default JobList;
