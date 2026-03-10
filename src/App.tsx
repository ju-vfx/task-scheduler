import JobList from "./components/JobList";
import NavBar from "./components/NavBar";
import { useState } from "react";
import WorkerList from "./components/WorkerList";

function App() {
  const [selection, setSelection] = useState("Jobs");
  const [isConnected, setIsConnected] = useState(false);
  const [workers, setWorkers] = useState<any>([]);
  const [jobs, setJobs] = useState<any>([]);

  const connectToWebsocket = () => {
    let socket = new WebSocket("ws://localhost:8080/api/registerClients");
    console.log("Attempting websocket connection");
    socket.onopen = () => {
      console.log("Successful connection to server");
      setIsConnected(true);
    };

    socket.onclose = (event) => {
      console.log("Connection to server closed: ", event);
      setIsConnected(false);
    };

    socket.onerror = (error) => {
      console.log("Socket Error: ", error);
      setIsConnected(false);
    };

    socket.onmessage = (msg) => {
      let data = JSON.parse(msg.data);
      if (data["message_type"] === "workers") {
        setWorkers(data["payload"]);
      }
      if (data["message_type"] === "jobs") {
        setJobs(data["payload"]);
      }
    };
  };

  if (!isConnected) {
    connectToWebsocket();
  }
  return (
    <>
      <NavBar
        onSelectItem={(selItem) => {
          setSelection(selItem);
        }}
      />
      {selection === "Jobs" ? (
        <JobList allJobs={jobs} />
      ) : (
        <WorkerList allWorkers={workers} />
      )}
      {isConnected === false && (
        <>
          <p>No connection to server</p>
          <button onClick={connectToWebsocket}>Reconnect</button>
        </>
      )}
    </>
  );
}

export default App;
