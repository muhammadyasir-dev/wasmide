import { BrowserRouter, Route, Routes } from "react-router-dom";
import CodeEditor from "./pages/CodeEditor";
import LandingPage from "./pages/LandingPage";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Output from "./pages/Output";
import { GoogleOAuthProvider } from "@react-oauth/google";
const App: React.FC = () => {
  return (
    <GoogleOAuthProvider clientId="<your_client_id>">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/editor" element={<CodeEditor />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/output" element={<Output />} />
        </Routes>
      </BrowserRouter>
    </GoogleOAuthProvider>
  );
};

export default App;
