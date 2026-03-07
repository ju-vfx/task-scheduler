import { useState } from "react";

const WorkerList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const apiUrl = "http://localhost:8080/api/";

  const fetchWorkers = async () => {
    const response = await fetch(apiUrl + "workers");
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

export default WorkerList;
