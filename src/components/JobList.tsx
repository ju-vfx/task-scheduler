import { useState } from "react";
import JobItem from "./JobItem";
import "./JobList.css";

const JobList = ({ allJobs }: { allJobs: any }) => {
  const [selectedJob, setSelectedJob] = useState("");
  if (allJobs.length < 1) {
    return <>No jobs available</>;
  }

  const handleSelectedJob = (id: string) => {
    if (selectedJob === id) {
      setSelectedJob("");
    } else {
      setSelectedJob(id);
    }
  };

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
        {allJobs.map((jobItem: any) => (
          <JobItem
            item={jobItem}
            onSelectItem={handleSelectedJob}
            selectedJob={selectedJob}
            key={jobItem.job_id}
          />
        ))}
      </tbody>
    </table>
  );
};

export default JobList;
