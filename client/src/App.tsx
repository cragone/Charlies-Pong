import "./App.css";
import Homepage from "./pages/Homepage";
import { Route, Routes, BrowserRouter as Router } from "react-router-dom";

function App() {
  return (
    <>
      <Router>
        <Routes>
          <Route path="/" element={<Homepage />} />
        </Routes>
      </Router>
    </>
  );
}

export default App;
