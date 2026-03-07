import { useState } from "react";

const JobList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const apiUrl = "http://localhost:8080/api/";

  const fetchJobs = async () => {
    const response = await fetch(apiUrl + "jobs");
    const data = await response.json();
    setDisplayItems(data.jobs);
  };

  return (
    <>
      <ul className="list-group">
        <li className="list-group-item"></li>
      </ul>
    </>
  );
};

export default JobList;
