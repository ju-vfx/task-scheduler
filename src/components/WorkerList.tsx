import { useState } from "react";
import { useEffect } from "react";

const WorkerList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const apiUrl = "http://localhost:8080/api/";

  const fetchWorkers = async () => {
    const response = await fetch(apiUrl + "workers");
    const data = await response.json();
    setDisplayItems(data);
  };

  useEffect(() => {
    fetchWorkers();
  }, []);

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Host</th>
          <th scope="col">Status</th>
          <th scope="col">Last Seen</th>
          <th scope="col">Connected At</th>
        </tr>
      </thead>
      <tbody>
        {displayItems.map((workerItem) => (
          <tr>
            <td>{workerItem.host}</td>
            <td>{workerItem.status}</td>
            <td>{workerItem.last_seen_at}</td>
            <td>{workerItem.connected_at}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default WorkerList;
