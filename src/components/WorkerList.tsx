import { useState } from "react";
import { useEffect } from "react";

const WorkerList = () => {
  const [displayItems, setDisplayItems] = useState([]);
  const apiUrl = "http://localhost:8080/api/";

  const fetchWorkers = async () => {
    try {
      const response = await fetch(apiUrl + "workers");
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const data = await response.json();
      setDisplayItems(data);
    } catch (error) {}
  };

  useEffect(() => {
    fetchWorkers();
    const interval = setInterval(() => {
      fetchWorkers();
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <table className="table">
      <thead>
        <tr href="#">
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
