import JobList from "./components/JobList";
import NavBar from "./components/NavBar";
import { useState } from "react";
import WorkerList from "./components/WorkerList";

function App() {
  const [selection, setSelection] = useState("Jobs");

  return (
    <>
      <NavBar
        onSelectItem={(selItem) => {
          setSelection(selItem);
        }}
      />
      {selection === "Jobs" ? <JobList /> : <WorkerList />}
    </>
  );
}

export default App;
