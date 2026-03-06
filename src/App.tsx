import NavBar from "./components/NavBar";
import ListGroup from "./components/ListGroup";
import { useState } from "react";

function App() {
  const [selection, setSelection] = useState("Jobs");

  const apiUrl = "http://localhost:8080/api/";
  const [displayItems, setDisplayItems] = useState([]);

  const fetchJobs = async () => {
    const response = await fetch(apiUrl + "jobs");
    const data = await response.json();
    setDisplayItems(data.jobs);
  };

  const fetchWorkers = async () => {
    setDisplayItems([]);
  };

  const handleSelectItem = (item: string) => {
    setSelection(item);
    if (item == "Jobs") {
      fetchJobs();
    } else {
      fetchWorkers();
    }
  };

  return (
    <>
      <NavBar onSelectItem={handleSelectItem} />
      <ListGroup items={displayItems} />
    </>
  );
}

export default App;
